package main

import (
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/http"
	"time"
)

func main() {
	network := "tcp"
	address := "13.208.212.66:1080"
	auth := proxy.Auth{
		User:     "test",
		Password: "1234",
	}

	client, _ := newHttpClient(network, address, &auth)

	// Make an HTTP GET request
	resp, err := client.Get("https://ifconfig.me")
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

type dialContextFunc func(ctx context.Context, network, address string) (net.Conn, error)

func newDialContext(network, address string, auth *proxy.Auth) (dialContextFunc, error) {
	baseDialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	if address != "" {
		dialSocksProxy, err := proxy.SOCKS5(network, address, auth, baseDialer)
		if err != nil {
			return nil, err
		}

		contextDialer, ok := dialSocksProxy.(proxy.ContextDialer)
		if !ok {
			return nil, err
		}

		return contextDialer.DialContext, nil
	} else {
		return baseDialer.DialContext, nil
	}
}

func newHttpClient(network, address string, auth *proxy.Auth) (*http.Client, error) {
	dialContext, err := newDialContext(network, address, auth)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			ForceAttemptHTTP2:     false, // disable http2
			DisableCompression:    true,  // To get the original response from the server, set Transport.DisableCompression to true.
		},
	}

	return client, nil
}
