package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	moduleName, err := getModuleName()
	if err != nil {
		fmt.Println("❌ Fail read go.mod:", err)
		return
	}

	if err := createFiles(moduleName); err != nil {
		fmt.Println("❌ Error creating files:", err)
		return
	}

	fmt.Println("✅ Project structure generated successfully.")
}

// ========== MODULE NAME PARSER ==========
func getModuleName() (string, error) {
	file, err := os.Open("go.mod")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}
	return "", fmt.Errorf("module name not found in go.mod")
}

// ========== FILE STRUCTURE CREATOR ==========
func createFiles(moduleName string) error {
	filesToCreate := map[string]string{
		"main.go": `package main

import (
	"log"
	"` + moduleName + `/database"
	"` + moduleName + `/router"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db, err := database.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = db.SqlDb.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	if err = router.Run(db); err != nil {
		log.Fatal(err)
	}
}
`,
		"database/database.go": `package database

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	GormDb *gorm.DB
	SqlDb  *sql.DB
}

func Open() (db DB, err error) {
	db.SqlDb, err = sql.Open("postgres", os.Getenv("DATABASE_DSN"))
	if err != nil {
		return
	}
	if db.GormDb, err = gorm.Open(postgres.New(postgres.Config{Conn: db.SqlDb}), &gorm.Config{}); err != nil {
		return
	}
	return
}
`,
		"router/router.go": `package router

import (
	"` + moduleName + `/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run(db database.DB) (err error) {
	router := gin.Default()
	router.Use(cors.New(corsConfig))
	return router.Run()
}
`,
		"router/cors.go": `package router

import "github.com/gin-contrib/cors"

var corsConfig = cors.Config{
	AllowAllOrigins: true,
	AllowHeaders: []string{"Authorization", "Content-Type"},
	AllowMethods: []string{"DELETE", "GET", "POST", "PUT"},
}
`,
		"library/response/error.go": `package response

import "github.com/gin-gonic/gin"

func Error(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, response{
		Data:    nil,
		Message: message,
	})
}
`,
		"library/response/response.go": `package response

type (
	page struct {
		Current int   ` + "`json:\"current\"`" + `
		Size    int   ` + "`json:\"size\"`" + `
		Total   int64 ` + "`json:\"total\"`" + `
	}

	response struct {
		Data    any    ` + "`json:\"data\"`" + `
		Message string ` + "`json:\"message\"`" + `
	}

	responseWithPage struct {
		Data    any    ` + "`json:\"data\"`" + `
		Message string ` + "`json:\"message\"`" + `
		Page    page   ` + "`json:\"page\"`" + `
	}
)
`,
		"library/response/success.go": `package response

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func Success(ctx *gin.Context, statusCode int, message string, data any) {
	if message == "" {
		message = http.StatusText(statusCode)
	}
	ctx.JSON(statusCode, response{
		Data:    data,
		Message: message,
	})
}
`,
		"library/response/success_with_page.go": `package response

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func SuccessWithPage(ctx *gin.Context, statusCode int, message string, data any, current, size int, total int64) {
	if message == "" {
		message = http.StatusText(statusCode)
	}
	ctx.JSON(statusCode, responseWithPage{
		Data:    data,
		Message: message,
		Page: page{
			Current: current,
			Size:    size,
			Total:   total,
		},
	})
}
`,
		"library/pagination/pagination.go": `package pagination

func Offset(limit, page *int) int {
	if *limit == 0 {
		*limit = 10
	}
	if *page == 0 {
		*page = 1
	}
	return (*page - 1) * *limit
}
`,
	}

	for path, content := range filesToCreate {
		dir := getDir(path)
		if dir != "" {
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create file %s: %w", path, err)
		}
		fmt.Println("📄 Created:", path)
	}

	return nil
}

// ========== HELPER ==========
func getDir(path string) string {
	if idx := strings.LastIndex(path, "/"); idx != -1 {
		return path[:idx]
	}
	return ""
}
