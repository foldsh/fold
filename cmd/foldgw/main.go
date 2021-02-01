package main

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/:service/*path", func(c *gin.Context) {
		handleRequest("GET", c)
	})
	r.PUT("/:service/*path", func(c *gin.Context) {
		handleRequest("PUT", c)
	})
	r.POST("/:service/*path", func(c *gin.Context) {
		handleRequest("POST", c)
	})
	r.DELETE("/:service/*path", func(c *gin.Context) {
		handleRequest("DELETE", c)
	})
	r.Run()
}

type proxyTarget struct {
	name string
	url  *url.URL
}

func getProxyTarget(c *gin.Context) proxyTarget {
	u := c.Request.URL
	u.Path = c.Param("path")
	return proxyTarget{name: c.Param("service"), url: u}
}

func handleRequest(method string, c *gin.Context) {
	pt := getProxyTarget(c)
	c.JSON(200, gin.H{
		"method":       method,
		"service_name": pt.name,
		"service_url":  pt.url.String(),
	})
}
