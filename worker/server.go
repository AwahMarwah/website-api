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
	productRepo "website-api/repository/product"
	product_variant "website-api/repository/product-variant"
	userRepo "website-api/repository/user"
	orderService "website-api/service/order"
	midtransProvider "website-api/third-party/provider/midtrans"
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

	if err := scheduler.Start(); err != nil {
		log.Fatalf("failed to start scheduler: %v", err)
	}
	defer scheduler.Shutdown()

	if err := server.Run(mux); err != nil {
		log.Fatal(err)
	}
}

func handleCancelExpiredOrders(ctx context.Context, t *asynq.Task) error {
	db, err := database.Open()
	if err != nil {
		return err
	}
	defer db.SqlDb.Close()

	svc := orderService.NewService(
		productRepo.NewRepo(db.GormDb),
		product_variant.NewRepo(db.GormDb),
		orderRepo.NewRepo(db.GormDb),
		userRepo.NewRepo(db.GormDb),
		transaction.NewTransactionManager(db.GormDb),
		midtransProvider.NewClient(),
		NewRedisClient(),
	)

	_, err = svc.CancelExpiredOrders()
	if err != nil {
		return err
	}

	log.Printf("cancel expired orders task finished at %s", time.Now().Format(time.RFC3339))
	return nil
}