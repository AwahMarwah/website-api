package middleware

import (
	"net/http"
	"website-api/database/transaction"
	"website-api/library/response"
	menuRepo "website-api/repository/menu"
	permissionRepo "website-api/repository/permission"
	roleRepo "website-api/repository/role"
	"website-api/service/menu"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Permission membatasi akses ke resource yang memerlukan permission tertentu.
// Harus dipanggil SETELAH AuthMiddleware (membaca role_name dari context).
func Permission(db *gorm.DB, name string) gin.HandlerFunc {
	svc := menu.NewService(
		menuRepo.NewRepo(db),
		roleRepo.NewRepo(db),
		permissionRepo.NewRepo(db),
		transaction.NewTransactionManager(db),
	)

	return func(ctx *gin.Context) {
		roleName, exists := ctx.Get("role_name")
		if !exists {
			response.Error(ctx, http.StatusForbidden, "forbidden")
			ctx.Abort()
			return
		}

		roleStr, ok := roleName.(string)
		if !ok {
			response.Error(ctx, http.StatusForbidden, "forbidden")
			ctx.Abort()
			return
		}

		allowed, err := svc.HasPermission(roleStr, name)
		if err != nil {
			response.Error(ctx, http.StatusForbidden, "forbidden")
			ctx.Abort()
			return
		}
		if !allowed {
			response.Error(ctx, http.StatusForbidden, "forbidden")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// RequireRole membatasi akses hanya untuk role tertentu (mis. admin, super_admin).
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool)
	for _, r := range roles {
		allowed[r] = true
	}
	return func(ctx *gin.Context) {
		roleName, exists := ctx.Get("role_name")
		if !exists {
			response.Error(ctx, http.StatusForbidden, "forbidden")
			ctx.Abort()
			return
		}
		roleStr, _ := roleName.(string)
		if !allowed[roleStr] {
			response.Error(ctx, http.StatusForbidden, "forbidden")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}