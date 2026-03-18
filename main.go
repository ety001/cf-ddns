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
	)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	// Get hostnames from environment variables
	ipv4Hostname := os.Getenv("IPV4_HOSTNAME")
	ipv6Hostname := os.Getenv("IPV6_HOSTNAME")

	// Validate that at least one hostname is provided
	if ipv4Hostname == "" && ipv6Hostname == "" {
		log.Panic("At least one of IPV4_HOSTNAME or IPV6_HOSTNAME must be set")
	}

	var dns *CFDNSUpdater
	var err error

	if dns, err = NewCFDNSUpdater(*cfZoneId, *cfApiKey, *cfEmail, log.WithField("component", "cf-dns-updater")); err != nil {
		log.Panic(err)
	}

	// Update A record (IPv4) if IPV4_HOSTNAME is set
	if ipv4Hostname != "" {
		var ipv4 net.IP
		if *ipAddress != "" {
			ipv4 = net.ParseIP(*ipAddress)
		} else {
			httpClient := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: *noVerify},
				},
			}
			ipv4Endpoint := getEnvOrDefault("IPV4_ENDPOINT", "https://ipv4.ip.sb")
			ipService := &HTTPBasedIPService{
				HttpClient:   httpClient,
				IPv4Endpoint: ipv4Endpoint,
				IPv6Endpoint: "", // Not needed for IPv4 only
			}
			var err error
			ipv4, err = ipService.GetExternalIP()
			if err != nil {
				log.Panic(err)
			}
		}
		err := dns.UpdateRecordA(ipv4Hostname, ipv4)
		if err != nil {
			log.Panic(err)
		}
		log.Infof("Updated A record for %s to %s", ipv4Hostname, ipv4)
	}

	// Update AAAA record (IPv6) if IPV6_HOSTNAME is set
	if ipv6Hostname != "" {
		var ipv6 net.IP
		if *ipv6Address != "" {
			ipv6 = net.ParseIP(*ipv6Address)
		} else {
			httpClient := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: *noVerify},
				},
			}
			ipv6Endpoint := getEnvOrDefault("IPV6_ENDPOINT", "https://ipv6.ip.sb")
			ipService := &HTTPBasedIPService{
				HttpClient:   httpClient,
				IPv4Endpoint: "", // Not needed for IPv6 only
				IPv6Endpoint: ipv6Endpoint,
			}
			var err error
			ipv6, err = ipService.GetExternalIPv6()
			if err != nil {
				log.Panic(err)
			}
		}
		err := dns.UpdateRecordAAAA(ipv6Hostname, ipv6)
		if err != nil {
			log.Panic(err)
		}
		log.Infof("Updated AAAA record for %s to %s", ipv6Hostname, ipv6)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
