package space

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Starter interface {
	Start()
	Stop()
}

type App struct {
	configLoader *ConfigLoader
	starters     []Starter
	wg           sync.WaitGroup
}

func NewApp(configLoader *ConfigLoader, starters []Starter) *App {
	return &App{
		configLoader: configLoader,
		starters:     starters,
	}
}

func (app *App) Run() {
	app.configLoader.LoadConfig()
	for _, v := range app.starters {
		app.wg.Go(v.Start)
	}

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no params) by default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	for _, v := range app.starters {
		go v.Stop()
	}
	app.wg.Wait()
}
