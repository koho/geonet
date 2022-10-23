package plugin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/koho/geonet/lib"
	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"net"
	"strconv"
)

const (
	typeTextOut = "text"
	descTextOut = "Convert data to plaintext format"
)

func init() {
	lib.RegisterFormatter(typeTextOut, &textOut{
		Description: descTextOut,
	})
}

type textOut struct {
	Description string
}

func (t *textOut) GetDescription() string {
	return t.Description
}

func (t *textOut) FormatGeoIP(c *gin.Context, cidrs []*router.CIDR) error {
	ipType, err := strconv.Atoi(c.DefaultQuery("type", "4"))
	if err != nil {
		return err
	}
	ipType -= 2
	if ipType < 0 {
		ipType = 0
	}
	for _, v2rayCIDR := range cidrs {
		if ip := v2rayCIDR.GetIp(); len(ip)>>ipType == 1 {
			ipStr := net.IP(ip).String() + "/" + fmt.Sprint(v2rayCIDR.GetPrefix())
			if _, err = c.Writer.WriteString(ipStr); err != nil {
				return err
			}
			if _, err = c.Writer.WriteString("\n"); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *textOut) FormatGeoSite(c *gin.Context, domains []*router.Domain) error {
	domainMap := make(map[string]bool)
	for _, site := range domains {
		if !domainMap[site.Value] {
			if _, err := c.Writer.WriteString(site.Value); err != nil {
				return err
			}
			if _, err := c.Writer.WriteString("\n"); err != nil {
				return err
			}
			domainMap[site.Value] = true
		}
	}
	return nil
}
