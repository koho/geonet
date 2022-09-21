package bind

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/koho/geonet/lib"
	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
)

const (
	typeDnsmasqOut  = "dnsmasq"
	descDnsmasqOut  = "Convert data to dnsmasq format"
	dnsmasqTemplate = "server=/%s/%s"
)

func init() {
	lib.RegisterFormatter(typeDnsmasqOut, &dnsmasqOut{
		Description: descDnsmasqOut,
	})
}

type dnsmasqOut struct {
	Description string
}

func (d *dnsmasqOut) GetDescription() string {
	return d.Description
}

func (d *dnsmasqOut) FormatGeoIP(c *gin.Context, cidrs []*router.CIDR) error {
	return lib.ErrNotImplemented
}

func (d *dnsmasqOut) FormatGeoSite(c *gin.Context, domains []*router.Domain) error {
	dns := c.DefaultQuery("dns", "114.114.114.114")
	domainMap := make(map[string]bool)
	for _, site := range domains {
		if !domainMap[site.Value] && site.Type != router.Domain_Regex {
			if _, err := c.Writer.WriteString(fmt.Sprintf(dnsmasqTemplate, site.Value, dns)); err != nil {
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
