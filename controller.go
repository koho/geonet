package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/koho/geonet/lib"
	_ "github.com/koho/geonet/plugin"
	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
)

func generateGeoip(c *gin.Context) {
	formatter, err := lib.GetFormatter(c.DefaultQuery("format", "text"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	country := c.DefaultQuery("country", "CN")
	source := "https://ghproxy.com/https://github.com/v2fly/geoip/releases/latest/download/geoip.dat"
	if country == "CN" {
		source = "https://ghproxy.com/https://github.com/v2fly/geoip/releases/latest/download/cn.dat"
	}
	resp, err := http.Get(source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("http status %d", resp.StatusCode)})
		return
	}
	geoipBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var geoipList router.GeoIPList
	if err = proto.Unmarshal(geoipBytes, &geoipList); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for _, geoip := range geoipList.Entry {
		if geoip.CountryCode == country {
			if ret, err := formatter.FormatGeoIP(c, geoip.Cidr); err == nil {
				c.String(http.StatusOK, ret)
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
	}
}

func generateGeosite(c *gin.Context) {
	formatter, err := lib.GetFormatter(c.DefaultQuery("format", "text"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := http.Get("https://ghproxy.com/https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("http status %d", resp.StatusCode)})
		return
	}
	geoSiteBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var geoSiteList router.GeoSiteList
	if err = proto.Unmarshal(geoSiteBytes, &geoSiteList); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for _, geoSite := range geoSiteList.Entry {
		if geoSite.CountryCode == c.DefaultQuery("country", "CN") {
			if ret, err := formatter.FormatGeoSite(c, geoSite.Domain); err == nil {
				c.String(http.StatusOK, ret)
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
	}
}
