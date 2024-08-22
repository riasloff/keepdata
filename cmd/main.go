package main

import (
	"context"
	"fmt"
	"keepdata/internal/alert"
	"keepdata/internal/config"
	"keepdata/internal/logger"
	"keepdata/internal/monitor"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var conf config.Conf

func main() {

	conf.Read()
	logger.Log.Info("чтение конфигурационного файла прошло успешно")

	logger.Setup(&conf)
	defer logger.Log.Sync()
	logger.Log.Info("инициализация системы логов прошла успешно")

	if err := run(); err != nil {
		logger.Log.Fatal("неизвестная ошибка", zap.Error(err))
	}
}

func run() error {

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/", homeHandler)
	router.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })
	router.GET("/about", aboutHandler)
	router.GET("/stream", streamHandler)
	router.GET("/filter", filterHandler)
	router.LoadHTMLGlob("templates/*")

	srv := &http.Server{
		Addr:    ":" + conf.Listen_port,
		Handler: router.Handler(),
	}

	errChan := make(chan error, 1)
	go func() { // веб сервер
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("ошибка работы веб сервера: %w", err)
		}
	}()
	logger.Log.Info("сервер запущен успешно", zap.String("port", conf.Listen_port))

	go func() { // мониторинг файловой системы
		if err := monitor.Monitor(conf.Track_files, conf.Track_dirs); err != nil {
			errChan <- fmt.Errorf("ошибка работы мониторинга файловой системы: %w", err)
		}
	}()

	go func() { // генератор сообщений о статусе
		for {
			logger.Log.Debug("сообщение о статусе каждые 30 секунд", zap.String("status", "ok"))
			time.Sleep(30 * time.Second)
		}
	}()

	alert.Alert()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-stopChan:
		logger.Log.Info("изящное завершение работы сервера", zap.String("signal", sig.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(ctx)
	case err := <-errChan:
		return err
	}
}
