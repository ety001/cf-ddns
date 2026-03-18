package main

import (
	"net"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"crypto/tls"
)

var log = logrus.New()
var Version string

type Config struct {
	ipAddress    string
	ipv6Address  string
	noVerify     bool
	interval     int
	cfEmail      string
	cfApiKey     string
	cfZoneId     string
	ipv4Hostname string
	ipv6Hostname string
}

func runUpdate(cfg Config, dns *CFDNSUpdater) {
	// Update A record (IPv4) if IPV4_HOSTNAME is set
	if cfg.ipv4Hostname != "" {
		var ipv4 net.IP
		if cfg.ipAddress != "" {
			ipv4 = net.ParseIP(cfg.ipAddress)
		} else {
			httpClient := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.noVerify},
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
				log.Errorf("Failed to get external IPv4: %v", err)
				return
			}
		}
		err := dns.UpdateRecordA(cfg.ipv4Hostname, ipv4)
		if err != nil {
			log.Errorf("Failed to update A record: %v", err)
			return
		}
		log.Infof("Checked A record for %s (current IP: %s)", cfg.ipv4Hostname, ipv4)
	}

	// Update AAAA record (IPv6) if IPV6_HOSTNAME is set
	if cfg.ipv6Hostname != "" {
		var ipv6 net.IP
		if cfg.ipv6Address != "" {
			ipv6 = net.ParseIP(cfg.ipv6Address)
		} else {
			httpClient := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.noVerify},
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
				log.Errorf("Failed to get external IPv6: %v", err)
				return
			}
		}
		err := dns.UpdateRecordAAAA(cfg.ipv6Hostname, ipv6)
		if err != nil {
			log.Errorf("Failed to update AAAA record: %v", err)
			return
		}
		log.Infof("Checked AAAA record for %s (current IP: %s)", cfg.ipv6Hostname, ipv6)
	}
}

func main() {
	var (
		app = kingpin.New("cf-ddns", "Cloudflare DynDNS Updater").Version(Version)

		ipAddress   = app.Flag("ip-address", "Skip resolving external IP and use provided IP").String()
		ipv6Address = app.Flag("ipv6-address", "Skip resolving external IPv6 and use provided IPv6").String()
		noVerify    = app.Flag("no-verify", "Don't verify ssl certificates").Bool()
		interval    = app.Flag("interval", "Run in loop mode, checking IP every N minutes (0 = run once)").Default("0").Int()

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

	cfg := Config{
		ipAddress:    *ipAddress,
		ipv6Address:  *ipv6Address,
		noVerify:     *noVerify,
		interval:     *interval,
		cfEmail:      *cfEmail,
		cfApiKey:     *cfApiKey,
		cfZoneId:     *cfZoneId,
		ipv4Hostname: ipv4Hostname,
		ipv6Hostname: ipv6Hostname,
	}

	// Run in loop mode if interval is set
	if cfg.interval > 0 {
		log.Infof("Starting loop mode, checking IP every %d minutes", cfg.interval)
		ticker := time.NewTicker(time.Duration(cfg.interval) * time.Minute)
		defer ticker.Stop()

		// Run once on start
		runUpdate(cfg, dns)

		for range ticker.C {
			runUpdate(cfg, dns)
		}
	} else {
		// Run once
		runUpdate(cfg, dns)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
