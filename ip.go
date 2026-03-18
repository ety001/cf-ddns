package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
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
	IP string `json:"ip"`
}

// parseIPResponse tries to parse the response body as either JSON or plain text
func parseIPResponse(body io.Reader) (string, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}

	content := strings.TrimSpace(string(data))

	// Try JSON format first: {"ip":"x.x.x.x"}
	var jsonResp IPResponse
	if err := json.Unmarshal(data, &jsonResp); err == nil && jsonResp.IP != "" {
		return jsonResp.IP, nil
	}

	// Fall back to plain text format
	if net.ParseIP(content) != nil {
		return content, nil
	}

	return "", fmt.Errorf("unable to parse IP from response: %s", content)
}

func (h *HTTPBasedIPService) GetExternalIP() (net.IP, error) {
	r, err := h.HttpClient.Get(h.IPv4Endpoint)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	ipStr, err := parseIPResponse(r.Body)
	if err != nil {
		return nil, err
	}
	return net.ParseIP(ipStr), nil
}

func (h *HTTPBasedIPService) GetExternalIPv6() (net.IP, error) {
	r, err := h.HttpClient.Get(h.IPv6Endpoint)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	ipStr, err := parseIPResponse(r.Body)
	if err != nil {
		return nil, err
	}
	return net.ParseIP(ipStr), nil
}
