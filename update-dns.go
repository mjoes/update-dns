package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	ZONE_ID string = os.Getenv("ZONE_ID")
	API_KEY string = os.Getenv("API_KEY")
	EMAIL   string = os.Getenv("EMAIL")
)

func get_ip() string {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		log.Fatalf("Error fetching IP: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	return string(body)
}

func main() {
	ip4 := get_ip()
	update_dns(os.Getenv("DNS_RECORD_ID_1"), "fools-paradise.com", ip4)
	update_dns(os.Getenv("DNS_RECORD_ID_2"), "www.fools-paradise.com", ip4)
}

func update_dns(dns_record_id string, name string, ip4 string) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", ZONE_ID, dns_record_id)

	payload := map[string]interface{}{
		"comment": "Domain verification record",
		"name":    name,
		"proxied": true,
		"ttl":     1,
		"content": ip4,
		"type":    "A",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Key", API_KEY)
	req.Header.Set("X-Auth-Email", EMAIL)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: received non-200 response status:", resp.Status)
		return
	}

	fmt.Println("Successfully sent PUT request:", name)
}
