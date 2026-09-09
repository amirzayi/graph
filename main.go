package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/amirzayi/graph/audit"
	"github.com/amirzayi/graph/models"
	"github.com/amirzayi/graph/router"
	"github.com/amirzayi/graph/task"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title           Graph swagger api
// @version         1.0
// @description     This is graph interview task api documents.

// @contact.name   developer
// @contact.url    https://amirzayi.github.io/
// @contact.email  mirzayi994@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8010
// @BasePath  /api/v1/tasks

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	_ = godotenv.Load()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatal(err)
	}
	_ = db.AutoMigrate(&models.Task{})

	repo := task.NewSQLRepository(db)
	auditLog := audit.NewIOWriter(os.Stdout)
	taskSvc := task.NewService(repo, auditLog)

	ginrouter := gin.Default()

	router.Bootstrap(ginrouter, taskSvc)

	httpsrv := http.Server{
		Addr:    ":8010",
		Handler: ginrouter,
	}
	log.Fatal(httpsrv.ListenAndServe())
}
