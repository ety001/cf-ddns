package main

import (
	"net"
	"os"

	"github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"crypto/tls"
)

var log = logrus.New()
var Version string


func main() {
	var (
		app = kingpin.New("cf-ddns", "Cloudflare DynDNS Updater").Version(Version)

		ipAddress   = app.Flag("ip-address", "Skip resolving external IP and use provided IP").String()
		ipv6Address = app.Flag("ipv6-address", "Skip resolving external IPv6 and use provided IPv6").String()
		noVerify    = app.Flag("no-verify", "Don't verify ssl certificates").Bool()

		cfEmail  = app.Flag("cf-email", "Cloudflare Email").Required().String()
		cfApiKey = app.Flag("cf-api-key", "Cloudflare API key").Required().String()
		cfZoneId = app.Flag("cf-zone-id", "Cloudflare Zone ID").Required().String()

		hostnames = app.Arg("hostnames", "Hostnames to update").Required().Strings()
	)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	var ip IPService
	var dns *CFDNSUpdater
	var err error

	if *ipAddress != "" || *ipv6Address != "" {
		ip = &FakeIPService{
			fakeIp:   net.ParseIP(*ipAddress),
			fakeIPv6: net.ParseIP(*ipv6Address),
		}
	} else {
		httpClient := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: *noVerify},
			},
		}
		ip = &IpifyIPService{HttpClient: httpClient}
	}

	if dns, err = NewCFDNSUpdater(*cfZoneId, *cfApiKey, *cfEmail, log.WithField("component", "cf-dns-updater")); err != nil {
		log.Panic(err)
	}

	// Update A record (IPv4)
	res, err := ip.GetExternalIP()
	if err != nil {
		log.Panic(err)
	}

	for _, hostname := range *hostnames {
		err := dns.UpdateRecordA(hostname, res)
		if err != nil {
			log.Panic(err)
		}
	}

	// Update AAAA record (IPv6)
	res6, err := ip.GetExternalIPv6()
	if err != nil {
		log.Warnf("Failed to get IPv6 address: %v (skipping IPv6 update)", err)
	} else {
		for _, hostname := range *hostnames {
			err := dns.UpdateRecordAAAA(hostname, res6)
			if err != nil {
				log.Panic(err)
			}
		}
	}
}
