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

  //static routing
  r.Static("/static", "../static")

  // routing list
  r.POST("/auth", hasLogged(), loginHandler(&cfg))
  r.POST("/logout", logoutHandler())

  v1 := r.Group("/v1")
  v1.Use(authorized(&cfg))
  {
    v1.POST("/schedule", scheduleHandler(&cfg))
    v1.POST("/contacts", contactsHandler(&cfg))
    v1.POST("/timer", timerHandler(&cfg))
  }

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
