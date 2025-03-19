package main

import (
	"awesomeProject/config"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
)

type App struct {
	conf    *config.Configuration
	logger  *zap.Logger
	httpSrv *http.Server
}

func newHttpServer(
	conf *config.Configuration,
	router *gin.Engine,
) *http.Server {
	return &http.Server{
		Addr:    ":" + conf.App.Port,
		Handler: router,
	}
}

func newApp(
	conf *config.Configuration,
	logger *zap.Logger,
	httpSrv *http.Server,
) *App {
	return &App{
		conf:    conf,
		logger:  logger,
		httpSrv: httpSrv,
	}
}

func (a *App) Start() error {
	// 启动 http server
	go func() {
		log.Printf("http server started")
		a.logger.Info("http server started")
		if err := a.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return nil
}

func (a *App) Stop(ctx context.Context) (err error) {
	log.Printf("http server has been stop")
	a.logger.Info("http server has been stop")
	if err = a.httpSrv.Shutdown(ctx); err != nil {
		return
	}
	return
}
