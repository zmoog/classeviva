package spaggiari

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	baseUrl = "https://web.spaggiari.eu/rest/v1"
)

type Identity struct {
	Ident     string `json:"ident"`
	ID        string
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Token     string `json:"token"`
	Release   string `json:"release"`
	Expire    string `json:"expire"`
}

func From(username, password string) (Adapter, error) {
	httpClient := http.Client{}

	adapter := Adapter{
		client: &httpClient,
		headers: map[string]string{
			"User-Agent":   "zorro/1.0",
			"Z-Dev-Apikey": "+zorro+",
			"Content-Type": "application/json",
		},
		identityProvider: IdentityProvider{
			Fetcher: IdentityFetcher{
				username: username,
				password: password,
				client:   &httpClient,
			},
			LoaderStorer: FilesystemLoaderStorer{},
		},
	}

	return adapter, nil
}

type Adapter struct {
	headers          map[string]string
	client           *http.Client
	identityProvider IdentityProvider
}

func (c Adapter) List() ([]Grade, error) {
	identity, err := c.identityProvider.Get()
	if err != nil {
		return []Grade{}, err
	}

	url := baseUrl + "/students/" + identity.ID + "/grades"

	req, err := c.newRequest("GET", url, nil, identity)
	if err != nil {
		return []Grade{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return []Grade{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return []Grade{}, fmt.Errorf("failed to fetch grades, status_code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Grade{}, err
	}

	envelope := map[string][]Grade{}

	err = json.Unmarshal(body, &envelope)
	if err != nil {
		return []Grade{}, err
	}

	return envelope["grades"], nil
}

func (c Adapter) newRequest(method, url string, body io.Reader, identity Identity) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for key, value := range c.headers {
		req.Header.Add(key, value)
	}

	req.Header.Add("Z-Auth-Token", identity.Token)
	return req, nil
}
