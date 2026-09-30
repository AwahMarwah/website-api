package worker

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"website-api/database"
	"website-api/database/transaction"
	orderRepo "website-api/repository/order"
	merchantRepo "website-api/repository/merchant"
	productRepo "website-api/repository/product"
	product_variant "website-api/repository/product-variant"
	reviewRepo "website-api/repository/review"
	userRepo "website-api/repository/user"
	userAddressRepo "website-api/repository/user_address"
	voucherRepo "website-api/repository/voucher"
	settlementRepo "website-api/repository/settlement"
	orderService "website-api/service/order"
	settlementService "website-api/service/settlement"
	midtransProvider "website-api/third-party/provider/midtrans"
	rajaongkirProvider "website-api/third-party/provider/rajaongkir"
	"website-api/task"

	"github.com/hibiken/asynq"
)

func StartWorker() {
	addr := fmt.Sprintf(
		"%s:%s",
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
	)

	redisOpt := asynq.RedisClientOpt{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()

	mux.HandleFunc(task.TypeSendResetPassword, HandleResetPassword)
	mux.HandleFunc(task.TypeSendPaymentSuccess, HandlePaymentSuccess)

	// Handler auto-cancel order yang sudah expired.
	// DB dibuka saat task dieksekusi agar worker tetap ringan.
	mux.HandleFunc(task.TypeCancelExpiredOrders, handleCancelExpiredOrders)

	// Rekonsiliasi pembayaran untuk order yang webhook-nya tidak pernah sampai.
	mux.HandleFunc(task.TypeReconcilePayments, handleReconcilePayments)

	// Scheduler untuk menjalankan auto-cancel order expired secara periodik.
	scheduler := asynq.NewScheduler(redisOpt, &asynq.SchedulerOpts{})
	cancelTask, err := task.NewCancelExpiredOrdersTask()
	if err != nil {
		log.Fatalf("failed to create cancel expired orders task: %v", err)
	}
	_, err = scheduler.Register("*/5 * * * *", cancelTask)
	if err != nil {
		log.Fatalf("failed to register cancel expired orders scheduler: %v", err)
	}

	// Rekonsiliasi berjalan lebih jarang karena setiap pemanggilan memakai kuota API provider.
	reconcileTask, err := task.NewReconcilePaymentsTask()
	if err != nil {
		log.Fatalf("failed to create reconcile payments task: %v", err)
	}
	if _, err = scheduler.Register("*/2 * * * *", reconcileTask); err != nil {
		log.Fatalf("failed to register reconcile payments scheduler: %v", err)
	}

	if err := scheduler.Start(); err != nil {
		log.Fatalf("failed to start scheduler: %v", err)
	}
	defer scheduler.Shutdown()

	if err := server.Run(mux); err != nil {
		log.Fatal(err)
	}
}

func handleCancelExpiredOrders(ctx context.Context, t *asynq.Task) error {
	svc, cleanup, err := newOrderService()
	if err != nil {
		return err
	}
	defer cleanup()

	if _, err = svc.CancelExpiredOrders(); err != nil {
		return err
	}

	log.Printf("cancel expired orders task finished at %s", time.Now().Format(time.RFC3339))
	return nil
}

func handleReconcilePayments(ctx context.Context, t *asynq.Task) error {
	svc, cleanup, err := newOrderService()
	if err != nil {
		return err
	}
	defer cleanup()

	updated, err := svc.ReconcilePayments()
	if err != nil {
		return err
	}

	if updated > 0 {
		log.Printf("reconcile payments: %d order diperbarui pada %s", updated, time.Now().Format(time.RFC3339))
	}
	return nil
}

// newOrderService menyusun order service dengan koneksi database baru.
// Koneksi dibuka per-task supaya worker utama tidak memegang koneksi database
// sepanjang waktu tidur.
func newOrderService() (orderService.IService, func(), error) {
	db, err := database.Open()
	if err != nil {
		return nil, func() {}, err
	}

	svc := orderService.NewService(
		productRepo.NewRepo(db.GormDb),
		product_variant.NewRepo(db.GormDb),
		orderRepo.NewRepo(db.GormDb),
		userRepo.NewRepo(db.GormDb),
		userAddressRepo.NewRepo(db.GormDb),
		merchantRepo.NewRepo(db.GormDb),
		reviewRepo.NewRepo(db.GormDb),
		voucherRepo.NewRepo(db.GormDb),
		settlementService.NewService(settlementRepo.NewRepo(db.GormDb)),
		transaction.NewTransactionManager(db.GormDb),
		midtransProvider.NewClient(),
		rajaongkirProvider.NewClient(),
		NewRedisClient(),
	)
	return svc, func() { db.SqlDb.Close() }, nil
}