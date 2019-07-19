package main

import (
  "github.com/jinzhu/gorm"
  "github.com/jinzhu/configor"
  "go.uber.org/zap"
)

type Cfg struct {
  DBConfig
  DBInstance    *gorm.DB
  SignatureKey  string
  Timezone      string     `env:"TIMEZONE" required:"true"`
  HttpAddress   string     `env:"HTTP_ADDRESS"`
  Debug         bool       `required:"true" default:"false"`
  Log           *zap.Logger
}

func (c *Cfg) init() {
  // initialize zap logger
  c.Log, _ = zap.NewProduction()
  defer c.Log.Sync()

  // initialize configor
  err := configor.New(&configor.Config{
    Environment:          "production",
    ErrorOnUnmatchedKeys: true,
    Debug:                true,
  }).Load(c, "config.yaml")

  if err != nil {
    c.Log.Error("unmatched configuration keys , ", zap.Error(err))
  }

  // initialize database
  c.DBInstance, err = c.initDB()
  if err != nil {
    panic(err.Error())
  }
}
