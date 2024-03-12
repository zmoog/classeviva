package spaggiari

import (
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func NewSpaggiariClient(identityProvider Provider) SpaggiariClient {
	return SpaggiariClient{
		identityProvider: identityProvider,
		httpClient:       &http.Client{},
	}
}

func New(username, password, identityStorePath string) (Adapter, error) {
	httpClient := http.Client{}
	identityProvider := IdentityProvider{
		Fetcher: IdentityFetcher{
			Username: username,
			Password: password,
			Client:   &httpClient,
		},
		LoaderStorer: FilesystemLoaderStorer{
			Path: identityStorePath,
		},
	}

	spaggiariClient := NewSpaggiariClient(identityProvider)

	adapter := Adapter{
		Agenda: agendaReceiver{
			Client:           spaggiariClient,
			IdentityProvider: identityProvider,
		},
		Grades: gradeReceiver{
			Client:           spaggiariClient,
			IdentityProvider: identityProvider,
		},
		Noticeboards: noticeboardsReceiver{
			Client:           spaggiariClient,
			IdentityProvider: identityProvider,
		},
	}

	return adapter, nil
}

type Adapter struct {
	Agenda       AgendaReceiver
	Grades       GradesReceiver
	Noticeboards NoticeboardsReceiver
}

//go:generate mockery --name Client
type Client interface {
	Get(url string, unmarshal Unmarshal) error
}

type SpaggiariClient struct {
	httpClient       HTTPClient
	identityProvider Provider
}

type Unmarshal func(body []byte) error

func (c SpaggiariClient) Get(url string, unmarshal Unmarshal) error {
	// This is a dummy function to make the linter happy
	identity, err := c.identityProvider.Get()
	if err != nil {
		return err
	}

	req, err := newRequest("GET", url, nil, identity)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to fetch grades, status_code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Trace(string(body))

	err = unmarshal(body)
	if err != nil {
		return err
	}

	return nil
}

func newRequest(method, url string, body io.Reader, identity Identity) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for key, value := range BaseHeaders {
		req.Header.Add(key, value)
	}

	req.Header.Add("Z-Auth-Token", identity.Token)

	return req, nil
}
