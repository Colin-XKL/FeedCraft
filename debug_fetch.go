package main

import (
	"fmt"
	"net"
	"net/url"
	"github.com/go-resty/resty/v2"
)

func validateURL(rawUrl string) error {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("invalid scheme: %s", u.Scheme)
	}

	ips, err := net.LookupIP(u.Hostname())
	if err != nil {
		return err
	}

	for _, ip := range ips {
		fmt.Printf("Resolved IP: %s, IsLoopback: %v, IsPrivate: %v\n", ip.String(), ip.IsLoopback(), ip.IsPrivate())
		if ip.IsLoopback() || ip.IsPrivate() {
			return fmt.Errorf("access to private IP %s is forbidden", ip.String())
		}
	}
	return nil
}

func main() {
	targetURL := "https://blog.colinx.one"
	fmt.Println("Testing URL:", targetURL)

	err := validateURL(targetURL)
	if err != nil {
		fmt.Println("Validation error:", err)
		return
	}
	fmt.Println("Validation success")

	client := resty.New()
	client.SetDebug(true)
	resp, err := client.R().
		SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36").
		Get(targetURL)

	if err != nil {
		fmt.Println("Fetch error:", err)
		return
	}
	fmt.Println("Fetch success, status:", resp.StatusCode())
	fmt.Printf("Content length: %d\n", len(resp.String()))
}
