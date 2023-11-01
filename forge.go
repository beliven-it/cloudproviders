package cloudproviders

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type forge struct {
	token string
	url   string
}

type forgeServerData struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Size      string    `json:"size"`
	Region    string    `json:"region"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}

type forgeSiteData struct {
	ID        int       `json:"id"`
	ServerID  int       `json:"server_id"`
	Domain    string    `json:"name"`
	Username  string    `json:"username"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type forgeServersResponse struct {
	Servers []forgeServerData `json:"servers"`
}

type forgeSitesResponse struct {
	Sites []forgeSiteData `json:"sites"`
}

func (f *forge) request(method string, endpoint string, queryParams map[string]interface{}, body []byte) ([]byte, error) {
	url := fmt.Sprintf("%s%s", f.url, endpoint)

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + f.token,
	}

	return providerRequest(method, url, queryParams, body, headers)
}

func (f *forge) listServers() (forgeServersResponse, error) {
	var forgeServersResponse forgeServersResponse

	data, err := f.request("GET", "/servers", nil, nil)
	if err != nil {
		return forgeServersResponse, err
	}

	err = json.Unmarshal(data, &forgeServersResponse)
	if err != nil {
		return forgeServersResponse, err
	}

	return forgeServersResponse, nil
}

func (f *forge) listSitesForServer(serverID int) ([]forgeSiteData, error) {
	endpoint := fmt.Sprintf("/servers/%d/sites", serverID)

	data, err := f.request("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var response forgeSitesResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	sites := []forgeSiteData{}
	for _, site := range response.Sites {
		site.ServerID = serverID
		sites = append(sites, site)
	}

	return sites, nil
}

func (f *forge) ListSites() ([]ProvideSiteResponse, error) {

	servers, err := f.listServers()
	if err != nil {
		return nil, err
	}

	sites := []ProvideSiteResponse{}
	wg := sync.WaitGroup{}
	for _, server := range servers.Servers {
		wg.Add(1)

		go func(s forgeServerData) {
			defer wg.Done()

			listOfSites, err := f.listSitesForServer(s.ID)
			if err != nil {
				return
			}

			for _, site := range listOfSites {
				siteConverted := ProvideSiteResponse(site)
				sites = append(sites, siteConverted)
			}

		}(server)
	}

	wg.Wait()
	return sites, nil
}

func (f *forge) ListServers() ([]ProviderServerResponse, error) {
	response, err := f.listServers()
	if err != nil {
		return nil, err
	}

	servers := []ProviderServerResponse{}
	for _, server := range response.Servers {
		servers = append(servers, ProviderServerResponse{
			ID:        server.ID,
			Name:      server.Name,
			IPAddress: server.IPAddress,
			Region:    server.Region,
			Tier:      server.Size,
			CreatedAt: server.CreatedAt,
		})
	}

	return servers, nil
}

func NewForge(options ProviderOptions) Provider {
	if options.URL == "" {
		options.URL = "https://forge.laravel.com/api/v1"
	}

	return &forge{
		token: options.Token,
		url:   options.URL,
	}
}
