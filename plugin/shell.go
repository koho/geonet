package plugin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/koho/geonet/lib"
	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"net"
	"strconv"
	"strings"
)

const (
	typeShellOut = "shell"
	descShellOut = "Convert data to shell format"
	ipv4Template = "ip route add %s %s"
	ipv6Template = "ip -6 route add %s %s"
)

func init() {
	lib.RegisterFormatter(typeShellOut, &shellOut{
		Description: descShellOut,
	})
}

type shellOut struct {
	Description string
}

func (s *shellOut) GetDescription() string {
	return s.Description
}

func (s *shellOut) FormatGeoIP(c *gin.Context, cidrs []*router.CIDR) (string, error) {
	var ret strings.Builder
	ipType, err := strconv.Atoi(c.DefaultQuery("type", "4"))
	if err != nil {
		return "", err
	}
	template := ipv4Template
	if ipType == 6 {
		template = ipv6Template
	}
	ipType -= 2
	if ipType < 0 {
		ipType = 0
	}
	for _, v2rayCIDR := range cidrs {
		if ip := v2rayCIDR.GetIp(); len(ip)>>ipType == 1 {
			ipStr := net.IP(ip).String() + "/" + fmt.Sprint(v2rayCIDR.GetPrefix())
			if _, err = ret.WriteString(fmt.Sprintf(template, ipStr, c.Query("opt"))); err != nil {
				return "", err
			}
			if _, err = ret.WriteString("\n"); err != nil {
				return "", err
			}
		}
	}
	return ret.String(), nil
}

func (s *shellOut) FormatGeoSite(c *gin.Context, domains []*router.Domain) (string, error) {
	return "", lib.ErrNotImplemented
}
