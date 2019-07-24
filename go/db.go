package main

import (
  "fmt"

  "github.com/golang-migrate/migrate/v4"
  "github.com/golang-migrate/migrate/v4/database/postgres"
  _ "github.com/golang-migrate/migrate/v4/source/file"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
  "github.com/lib/pq"
)

type contacts struct{
  gorm.Model
  Userid    interface{} `gorm:"type:integer"`
  Firstname string `form:"firstname" binding:"required" gorm:"type:varchar(50)"`
  Lastname  string `form:"lastname" binding:"required" gorm:"type:varchar(50)"`
  Phone     string `form:"phone" binding:"required" gorm:"type:varchar(50)"`
  Email     string `form:"email" binding:"required" gorm:"type:varchar(50)"`
}

type logins struct{
  gorm.Model
  Username string
  Password string
}

type jadwals struct{
  gorm.Model
  Confirmed   bool `gorm:"-"`
  Userid      interface{}
  Kegiatan    string
  Jam         string
  ContactList pq.Int64Array `gorm:"type:integer[]"`
}

type DBConfig struct {
  Host            string `env:"DB_HOST" required:"true"`
  Port						int 	 `env:"DB_PORT" required:"true"`
	Name            string `env:"DB_NAME" required:"true"`
	User            string `env:"DB_USER" required:"true"`
	Password        string `env:"DB_PASSWORD" required:"true"`
	ApplicationName string `env:"DB_APP_NAME"`
	Logging         bool   `env:"DB_LOGGING" required:"true" default:"false"`
  ConnectTimeout  int    `default:"30"`
	MaxOpenConn     int    `env:"DB_MAX_OPEN_CONN" default:"50"`
	MaxIdleConn     int    `env:"DB_MAX_IDLE_CONN" default:"10"`
}

func (d DBConfig) makeDBInfo() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s connect_timeout=%d application_name=%s",
		d.Host, d.Port, d.User, d.Name, d.Password, d.ConnectTimeout, d.ApplicationName)
}

func (d DBConfig) initDB() (*gorm.DB, error) {
	db, err := gorm.Open("postgres", d.makeDBInfo())
  if err != nil {
    return nil, err
  }

  defer db.AutoMigrate(&logins{}, &jadwals{}, &contacts{})

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

  return db, nil
}
