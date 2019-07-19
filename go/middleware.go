package main

import (
  "fmt"

  jwt "github.com/dgrijalva/jwt-go"
  "github.com/gin-gonic/gin"
)

func respond(c *gin.Context, code int, message interface{}) {
  c.AbortWithStatusJSON(code, message)
}

func hasLogged() gin.HandlerFunc {
  return func(c *gin.Context) {
    _, err := c.Cookie("token")
    if err == nil {
      respond(c, 200, gin.H{
        "msg": "you are logged in",
      })
    }
  }
}

func authorized(cfg *Cfg) gin.HandlerFunc {
  return func (c *gin.Context) {
    tokenString, _ := c.Cookie("token")

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
      if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("signing method invalid: %v", token.Header["alg"])
      }
      return []byte(cfg.SignatureKey), nil
    })
    if err != nil {
      respond(c, 400, gin.H{
        "msg":err.Error(),
        "response":false,
      })
      return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
      respond(c, 400, gin.H{
        "msg":err.Error(),
        "response":false,
      })
      return
    }

    c.Set("UserData", claims["UserData"])
  }
}

func cors() gin.HandlerFunc {
  return func(c *gin.Context) {
    c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
    c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, "+
      "X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
  }
}
