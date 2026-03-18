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

type HTTPBasedIPService struct {
	HttpClient   *http.Client
	IPv4Endpoint string
	IPv6Endpoint string
}

type IPResponse struct {
	IP string
}

func (h *HTTPBasedIPService) GetExternalIP() (net.IP, error) {
	r, err := h.HttpClient.Get(h.IPv4Endpoint)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	var resp IPResponse
	json.NewDecoder(r.Body).Decode(&resp)
	return net.ParseIP(resp.IP), nil
}

func (h *HTTPBasedIPService) GetExternalIPv6() (net.IP, error) {
	r, err := h.HttpClient.Get(h.IPv6Endpoint)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	var resp IPResponse
	json.NewDecoder(r.Body).Decode(&resp)
	return net.ParseIP(resp.IP), nil
}
