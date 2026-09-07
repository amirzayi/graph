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
