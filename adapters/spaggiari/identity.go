package spaggiari

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

var (
	noIdentity = Identity{}
)

type Provider interface {
	Get() (Identity, error)
}

type IdentityProvider struct {
	Fetcher      Fetcher
	LoaderStorer LoaderStorer
}

func (p IdentityProvider) Get() (Identity, error) {

	identity, exists, err := p.LoaderStorer.Load()
	if err != nil {
		return Identity{}, err
	}

	if exists {
		return identity, nil
	}

	identity, err = p.Fetcher.Fetch()
	if err != nil {
		return Identity{}, err
	}

	err = p.LoaderStorer.Store(identity)
	if err != nil {
		return Identity{}, err
	}

	return identity, nil
}

type LoaderStorer interface {
	Load() (Identity, bool, error)
	Store(Identity) error
}

type IdentityLoaderStorer struct {
	identity Identity
}

func (ls IdentityLoaderStorer) Load() (Identity, bool, error) {
	if ls.identity == noIdentity {
		fmt.Println("identity is not available in the store")
		return noIdentity, false, nil
	}

	fmt.Println("returning identity from the store")
	return ls.identity, true, nil
}

func (ls IdentityLoaderStorer) Store(identity Identity) error {
	ls.identity = identity
	return nil
}

type Fetcher interface {
	Fetch() (Identity, error)
}

type IdentityFetcher struct {
	username string
	password string
	client   *http.Client
}

func (f IdentityFetcher) Fetch() (Identity, error) {
	fmt.Println("fetching new identity")

	creds := map[string]string{
		"uid":  f.username,
		"pass": f.password,
	}

	payload, err := json.Marshal(creds)
	if err != nil {
		return Identity{}, err
	}

	req, err := http.NewRequest("POST", baseUrl+"/auth/login/", bytes.NewBuffer(payload))
	if err != nil {
		return Identity{}, err
	}

	req.Header.Add("User-Agent", "zorro/1.0")
	req.Header.Add("Z-Dev-Apikey", "+zorro+")
	req.Header.Add("Content-Type", "application/json")

	resp, err := f.client.Do(req)
	if err != nil {
		return Identity{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Identity{}, err
	}

	fmt.Println(string(body))

	identity := Identity{}

	err = json.Unmarshal(body, &identity)
	if err != nil {
		return Identity{}, err
	}

	// The identity ID is made of the `ident` without the leading
	// and trailing characters.
	// For example, with
	//   `ident = G9123456R`
	//   `id = 9123456`
	// without the leading `G` and the trailing `R`.
	//
	// The ID is required to make calls to other endpoints, like grades,
	// agenda, and so on.
	m := regexp.MustCompile("\\D")
	identity.ID = m.ReplaceAllString(identity.Ident, "")

	return identity, nil
}
