package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/iosmanthus/geomatch"
	v2ray "v2ray.com/core/app/router"
)

var (
	geosite  = flag.String("geosite", "", "geosite.dat file location")
	name     = flag.String("name", "", "name of pre-defined domain group")
	upstream = flag.String("upstream", "114.114.114.114", "dns upstream")
)

func ExtractDomainList(location string, name string) ([]string, error) {
	data, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, err
	}

	geoList := new(v2ray.GeoSiteList)
	if err = proto.Unmarshal(data, geoList); err != nil {
		return nil, err
	}

	var list []*v2ray.Domain
	if list, err = geomatch.ExtractDomainList(name, geoList); err != nil {
		return nil, err
	}

	domains := make([]string, 0, len(list))
	for _, domain := range list {
		if domain.Type == v2ray.Domain_Domain {
			domains = append(domains, domain.Value)
		}
	}

	return domains, nil
}

func main() {
	flag.Parse()
	var domains []string
	domains, err := ExtractDomainList(*geosite, *name)

	if err != nil {
		log.Fatal(err)
	}

	if len(domains) == 0 {
		log.Fatal("domain list is empty")
	}

	rule := "["
	for _, domain := range domains {
		rule += "/" + domain
	}
	rule += "/]" + *upstream
	fmt.Println(rule)
}
