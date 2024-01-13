package main

import (
	"github.com/koho/geonet/plugin"
	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	ipv4List, err := readList("https://raw.githubusercontent.com/17mon/china_ip_list/master/china_ip_list.txt")
	if err != nil {
		log.Fatal(err)
	}
	ipv6List, err := readList("https://ispip.clang.cn/all_cn_ipv6.txt")
	if err != nil {
		log.Fatal(err)
	}
	writer, err := mmdbwriter.New(
		mmdbwriter.Options{
			DatabaseType: "GeoLite2-Country",
			RecordSize:   24,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	ipv4, err := insertIPList(writer, append(ipv4List, "1.12.12.12/32", "120.53.53.53/32"))
	if err != nil {
		log.Fatal(err)
	}
	if _, err = insertIPList(writer, ipv6List); err != nil {
		log.Fatal(err)
	}
	script, err := plugin.Format(plugin.RosTemplate, plugin.RosScript{
		CIDRs:   strings.Join(ipv4, ";"),
		Gateway: "[/ip route get [/ip route find routing-table=main dst-address=0.0.0.0/0] gateway]",
		Table:   "[/routing table get [/routing table find comment~\"cn\"] name]",
		Delay:   "50",
	})
	if err != nil {
		log.Fatal(err)
	}
	if err = os.MkdirAll("dist", 0750); err != nil {
		log.Fatal(err)
	}
	if err = os.WriteFile("dist/cn-ipv4.txt", []byte(strings.Join(ipv4, "\n")), 0660); err != nil {
		log.Fatal(err)
	}
	if err = os.WriteFile("dist/cn-ipv6.txt", []byte(strings.Join(ipv6List, "\n")), 0660); err != nil {
		log.Fatal(err)
	}
	if err = os.WriteFile("dist/cn-ipv4.rsc", []byte(script), 0660); err != nil {
		log.Fatal(err)
	}
	fh, err := os.Create("dist/cn.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	if _, err = writer.WriteTo(fh); err != nil {
		log.Fatal(err)
	}
}

func readList(source string) ([]string, error) {
	resp, err := http.Get(source)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, os.ErrInvalid
	}
	ipBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(ipBytes), "\n"), nil
}

func readDat(source string) (*router.GeoIPList, error) {
	resp, err := http.Get(source)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, os.ErrInvalid
	}
	geoipBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var geoipList router.GeoIPList
	if err = proto.Unmarshal(geoipBytes, &geoipList); err != nil {
		return nil, err
	}
	return &geoipList, nil
}

func insertIPList(writer *mmdbwriter.Tree, ipList []string) ([]string, error) {
	inserted := make([]string, 0)
	for _, ip := range ipList {
		ipStr := strings.TrimSpace(ip)
		if ipStr == "" {
			continue
		}
		_, network, err := net.ParseCIDR(ipStr)
		if err != nil {
			return inserted, err
		}
		record := mmdbtype.Map{
			"country": mmdbtype.Map{
				"iso_code": mmdbtype.String("CN"),
			},
		}
		if err = writer.Insert(network, record); err != nil {
			return inserted, err
		}
		inserted = append(inserted, "\""+ipStr+"\"")
	}
	return inserted, nil
}
