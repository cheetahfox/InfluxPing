package main

import (
	"fmt"
	"os"
	"log"
	"net"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Config struct {
	InfluxdbServer string
	InfluxDB string
	InfluxUsername string
	InfluxPassword string
	Ipv6Allowed bool
	PtpIpPing bool
	Hosts []string
}

func getConfigHosts(config Config) []net.IP {
	/*
	This function looks at the configuration we got from the Yml file 
	does validation and DNS resolution of the hosts in the config file. 
	Basically using the dns lookup to do validation
	*/
	var hosts []net.IP

	for i := range(config.Hosts) {
		resovledip, err := net.LookupIP(config.Hosts[i])
		if err != nil {
			fmt.Printf("Unable to resovle: %s\n", config.Hosts[i])
		}
		hosts = append(hosts, resovledip[0])
	}
	return hosts
}

func main() {
	var PingHost []net.IP

	if len(os.Args) <= 1 {
		log.Fatal("No Configuration file specified")
	}
	filename := os.Args[1]
	fmt.Println(filename)
	var config Config

	source, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(source, &config)
	if err != nil {
		log.Fatal(err)
	}

	PingHost = getConfigHosts(config)

	for i := range(PingHost) {
		fmt.Println(PingHost[i])
	}
}
