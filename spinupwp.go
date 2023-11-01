package cloudserverproviders

import (
	"encoding/json"
	"fmt"
	"time"
)

type spinupwp struct {
	token string
	url   string
}

type spinupwpPagination struct {
	Next *string `json:"next"`
}

type spinupwpSiteData struct {
	ID        int       `json:"id"`
	ServerID  int       `json:"server_id"`
	Domain    string    `json:"domain"`
	Status    string    `json:"status"`
	Username  string    `json:"site_user"`
	CreatedAt time.Time `json:"created_at"`
}

type spinupwpServerData struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Size      string    `json:"size"`
	Region    string    `json:"region"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}

type spinupwpSitesResponse struct {
	Data       []spinupwpSiteData `json:"data"`
	Pagination spinupwpPagination `json:"pagination"`
}

type spinupwpServersResponse struct {
	Data       []spinupwpServerData `json:"data"`
	Pagination spinupwpPagination   `json:"pagination"`
}

func (s *spinupwp) request(method string, endpoint string, queryParams map[string]interface{}, body []byte) ([]byte, error) {
	url := fmt.Sprintf("%s%s", s.url, endpoint)

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + s.token,
	}

	return providerRequest(method, url, queryParams, body, headers)
}

func (s *spinupwp) listSites(nextPage int, results []spinupwpSiteData) ([]spinupwpSiteData, error) {
	queryParams := map[string]interface{}{
		"page": nextPage,
	}

	data, err := s.request("GET", "/sites", queryParams, nil)
	if err != nil {
		return results, err
	}

	var response spinupwpSitesResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return results, err
	}

	results = append(results, response.Data...)

	if response.Pagination.Next != nil {
		nextPage = nextPage + 1
		return s.listSites(nextPage, results)
	}

	return results, nil
}

func (s *spinupwp) listServers(nextPage int, results []spinupwpServerData) ([]spinupwpServerData, error) {
	queryParams := map[string]interface{}{
		"page": nextPage,
	}

	data, err := s.request("GET", "/servers", queryParams, nil)
	if err != nil {
		return results, err
	}

	var response spinupwpServersResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return results, err
	}

	results = append(results, response.Data...)

	if response.Pagination.Next != nil {
		nextPage = nextPage + 1
		return s.listServers(nextPage, results)
	}

	return results, nil
}

func (s *spinupwp) ListSites() ([]ProvideSiteResponse, error) {
	data, err := s.listSites(1, []spinupwpSiteData{})
	if err != nil {
		return nil, err
	}

	responses := []ProvideSiteResponse{}
	for _, site := range data {
		responses = append(responses, ProvideSiteResponse{
			ID:        site.ID,
			ServerID:  site.ServerID,
			Domain:    site.Domain,
			Username:  site.Username,
			Status:    site.Status,
			CreatedAt: site.CreatedAt,
		})
	}

	return responses, nil

}

func (s *spinupwp) ListServers() ([]ProviderServerResponse, error) {
	data, err := s.listServers(1, []spinupwpServerData{})
	if err != nil {
		return nil, err
	}

	responses := []ProviderServerResponse{}
	for _, server := range data {
		responses = append(responses, ProviderServerResponse{
			ID:        server.ID,
			Name:      server.Name,
			IPAddress: server.IPAddress,
			Region:    server.Region,
			Tier:      server.Size,
			CreatedAt: time.Now(),
		})
	}

	return responses, nil
}

func NewSpinupwp(options ProviderOptions) Provider {
	if options.URL == "" {
		options.URL = "https://api.spinupwp.app/v1"
	}

	return &spinupwp{
		token: options.Token,
		url:   options.URL,
	}
}
