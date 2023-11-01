package cloudserverproviders

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Provider interface {
	ListSites() ([]ProvideSiteResponse, error)
	ListServers() ([]ProviderServerResponse, error)
}

type ProviderServerResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	IPAddress string    `json:"ip_address"`
	Region    string    `json:"region"`
	Tier      string    `json:"tier"`
	CreatedAt time.Time `json:"created_at"`
}

type ProvideSiteResponse struct {
	ID        int       `json:"id"`
	ServerID  int       `json:"server_id"`
	Domain    string    `json:"domain"`
	Username  string    `json:"username"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type ProviderOptions struct {
	Token string
	URL   string
}

func providerRequest(method string, url string, queryParams map[string]interface{}, body []byte, headers map[string]string) ([]byte, error) {
	bodyReader := bufio.NewReader(bytes.NewBuffer(body))

	request, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	q := request.URL.Query()
	for k, v := range queryParams {
		q.Add(k, fmt.Sprintf("%v", v))
	}

	request.URL.RawQuery = q.Encode()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return io.ReadAll(response.Body)
}

func NewProvider(provider string, options ProviderOptions) Provider {
	switch provider {
	case "spinupwp":
		return NewSpinupwp(options)
	case "forge":
		return NewForge(options)
	default:
		return nil
	}
}
