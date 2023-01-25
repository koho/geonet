package plugin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/koho/geonet/lib"
	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"net"
	"net/http"
	"strconv"
)

const (
	typeMMDBOut = "mmdb"
	descMMDBOut = "Convert data to mmdb format"
)

func init() {
	lib.RegisterFormatter(typeMMDBOut, &mmdbOut{
		Description: descMMDBOut,
	})
}

type mmdbOut struct {
	Description string
}

func (m *mmdbOut) GetDescription() string {
	return m.Description
}

func (m *mmdbOut) FormatGeoIP(c *gin.Context, cidrs []*router.CIDR, countryCode string) error {
	writer, err := mmdbwriter.New(
		mmdbwriter.Options{
			DatabaseType: "GeoLite2-Country",
			RecordSize:   24,
		},
	)
	if err != nil {
		return err
	}
	ipType, err := strconv.Atoi(c.DefaultQuery("type", "0"))
	if err != nil {
		return err
	}
	ipType -= 2
	if ipType < 0 {
		ipType = 0
	}
	for _, v2rayCIDR := range cidrs {
		if ip := v2rayCIDR.GetIp(); ipType == 0 || len(ip)>>ipType == 1 {
			ipStr := net.IP(ip).String() + "/" + fmt.Sprint(v2rayCIDR.GetPrefix())
			_, network, err := net.ParseCIDR(ipStr)
			if err != nil {
				return err
			}
			record := mmdbtype.Map{}
			country := mmdbtype.Map{}
			record["country"] = country
			country["iso_code"] = mmdbtype.String(countryCode)
			if err = writer.Insert(network, record); err != nil {
				return err
			}
		}
	}
	c.Status(http.StatusOK)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.mmdb", countryCode))
	_, err = writer.WriteTo(c.Writer)
	return err
}

func (m *mmdbOut) FormatGeoSite(c *gin.Context, domains []*router.Domain, countryCode string) error {
	return lib.ErrNotImplemented
}
