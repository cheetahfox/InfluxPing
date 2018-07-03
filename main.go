package main

import (
	"fmt"
	"os"
	"log"
	"net"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/sparrc/go-ping"
	"time"
)

type Config struct {
	InfluxdbServer string
	InfluxDB string
	InfluxUsername string
	InfluxPassword string
	Ipv6Allowed bool
	PtpIpPing bool
	ProbeCount int
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
		} else {
			hosts = append(hosts, resovledip[0])
		}
	}
	return hosts
}

func main() {
	var PingHost []net.IP

	var interval time.Duration = 100 * time.Millisecond
	var timeout time.Duration = 10000 * time.Millisecond

	if len(os.Args) <= 1 {
		log.Fatal("No Configuration file specified")
	}
	filename := os.Args[1]
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
		pinger, err := ping.NewPinger(PingHost[i].String())
		if err != nil {
			log.Fatal(err)
		}

		pinger.OnRecv = func(pkt *ping.Packet) {
			fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n", pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
		}

		pinger.OnFinish = func(stats *ping.Statistics) {
			fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
			fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n", stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
			fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n", stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
		}

		pinger.Count    = config.ProbeCount
		pinger.Interval = interval
		pinger.Timeout  = timeout

		pinger.SetPrivileged(true)

		pinger.Run()
	}
}
