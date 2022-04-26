package spaggiari

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	baseUrl = "https://web.spaggiari.eu/rest/v1"
)

func From(username, password, identityStorePath string) (Adapter, error) {
	httpClient := http.Client{}

	adapter := spaggiariAdapter{
		client: &httpClient,
		headers: map[string]string{
			"User-Agent":   "zorro/1.0",
			"Z-Dev-Apikey": "+zorro+",
			"Content-Type": "application/json",
		},
		IdentityProvider: IdentityProvider{
			Fetcher: IdentityFetcher{
				username: username,
				password: password,
				client:   &httpClient,
			},
			LoaderStorer: FilesystemLoaderStorer{
				Path: identityStorePath,
			},
		},
	}

	return adapter, nil
}

type Adapter interface {
	List() ([]Grade, error)
	ListAgenda(since, until time.Time) ([]AgendaEntry, error)
}

type spaggiariAdapter struct {
	headers          map[string]string
	client           *http.Client
	IdentityProvider IdentityProvider
}

func (c spaggiariAdapter) List() ([]Grade, error) {
	identity, err := c.IdentityProvider.Get()
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

func (c spaggiariAdapter) ListAgenda(since, until time.Time) ([]AgendaEntry, error) {
	identity, err := c.IdentityProvider.Get()
	if err != nil {
		return []AgendaEntry{}, err
	}

	_since := since.Format("20060102")
	_until := until.Format("20060102")

	url := baseUrl + "/students/" + identity.ID + "/agenda/all/" + _since + "/" + _until
	// fmt.Println(url)
	log.Debug(url)

	req, err := c.newRequest("GET", url, nil, identity)
	if err != nil {
		return []AgendaEntry{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return []AgendaEntry{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return []AgendaEntry{}, fmt.Errorf("failed to fetch agenda entries, status_code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []AgendaEntry{}, err
	}

	envelope := map[string][]AgendaEntry{}

	err = json.Unmarshal(body, &envelope)
	if err != nil {
		return []AgendaEntry{}, err
	}

	return envelope["agenda"], nil
}

func (c spaggiariAdapter) newRequest(method, url string, body io.Reader, identity Identity) (*http.Request, error) {
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
