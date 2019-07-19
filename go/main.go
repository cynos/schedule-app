package main

import (
  "context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

  "github.com/gin-gonic/gin"
)

func main() {
  cfg := Cfg{}
  cfg.init()

  gin.SetMode(gin.ReleaseMode)
  r := gin.New()

  if cfg.Debug {
    r.Use(gin.Logger())
    r.Use(gin.Recovery())
  }

  r.Use(cors())

  // end-point list
  r.POST("/auth", hasLogged(), loginHandlers(&cfg))
  r.POST("/action-schedule", authorized(&cfg), actionHandlers(&cfg))
  r.POST("/logout", logoutHandlers())

  //static routing
  r.Static("/static", "../static")

  srv := &http.Server{
    Addr: cfg.HttpAddress,
    Handler: r,
  }

  go func() {
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
      panic(err.Error())
    }
  }()

  // Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
  quit := make(chan os.Signal, 2)

  signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
  <-quit
  cfg.Log.Info("Shutdown Server ..")

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  if err := srv.Shutdown(ctx); err != nil {
    cfg.Log.Info("can't shutdown server")
  }
  cfg.Log.Info("Server exiting")
}
