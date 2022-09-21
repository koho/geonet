package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/geoip", generateGeoip)
	r.GET("/geosite", generateGeosite)
	r.Run()
}
