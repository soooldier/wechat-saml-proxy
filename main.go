package main

import (
	"io/ioutil"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"wechat-saml-proxy/service"
)

func main() {
	r := gin.Default()
	store := memstore.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))
	r.GET("/", func(c *gin.Context) {
		b, err := ioutil.ReadFile("./index.html")
		if err != nil {
			c.HTML(500, "Internal Error", nil)
			c.Abort()
		}
		c.HTML(200, string(b), nil)
	})
	r.GET("/MP_verify_32iWga2EVle6QTQm.txt", func(c *gin.Context) {
		c.HTML(200, string("32iWga2EVle6QTQm"), nil)
	})
	r.GET("/api/callback", service.LoginHandler)
	r.GET("/api/saml", func(c *gin.Context) {
		openid4interface := sessions.Default(c).Get("openid4wechat")
		openid4string := openid4interface.(string)
		c.HTML(200, openid4string, nil)
	})
	r.Run(":80")
}
