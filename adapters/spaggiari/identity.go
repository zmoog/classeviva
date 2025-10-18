package spaggiari

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
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
	log.Debug("Looking for an existing identity..")
	identity, exists, err := p.LoaderStorer.Load()
	if err != nil {
		return Identity{}, err
	}

	log.Debug("exists: ", exists)

	now := time.Now().Format(time.RFC3339)
	log.Debug("expire", identity.Expire)
	log.Debug("now", now)

	if exists && now >= identity.Release && now < identity.Expire {
		log.Debug("reusing existing identity")
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

func (ls *InMemoryLoaderStorer) Load() (Identity, bool, error) {
	if ls.identity == noIdentity {
		log.Debug("identity is not available in the store")
		return noIdentity, false, nil
	}

	log.Debug("returning identity from the store")
	return ls.identity, true, nil
}

func (ls *InMemoryLoaderStorer) Store(identity Identity) error {
	ls.identity = identity
	return nil
}

// FilesystemLoaderStorer loads and stores an Identity using the file system as a backing storage.
type FilesystemLoaderStorer struct {
	Path string
}

func (ls FilesystemLoaderStorer) Load() (Identity, bool, error) {
	path, err := ls.getSettingsDir()
	if err != nil {
		return noIdentity, false, err
	}

	configFilePath := filepath.Join(path, "identity.json")

	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		log.Debugf("identity file [%s] does not exist", configFilePath)
		return noIdentity, false, nil
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return noIdentity, false, err
	}

	log.Debugf("loading identity from: [%s]", configFilePath)
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

	err = os.WriteFile(filepath.Join(path, "identity.json"), data, 0700)
	if err != nil {
		return err
	}

	return err
}

func (ls FilesystemLoaderStorer) getSettingsDir() (string, error) {
	path := filepath.Join(ls.Path, ".classeviva")
	log.Debugf("settings path is: [%s]", path)

	err := os.MkdirAll(path, 0700)
	if err != nil {
		return "", err
	}

	return path, err
}

type Fetcher interface {
	Fetch() (Identity, error)
}

type IdentityFetcher struct {
	Client   HTTPClient
	Username string
	Password string
}

func (f IdentityFetcher) Fetch() (Identity, error) {
	log.Debug("fetching new identity")

	creds := map[string]string{
		"uid":  f.Username,
		"pass": f.Password,
	}

	payload, err := json.Marshal(creds)
	if err != nil {
		return Identity{}, err
	}

	req, err := http.NewRequest("POST", baseUrl+"/auth/login/", bytes.NewBuffer(payload))
	if err != nil {
		return Identity{}, err
	}

	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Z-Dev-Apikey", apiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := f.Client.Do(req)
	if err != nil {
		return Identity{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		switch resp.StatusCode {
		case 403:
			return Identity{}, fmt.Errorf("fetcher: access denied to Classeviva API (status_code: %v). Hit: https://web.spaggiari.eu is not available to call from cloud provider.", resp.StatusCode)
		default:
			return Identity{}, fmt.Errorf("fetcher: failed to fetch identity (status_code: %v)", resp.StatusCode)
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Identity{}, err
	}
	log.Trace(string(body))

	identity := Identity{}

	err = json.Unmarshal(body, &identity)
	if err != nil {
		return Identity{}, fmt.Errorf("fetcher: failed to unmarshal identity %w", err)
	}

	// The student ID can be obtained from the `ident` field
	// value, by removing the leading and trailing characters.
	//
	// For example, with `ident = "G9123456R"` the student ID is
	// `9123456` (without the leading `G` and the trailing `R`).
	//
	// The student ID is required to make calls to other endpoints,
	// such as grades, agenda, noticeboards, and so on.
	m := regexp.MustCompile(`\D`)
	identity.ID = m.ReplaceAllString(identity.Ident, "")

	return identity, nil
}
