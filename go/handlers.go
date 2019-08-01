package main

import (
  "encoding/json"
  "fmt"
  "log"
  "net/smtp"
  "net/url"
  "strconv"
  "strings"
  "time"


  jwt "github.com/dgrijalva/jwt-go"
  "github.com/gin-gonic/gin"
  "github.com/jinzhu/gorm"
)

func responseError(c *gin.Context, code int, msg interface{}) {
  c.AbortWithStatusJSON(code, gin.H{"response":false, "msg":msg})
}

func methodNotFound(c *gin.Context) {
  c.JSON(400, gin.H{
    "response": false,
    "msg": "method not found",
  })
}

func logoutHandler() gin.HandlerFunc {
  return func(c *gin.Context) {
    c.SetCookie("token", "", 0, "/", "", false, false)
    c.JSON(200, gin.H{"response":true})
  }
}

func loginHandler(cfg *Cfg) gin.HandlerFunc {
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
        UserData: gin.H{"id":user[0].ID},
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

func scheduleHandler(cfg *Cfg) gin.HandlerFunc {
  return func(c *gin.Context) {
    var (
      jadwal, val jadwals
      db       = cfg.DBInstance
      UserData = c.MustGet("UserData").(map[string]interface{})
      idjadwal = c.DefaultPostForm("idjadwal", "0")
      method   = c.PostForm("method")
    )

    create := func(){
      var data = jadwals{
        Kegiatan: c.PostForm("Kegiatan"),
        Jam:      c.PostForm("Jam"),
        Userid:   UserData["id"],
      }

      db.Create(&data)
      c.JSON(200, gin.H{
        "response": true,
        "msg":   "successfully insert data",
      })
    }

    read := func(){
      var data []jadwals

      db.Where("userid = ?", UserData["id"]).Find(&data)
      c.JSON(200, gin.H{
        "response": true,
        "data": data,
      })
    }

    update := func(){
      var (
        data []jadwals
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
        c.JSON(500, gin.H{
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

    delete := func(){
      var id = c.PostForm("id")

      err := db.Exec("DELETE FROM jadwals WHERE id=? AND userid=?", id, UserData["id"]).Error

      if err != nil {
        c.JSON(500, gin.H{
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

    getContact := func(){
      err := db.Find(&val, idjadwal).Error
      if err != nil {
        responseError(c, 400, err.Error())
        return
      } else {
        var contact []contacts
        err = db.Where([]int64(val.ContactList)).Find(&contact).Error
        if err != nil {
          responseError(c, 400, err.Error())
        } else {
          c.JSON(200, gin.H{"response": true, "data": contact})
        }
      }
    }

    getContactOptions := func(){
      err := db.Find(&val, idjadwal).Error
      if err != nil {
        responseError(c, 400, err.Error())
        return
      } else {
        var contact []contacts
        err = db.Not([]int64(val.ContactList)).Find(&contact).Error
        if err != nil {
          responseError(c, 400, err.Error())
        } else {
          c.JSON(200, gin.H{"response": true, "data": contact})
        }
      }
    }

    addContact := func(){
      var idcontact []string
      _ = json.Unmarshal([]byte(c.PostForm("idcontact")), &idcontact)

      if len(idcontact) < 1 {
        responseError(c, 400, "invalid request")
        return
      }

      err := db.Find(&val, idjadwal).Error
      if err != nil {
        responseError(c, 400, err.Error())
        return
      } else {
        for _, num := range idcontact {
          id, _ := strconv.ParseInt(num, 10, 64)
          val.ContactList = append(val.ContactList, id)
        }
        val.Confirmed = true
        db.Model(&jadwal).Updates(val)
        c.JSON(200, gin.H{"response":true, "msg":"successfully add contact"})
      }
    }

    removeContact := func(){
      idcontact, _ := strconv.ParseInt(c.DefaultPostForm("idcontact", "0"), 10, 64)
      if idcontact == 0 {
        responseError(c, 400, "invalid request")
      }

      err := db.Exec("UPDATE jadwals SET contact_list = array_remove(contact_list, ?) WHERE id=?", idcontact, idjadwal).Error
      if err != nil {
        responseError(c, 400, err.Error())
      } else {
        c.JSON(200, gin.H{"response":true, "msg":"successfully remove contact"})
      }
    }

    switch method {
      case "create": create()
      case "read"  : read()
      case "update": update()
      case "delete": delete()

      case "get-contact": getContact()
      case "get-contact-options": getContactOptions()
      case "add-contact": addContact()
      case "remove-contact": removeContact()
      default: methodNotFound(c)
    }
  }
}

func contactsHandler(cfg *Cfg) gin.HandlerFunc {
  return func(c *gin.Context) {
    var (
      contact, val contacts
      db        = cfg.DBInstance
      userData  = c.MustGet("UserData").(map[string]interface{})
      idcontact = c.DefaultPostForm("idcontact", "0")
      method    = c.PostForm("method")
    )

    create := func(){
      if err := c.ShouldBind(&contact); err != nil {
        responseError(c, 400, err.Error())
        return
      }

      contact.Userid = userData["id"]

      err := db.Create(&contact).Error
      if err != nil {
        responseError(c, 400, err.Error())
      } else {
        c.JSON(200, gin.H{"response": true, "msg": "successfully insert data"})
      }
    }

    read := func(){
      var contact []contacts

      err := db.Where("userid = ?", userData["id"]).Find(&contact).Error
      if err != nil {
        responseError(c, 400, err.Error())
      } else {
        c.JSON(200, gin.H{"response": true, "data": contact})
      }
    }

    update := func(){
      if err := c.ShouldBind(&val); err != nil {
        responseError(c, 400, err.Error())
        return
      }

      err := db.Find(&contact, idcontact).Error
      if err != nil {
        responseError(c, 400, err.Error())
        return
      } else {
        err = db.Model(&contact).Updates(val).Error
        if err != nil {
          responseError(c, 400, err.Error())
          return
        }
        c.JSON(200, gin.H{"response": true, "msg": "successfully update data"})
      }
    }

    delete := func(){
      err := db.Find(&contact, idcontact).Error
      if err != nil {
        responseError(c, 400, err.Error())
        return
      } else {
        db.Delete(&contact)
        c.JSON(200, gin.H{"response": true, "msg": "successfully delete data"})
      }
    }

    switch method {
      case "create": create()
      case "read"  : read()
      case "update": update()
      case "delete": delete()
      default: methodNotFound(c)
    }
  }
}

func timerHandler(cfg *Cfg) gin.HandlerFunc {
  return func(c *gin.Context) {
    var (
      db       = cfg.DBInstance
      UserData = c.MustGet("UserData").(map[string]interface{})
      method   = c.PostForm("method")
    )

    get := func() {
      var data []timer
      db.Where("userid = ?", UserData["id"]).Find(&data)
      c.JSON(200, gin.H{
        "response": true,
        "data": data,
      })
    }

    add := func() {
      var data = timer{
        Userid: UserData["id"],
        Time:   c.DefaultPostForm("time", ""),
        Title:  c.DefaultPostForm("title", ""),
      }

      db.Create(&data)
      c.JSON(200, gin.H{
        "response": true,
        "msg":   "successfully insert data",
      })
    }

    switch method {
      case "add": add()
      case "get": get()
      default: methodNotFound(c)
    }
  }
}

// hooks gorm ================
func (j *jadwals) AfterUpdate(tx *gorm.DB) (err error) {
  if j.Confirmed {

    // get mail list
    var (
      contact []contacts
      email []string
    )
    tx.Where([]int64(j.ContactList)).Find(&contact)
    for _, val := range contact {
      email = append(email, val.Email)
    }

    // send mail
    const CONFIG_SMTP_HOST = "smtp.gmail.com"
    const CONFIG_SMTP_PORT = 587
    const CONFIG_EMAIL     = "cynomous@gmail.com"
    const CONFIG_PASSWORD  = ""

    sendMail := func(to []string, cc []string, subject, message string) error {
      body := "From: " + CONFIG_EMAIL + "\n" +
              "To: " + strings.Join(to, ",") + "\n" +
              "Cc: " + strings.Join(cc, ",") + "\n" +
              "Subject: " + subject + "\n\n" +
              message

      auth := smtp.PlainAuth("", CONFIG_EMAIL, CONFIG_PASSWORD, CONFIG_SMTP_HOST)
      smtpAddr := fmt.Sprintf("%s:%d", CONFIG_SMTP_HOST, CONFIG_SMTP_PORT)

      err := smtp.SendMail(smtpAddr, auth, CONFIG_EMAIL, append(to, cc...), []byte(body))
      if err != nil {
          return err
      }
      return nil
    }

    to := email
    cc := []string{}
    subject := "Notifikasi Jadwal"
    message := "Hello, anda telah ditambahkan kedalam jadwal, berikut rincian jadwal anda \n" +
                "Kegiatan : " + j.Kegiatan + "\n" +
                "Jam : " + j.Jam

    err := sendMail(to, cc, subject, message)
    if err != nil {
      log.Fatal(err.Error())
    }
    log.Println("Mail sent!")
  }
  return
}
