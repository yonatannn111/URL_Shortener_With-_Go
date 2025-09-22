package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ipApiResp struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Status  string `json:"status"`
}

func lookupGeo(ip string) (country, city string) {
	if ip == "" { return "", "" }
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	client := http.Client{Timeout: 500 * time.Millisecond}
	resp, err := client.Get(url)
	if err != nil { return "", "" }
	defer resp.Body.Close()
	var r ipApiResp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil { return "", "" }
	if r.Status != "success" { return "", "" }
	return r.Country, r.City
}
