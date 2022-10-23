package plugin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/koho/geonet/lib"
	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
)

const (
	typeCoreDnsOut  = "coredns"
	descCoreDnsOut  = "Convert data to coredns format"
	coreDnsTemplate = "except"
)

func init() {
	lib.RegisterFormatter(typeCoreDnsOut, &coreDnsOut{
		Description: descCoreDnsOut,
	})
}

type coreDnsOut struct {
	Description string
}

func (cd *coreDnsOut) GetDescription() string {
	return cd.Description
}

func (cd *coreDnsOut) FormatGeoIP(c *gin.Context, cidrs []*router.CIDR) error {
	return lib.ErrNotImplemented
}

func (cd *coreDnsOut) FormatGeoSite(c *gin.Context, domains []*router.Domain) error {
	domainMap := make(map[string]bool)
	for _, site := range domains {
		if !domainMap[site.Value] && site.Type != router.Domain_Regex {
			if len(domainMap) == 0 {
				if _, err := c.Writer.WriteString(coreDnsTemplate); err != nil {
					return err
				}
			}
			if _, err := c.Writer.WriteString(fmt.Sprintf(" %s", site.Value)); err != nil {
				return err
			}
			domainMap[site.Value] = true
		}
	}
	return nil
}
