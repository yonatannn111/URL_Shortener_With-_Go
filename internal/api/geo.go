package api

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

type IpApiResp struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Status  string `json:"status"`
}

// lookupGeo returns country and city for a given IP
func lookupGeo(ip string) (country, city string) {
	ip = strings.Split(ip, ":")[0] // remove port if exists

	// Skip local addresses
	if ip == "127.0.0.1" || ip == "::1" || net.ParseIP(ip).IsLoopback() {
		return "", ""
	}

	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	client := http.Client{Timeout: 500 * time.Millisecond}

	resp, err := client.Get(url)
	if err != nil {
		return "", ""
	}
	defer resp.Body.Close()

	var r IpApiResp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", ""
	}
	if r.Status != "success" {
		return "", ""
	}

	return r.Country, r.City
}
