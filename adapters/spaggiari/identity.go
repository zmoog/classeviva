package spaggiari

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

var (
	noIdentity Identity
)

type Provider interface {
	Get() (Identity, error)
}

type IdentityProvider struct {
	Fetcher      Fetcher
	LoaderStorer LoaderStorer
}

func (p IdentityProvider) Get() (Identity, error) {
	fmt.Println("checking loaderstorer for identity")

	identity, exists, err := p.LoaderStorer.Load()
	if err != nil {
		return Identity{}, err
	}

	fmt.Println("exists: ", exists)

	now := time.Now().Format(time.RFC3339)

	fmt.Println("expire", identity.Expire)
	fmt.Println("now", now)

	if exists && now >= identity.Release && now < identity.Expire {
		fmt.Println("reusing existing identity")
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

type InMemoryLoaderStorer struct {
	identity Identity
}

func (ls InMemoryLoaderStorer) Load() (Identity, bool, error) {
	if ls.identity == noIdentity {
		fmt.Println("identity is not available in the store")
		return noIdentity, false, nil
	}

	fmt.Println("returning identity from the store")
	return ls.identity, true, nil
}

func (ls InMemoryLoaderStorer) Store(identity Identity) error {
	ls.identity = identity
	return nil
}

// FilesystemLoaderStorer loads and stores an Identity using the file system as a backing storage.
type FilesystemLoaderStorer struct{}

func (ls FilesystemLoaderStorer) Load() (Identity, bool, error) {
	path, err := ls.getSettingsDir()
	if err != nil {
		return noIdentity, false, err
	}

	configFilePath := filepath.Join(path, "identity.json")

	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		return noIdentity, false, nil
	}

	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return noIdentity, false, err
	}

	identity := Identity{}
	err = json.Unmarshal(data, &identity)
	if err != nil {
		return noIdentity, false, err
	}

	return identity, true, nil
}

func (ls FilesystemLoaderStorer) Store(identity Identity) error {
	path, err := ls.getSettingsDir()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(identity, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(path, "identity.json"), data, 0700)
	if err != nil {
		return err
	}

	return err
}

func (ls FilesystemLoaderStorer) getSettingsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(home, ".classeviva")
	err = os.MkdirAll(path, 0700)
	if err != nil {
		return "", err
	}

	return path, err
}

type Fetcher interface {
	Fetch() (Identity, error)
}

type IdentityFetcher struct {
	client   *http.Client
	username string
	password string
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
