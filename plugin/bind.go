package bind

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/koho/geonet/lib"
	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
)

const (
	typeBindOut  = "bind"
	descBindOut  = "Convert data to bind format"
	bindTemplate = "zone \"%s\" {\n    include \"%s\";\n};"
)

func init() {
	lib.RegisterFormatter(typeBindOut, &bindOut{
		Description: descBindOut,
	})
}

type bindOut struct {
	Description string
}

func (b *bindOut) GetDescription() string {
	return b.Description
}

func (b *bindOut) FormatGeoIP(c *gin.Context, cidrs []*router.CIDR) error {
	return lib.ErrNotImplemented
}

func (b *bindOut) FormatGeoSite(c *gin.Context, domains []*router.Domain) error {
	inc := c.DefaultQuery("include", "/etc/bind/named.zones")
	domainMap := make(map[string]bool)
	for _, site := range domains {
		if !domainMap[site.Value] && site.Type != router.Domain_Regex {
			if _, err := c.Writer.WriteString(fmt.Sprintf(bindTemplate, site.Value, inc)); err != nil {
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
