package plugin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/koho/geonet/lib"
	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"net"
	"net/http"
	"strconv"
	"strings"
)

const (
	typeRosOut  = "ros"
	descRosOut  = "Convert data to RouterOS format"
	rosTemplate = "add distance=1 dst-address=%s gateway=%s routing-table=%s"
)

func init() {
	lib.RegisterFormatter(typeRosOut, &rosOut{
		Description: descRosOut,
	})
}

type rosOut struct {
	Description string
}

func (r *rosOut) GetDescription() string {
	return r.Description
}

func (r *rosOut) FormatGeoIP(c *gin.Context, cidrs []*router.CIDR, countryCode string) error {
	var ret strings.Builder
	ipType, err := strconv.Atoi(c.DefaultQuery("type", "4"))
	if err != nil {
		return err
	}
	ipType -= 2
	if ipType < 0 {
		ipType = 0
	}
	if _, err = ret.WriteString("/ip route\n"); err != nil {
		return err
	}
	for _, v2rayCIDR := range cidrs {
		if ip := v2rayCIDR.GetIp(); len(ip)>>ipType == 1 {
			ipStr := net.IP(ip).String() + "/" + fmt.Sprint(v2rayCIDR.GetPrefix())
			if _, err = ret.WriteString(fmt.Sprintf(rosTemplate, ipStr, c.Query("gw"), c.DefaultQuery("table", "main"))); err != nil {
				return err
			}
			if _, err = ret.WriteString("\n"); err != nil {
				return err
			}
		}
	}
	c.String(http.StatusOK, ret.String())
	return nil
}

func (r *rosOut) FormatGeoSite(c *gin.Context, domains []*router.Domain, countryCode string) error {
	return lib.ErrNotImplemented
}
