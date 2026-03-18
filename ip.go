package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

type IPService interface {
	GetExternalIP() (net.IP, error)
	GetExternalIPv6() (net.IP, error)
}

type FakeIPService struct {
	fakeIp   net.IP
	fakeIPv6 net.IP
}

func (f *FakeIPService) GetExternalIP() (net.IP, error) {
	if f.fakeIp == nil {
		return nil, fmt.Errorf("FakeIPService: No IP specified")
	}
	return f.fakeIp, nil
}

func (f *FakeIPService) GetExternalIPv6() (net.IP, error) {
	if f.fakeIPv6 == nil {
		return nil, fmt.Errorf("FakeIPService: No IPv6 specified")
	}
	return f.fakeIPv6, nil
}

type IpifyIPService struct {
	HttpClient *http.Client
}

type IpifyAPIResponse struct {
	IP string
}

func (i *IpifyIPService) GetExternalIP() (net.IP, error) {
	r, err := i.HttpClient.Get("https://api.ipify.org?format=json")
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	var resp IpifyAPIResponse
	json.NewDecoder(r.Body).Decode(&resp)
	return net.ParseIP(resp.IP), nil
}

func (i *IpifyIPService) GetExternalIPv6() (net.IP, error) {
	r, err := i.HttpClient.Get("https://api64.ipify.org?format=json")
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	var resp IpifyAPIResponse
	json.NewDecoder(r.Body).Decode(&resp)
	return net.ParseIP(resp.IP), nil
}
