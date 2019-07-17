package main

import (
  "fmt"

  "github.com/golang-migrate/migrate/v4"
  "github.com/golang-migrate/migrate/v4/database/postgres"
  _ "github.com/golang-migrate/migrate/v4/source/file"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
)

type DBConfig struct {
  Host            string `env:"DB_HOST" required:"true"`
  Port						int 	 `env:"DB_PORT" required:"true"`
	Name            string `env:"DB_NAME" required:"true"`
	User            string `env:"DB_USER" required:"true"`
	Password        string `env:"DB_PASSWORD" required:"true"`
	ApplicationName string `env:"DB_APP_NAME"`
	Logging         bool   `env:"DB_LOGGING" required:"true" default:"false"`
  ConnectTimeout  int    `default:"30"`
	MaxOpenConn     int    `default:"50"`
	MaxIdleConn     int    `default:"10"`
}

func (d DBConfig) makeDBInfo() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s connect_timeout=%d application_name=%s",
		d.Host, d.Port, d.User, d.Name, d.Password, d.ConnectTimeout, d.ApplicationName)
}

func (d DBConfig) initDB() (*gorm.DB, error) {
	db, errConn := gorm.Open("postgres", d.makeDBInfo())

  // migration =============
  driver, _ := postgres.WithInstance(db.DB(), &postgres.Config{})

  m, err := migrate.NewWithDatabaseInstance(
    "file://migrations",
    "postgres", driver)
  if err != nil {
    return nil, err
  }

  err = m.Up()
  if err != nil && err != migrate.ErrNoChange {
    return nil, err
  }
  // migration =============


  db.DB().SetMaxIdleConns(d.MaxIdleConn)
  db.DB().SetMaxOpenConns(d.MaxOpenConn)

  if errConn != nil {
    return nil, errConn
  } else {
    return db, nil
  }
}
