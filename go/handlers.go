package main

import (
  "time"

  "github.com/gin-gonic/gin"
  "github.com/jinzhu/gorm"
)

type jadwals struct{
  gorm.Model
  Kegiatan  string
  Jam       string
}

func actionHandlers(cfg *Cfg) gin.HandlerFunc {
  return func(c *gin.Context) {
    db := cfg.DBInstance
    db.AutoMigrate(&jadwals{})

    var methodNotFound = func() {
      c.JSON(400, gin.H{
        "response": false,
        "msg": "method not found",
      })
    }

    method := c.PostForm("method")
    if method == "" {
      methodNotFound()
      return
    }

    switch method {
      case "create": create(c, db)
      case "read"  : read(c, db)
      case "update": update(c, db)
      case "delete": delete(c, db)
      default:
        methodNotFound()
    }
  }
}

func create(c *gin.Context, db *gorm.DB) {
  var data = jadwals{
    Kegiatan: c.PostForm("Kegiatan"),
    Jam: c.PostForm("Jam"),
  }
  db.Create(&data)

  c.JSON(200, gin.H{
    "response": true,
    "msg":   "berhasil insert data",
  })
}

func read(c *gin.Context, db *gorm.DB) {
  var data []jadwals
  db.Find(&data)

  c.JSON(200, gin.H{
    "response": true,
    "data": data,
  })
}

func update(c *gin.Context, db *gorm.DB) {
  id       := c.PostForm("id")
  kegiatan := c.PostForm("Kegiatan")
  jam      := c.PostForm("Jam")

  var data jadwals

  err := db.First(&data, id).Error
  if err != nil {
    c.JSON(200, gin.H{
      "response": false,
      "msg": "data not found",
    })
    return
  }

  err = db.Exec("UPDATE jadwals SET kegiatan=?, jam=?, updated_at=? WHERE id=?", kegiatan, jam, time.Now(), id).Error

  if err != nil {
    c.JSON(200, gin.H{
      "response": false,
      "msg": "update failed | "+err.Error(),
    })
  } else {
    c.JSON(200, gin.H{
      "response": true,
      "msg": "successfully update data",
    })
  }
}

func delete(c *gin.Context, db *gorm.DB) {
  id := c.PostForm("id")
  err := db.Exec("DELETE FROM jadwals WHERE id=?", id).Error
  if err != nil {
    c.JSON(200, gin.H{
      "response": false,
      "msg": "delete failed | "+err.Error(),
    })
  } else {
    c.JSON(200, gin.H{
      "response": true,
      "msg": "successfully delete data",
    })
  }
}
