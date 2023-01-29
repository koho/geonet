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
	"text/template"
)

const (
	typeRosOut  = "ros"
	descRosOut  = "Convert data to RouterOS format"
	RosTemplate = `:local table [[.Table]];
:local gateway [[.Gateway]];
:local cidrs {[[.CIDRs]]};
/log info "syncing routing table: $table";
:foreach cidr in=$cidrs do={
    :if ([:len [/ip route find dst-address=$cidr gateway=$gateway routing-table=$table]] = 0) do={/ip route add distance=1 dst-address=$cidr gateway=$gateway routing-table=$table;}
    :delay [[.Delay]]ms;
}
:foreach i in=[/ip route find gateway=$gateway routing-table=$table (comment).""=""] do={
    :local p [:find $cidrs [/ip route get $i dst-address]];
    :if ([:type $p]="nil") do={
        /ip route remove $i;
    }
    :delay [[.Delay]]ms;
}
/log info "updated routing table: $table";
`
)

func init() {
	lib.RegisterFormatter(typeRosOut, &rosOut{
		Description: descRosOut,
	})
}

type rosOut struct {
	Description string
}

type RosScript struct {
	CIDRs   string
	Gateway string
	Table   string
	Delay   string
}

func (r *rosOut) GetDescription() string {
	return r.Description
}

func (r *rosOut) FormatGeoIP(c *gin.Context, cidrs []*router.CIDR, countryCode string) error {
	gw := c.Query("gw")
	if gw == "" {
		return lib.ErrInvalidParameter
	}
	ipType, err := strconv.Atoi(c.DefaultQuery("type", "4"))
	if err != nil {
		return err
	}
	ipType -= 2
	if ipType < 0 {
		ipType = 0
	}
	ipList := make([]string, 0)
	for _, v2rayCIDR := range cidrs {
		if ip := v2rayCIDR.GetIp(); len(ip)>>ipType == 1 {
			ipStr := net.IP(ip).String() + "/" + fmt.Sprint(v2rayCIDR.GetPrefix())
			ipList = append(ipList, "\""+ipStr+"\"")
		}
	}
	script, err := Format(RosTemplate, RosScript{
		CIDRs:   strings.Join(ipList, ";"),
		Gateway: "\"" + gw + "\"",
		Table:   "\"" + c.DefaultQuery("table", "main") + "\"",
		Delay:   c.DefaultQuery("delay", "50"),
	})
	if err != nil {
		return err
	}
	c.String(http.StatusOK, script)
	return nil
}

func (r *rosOut) FormatGeoSite(c *gin.Context, domains []*router.Domain, countryCode string) error {
	return lib.ErrNotImplemented
}

func Format(s string, v interface{}) (string, error) {
	t, b := new(template.Template), new(strings.Builder)
	t.Delims("[[", "]]")
	tp, err := t.Parse(s)
	if err != nil {
		return "", err
	}
	if err = tp.Execute(b, v); err != nil {
		return "", err
	}
	return b.String(), nil
}
