package rajaongkir

import (
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Client struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewClient() *Client {
	proxyURLString := os.Getenv("HTTP_PROXY")

	transport := &http.Transport{}

	if proxyURLString != "" {
		proxyURL, err := url.Parse(proxyURLString)
		if err != nil {
			panic(err)
		}

		transport.Proxy = http.ProxyURL(proxyURL)
	}

	return &Client{
		apiKey: os.Getenv("RAJAONGKIR_API_KEY"),
		baseURL: strings.TrimRight(
			os.Getenv("RAJAONGKIR_BASE_URL"),
			"/",
		),
		client: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
	}
}
