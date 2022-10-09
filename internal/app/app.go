package app

import (
	"encoding/json"
	"fmt"
	_ "github.com/igilgyrg/betera-test/docs"
	"github.com/igilgyrg/betera-test/internal/config"
	"github.com/igilgyrg/betera-test/pkg/logging"
	"github.com/igilgyrg/betera-test/pkg/nasa"
	postgres "github.com/igilgyrg/betera-test/pkg/storage/posgres"
	s3Store "github.com/igilgyrg/betera-test/pkg/storage/s3"
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	timeout        = 10
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

type App struct {
	cfg        *config.Config
	postgresDB *gorm.DB
	nasaClient nasa.NASAClient
	s3Storage  s3Store.S3Storage
	echo       *echo.Echo
}

func New() *App {
	logging.Log().Info("init config file")
	cfg := config.NewConfig()
	if cfg == nil {
		log.Fatal("config file have not parsed")
	}

	logging.Log().Info("init s3Storage client")
	s3Storage, err := s3Store.NewS3Storage(s3Store.NewS3Config(cfg.AWSAccessKey, cfg.AWSSecretKey, cfg.AWSRegion, cfg.AWSBucketName))
	if err != nil {
		log.Fatalf("error of create S3 client %v", err)
	}

	logging.Log().Info("init nasa client")
	nasaClient := nasa.NewClient(cfg.NASAApiKey)

	logging.Log().Info("init postgres client")
	postgresDB, err := postgres.NewPostgresClient(postgres.NewPostgresConfig(cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUser, cfg.DBPassword))
	if err != nil {
		log.Fatalf("error of create postgres client %v", err)
	}

	return &App{
		cfg:        cfg,
		postgresDB: postgresDB,
		s3Storage:  s3Storage,
		nasaClient: nasaClient,
		echo:       echo.New(),
	}
}

func (a *App) Start() {
	logging.Log().Info("init health route")

	a.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	logging.Log().Info("init swagger router")
	a.echo.GET("/swagger/*", echoSwagger.WrapHandler)

	httpServer := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", a.cfg.Host, a.cfg.Port),
		ReadTimeout:    timeout * time.Second,
		WriteTimeout:   timeout * time.Second,
		MaxHeaderBytes: maxHeaderBytes,
	}

	logging.Log().Info("starting server...")
	go func() {
		if err := a.echo.StartServer(httpServer); err != nil {
			log.Fatalf("error start server, %w", err)
		}
	}()

	if err := a.MapHandlers(a.echo, a.cfg, ctxTimeout*time.Second); err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	logging.Log().Info("server have started")

	<-quit
}

// General godoc
// @Summary Health check
// @Description Health check endpoint
// @Tags general
// @Accept  json
// @Produce  json
// @Success 200
// @Router /health [get]
func health(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	responseBody := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	bytes, _ := json.Marshal(responseBody)

	res.Write(bytes)
}
