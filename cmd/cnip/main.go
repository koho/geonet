package main

import (
	"fmt"
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
	ipv6List, err := readDat("https://ghproxy.com/https://github.com/v2fly/geoip/releases/latest/download/cn.dat")
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
	ipv4ListF := make([]string, 0)
	for _, ipv4 := range ipv4List {
		ipv4Str := strings.TrimSpace(ipv4)
		if ipv4Str == "" {
			continue
		}
		_, network, err := net.ParseCIDR(ipv4Str)
		if err != nil {
			log.Fatal(err)
		}
		record := mmdbtype.Map{}
		country := mmdbtype.Map{}
		record["country"] = country
		country["iso_code"] = mmdbtype.String("CN")
		if err = writer.Insert(network, record); err != nil {
			log.Fatal(err)
		}
		ipv4ListF = append(ipv4ListF, "\""+ipv4Str+"\"")
	}
	for _, geoip := range ipv6List.Entry {
		if geoip.CountryCode == "CN" {
			for _, v2rayCIDR := range geoip.Cidr {
				if ip := v2rayCIDR.GetIp(); len(ip) == 16 {
					ipStr := net.IP(ip).String() + "/" + fmt.Sprint(v2rayCIDR.GetPrefix())
					_, network, err := net.ParseCIDR(ipStr)
					if err != nil {
						log.Fatal(err)
					}
					record := mmdbtype.Map{}
					country := mmdbtype.Map{}
					record["country"] = country
					country["iso_code"] = mmdbtype.String("CN")
					if err = writer.Insert(network, record); err != nil {
						log.Fatal(err)
					}
				}
			}
			break
		}
	}
	script, err := plugin.Format(plugin.RosTemplate, plugin.RosScript{
		CIDRs:   strings.Join(ipv4ListF, ";"),
		Gateway: "[/ip route get [/ip route find routing-table=main dst-address=0.0.0.0/0] gateway]",
		Table:   "[/routing table get [/routing table find comment~\"cn\"] name]",
		Delay:   "50",
	})
	if err != nil {
		log.Fatal(err)
	}
	if err = os.MkdirAll("dist", 0666); err != nil {
		log.Fatal(err)
	}
	if err = os.WriteFile("dist/cn-ipv4.txt", []byte(strings.Join(ipv4List, "\n")), 0666); err != nil {
		log.Fatal(err)
	}
	if err = os.WriteFile("dist/cn-ipv4.rsc", []byte(script), 0666); err != nil {
		log.Fatal(err)
	}
	fh, err := os.Create("dist/cn.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	_, err = writer.WriteTo(fh)
	if err != nil {
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
