package main

import (
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"net/http"
)

func main() {
	auth := proxy.Auth{
		User:     "test",
		Password: "1234",
	}

	// Set up SOCKS5 proxy dialer
	dialer, err := proxy.SOCKS5("tcp", "13.208.212.66:1080", &auth, proxy.Direct)
	if err != nil {
		fmt.Println("Error creating SOCKS5 dialer:", err)
		return
	}

	// Create a transport using the SOCKS5 proxy dialer
	tr := &http.Transport{Dial: dialer.Dial}

	// Create an HTTP client with the custom transport
	client := &http.Client{Transport: tr}

	// Make an HTTP GET request
	resp, err := client.Get("https://ifconfig.me/ip")
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return
	}

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Request failed with status:", resp.StatusCode)
		return
	}

	body, _ := getBody(resp)

	fmt.Println(body)
}

func getBody(response *http.Response) (string, error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
