package plugin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/koho/geonet/lib"
	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"strings"
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

func (cd *coreDnsOut) FormatGeoIP(c *gin.Context, cidrs []*router.CIDR) (string, error) {
	return "", lib.ErrNotImplemented
}

func (cd *coreDnsOut) FormatGeoSite(c *gin.Context, domains []*router.Domain) (string, error) {
	var ret strings.Builder
	domainMap := make(map[string]bool)
	for _, site := range domains {
		if !domainMap[site.Value] && site.Type != router.Domain_Regex {
			if len(domainMap) == 0 {
				if _, err := ret.WriteString(coreDnsTemplate); err != nil {
					return "", err
				}
			}
			if _, err := ret.WriteString(fmt.Sprintf(" %s", site.Value)); err != nil {
				return "", err
			}
			domainMap[site.Value] = true
		}
	}
	return ret.String(), nil
}
