package main

import (
	"log"
	"net/http"
	"os"

	"github.com/amirzayi/graph/audit"
	"github.com/amirzayi/graph/models"
	"github.com/amirzayi/graph/router"
	"github.com/amirzayi/graph/task"
	"github.com/gin-gonic/gin"
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
	db, err := gorm.Open(postgres.Open("postgres://amir:mirzaei@localhost:5432/task?sslmode=disable"))
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&models.Task{})

	repo := task.NewSQLRepository(db)
	auditLog := audit.NewIOWriter(os.Stdout)
	taskSvc := task.NewService(repo, auditLog)

	ginrouter := gin.Default()

	router.Bootstrap(ginrouter, router.Dependencies{
		TaskService: taskSvc,
	})

	httpsrv := http.Server{
		Addr:    ":8010",
		Handler: ginrouter,
	}
	log.Fatal(httpsrv.ListenAndServe())
}
