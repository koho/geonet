package lib

import (
	"github.com/gin-gonic/gin"
	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
)

type Formatter interface {
	GetDescription() string
	FormatGeoIP(*gin.Context, []*router.CIDR) (string, error)
	FormatGeoSite(*gin.Context, []*router.Domain) (string, error)
}
