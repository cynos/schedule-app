package main

import (
  "net/url"
  "time"

  jwt "github.com/dgrijalva/jwt-go"
  "github.com/gin-gonic/gin"
  "github.com/jinzhu/gorm"
)

func logoutHandlers() gin.HandlerFunc {
  return func(c *gin.Context) {
    c.SetCookie("token", "", 0, "/", "", false, false)
    c.JSON(200, gin.H{"response":true})
  }
}

func loginHandlers(cfg *Cfg) gin.HandlerFunc {
  return func(c *gin.Context) {
    db := cfg.DBInstance

    var (
      user []logins
      data, _ = url.ParseQuery(c.PostForm("data"))
    )

    err := db.Where("username = ? AND password = ?", data["username"], data["password"]).Find(&user).Error
    if err != nil {
      c.JSON(500, gin.H{
        "msg":err.Error(),
        "response":false,
      })
      return
    }

    if len(user) > 0 {
      var claimsData = struct {
        jwt.StandardClaims
        UserData gin.H
      }{
        StandardClaims : jwt.StandardClaims{
          Issuer: "lumoshive",
          ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
        },
        UserData: gin.H{"id":user[0].ID, "name":user[0].Name},
      }

      token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsData)
      ss, _ := token.SignedString([]byte(cfg.SignatureKey))

      c.SetCookie("token", ss, 3600, "/", "", false, true)
      c.SetCookie("login", "true", 3600, "/", "", false, false)

      c.JSON(200, gin.H{
        "msg": "login success",
        "response":true,
      })
    } else {
      c.JSON(200, gin.H{
        "msg":"user not found",
        "response":false,
      })
    }
  }
}

func actionHandlers(cfg *Cfg) gin.HandlerFunc {
  return func(c *gin.Context) {
    db := cfg.DBInstance

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
      default: methodNotFound()
    }
  }
}

func create(c *gin.Context, db *gorm.DB) {
  var (
    UserData  = c.MustGet("UserData").(map[string]interface{})
    data      = jadwals{
                Kegiatan: c.PostForm("Kegiatan"),
                Jam:      c.PostForm("Jam"),
                Userid:   UserData["id"],
              }
  )

  db.Create(&data)
  c.JSON(200, gin.H{
    "response": true,
    "msg":   "berhasil insert data",
  })
}

func read(c *gin.Context, db *gorm.DB) {
  var (
    data []jadwals
    UserData = c.MustGet("UserData").(map[string]interface{})
  )

  db.Where("userid = ?", UserData["id"]).Find(&data)
  c.JSON(200, gin.H{
    "response": true,
    "data": data,
  })
}

func update(c *gin.Context, db *gorm.DB) {
  var (
    data []jadwals
    UserData = c.MustGet("UserData").(map[string]interface{})
    id       = c.PostForm("id")
    kegiatan = c.PostForm("Kegiatan")
    jam      = c.PostForm("Jam")
  )

  db.Where("id=? AND userid=?", id, UserData["id"]).Find(&data)

  if len(data) < 1 {
    c.JSON(200, gin.H{
      "response": false,
      "msg": "data not found",
    })
    return
  }

  err := db.Exec("UPDATE jadwals SET kegiatan=?, jam=?, updated_at=? WHERE id=? AND userid=?", kegiatan, jam, time.Now(), id, UserData["id"]).Error

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
  var (
    id       = c.PostForm("id")
    UserData = c.MustGet("UserData").(map[string]interface{})
  )

  err := db.Exec("DELETE FROM jadwals WHERE id=? AND userid=?", id, UserData["id"]).Error

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
