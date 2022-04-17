package spaggiari

import (
	"bytes"
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
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Token     string `json:"token"`
	Release   string `json:"release"`
	Expire    string `json:"expire"`
}

func From(username, password string) (Client, error) {

	client := Client{
		headers: map[string]string{
			"User-Agent":   "zorro/1.0",
			"Z-Dev-Apikey": "+zorro+",
			"Content-Type": "application/json",
		},
		client: &http.Client{},
	}

	identity, err := client.getIdentity(username, password)
	if err != nil {
		return Client{}, err
	}

	client.headers["Z-Auth-Token"] = identity.Token

	return client, nil
}

type Client struct {
	headers map[string]string
	client  *http.Client
}

func (c Client) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for key, value := range c.headers {
		req.Header.Add(key, value)
	}

	return req, nil
}

func (c Client) getIdentity(username, password string) (Identity, error) {

	creds := map[string]string{
		"uid":  username,
		"pass": password,
	}

	payload, err := json.Marshal(creds)
	if err != nil {
		return Identity{}, err
	}

	req, err := c.newRequest("POST", baseUrl+"/auth/login/", bytes.NewBuffer(payload))
	if err != nil {
		return Identity{}, err
	}

	resp, err := c.client.Do(req)
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

	return identity, nil
}
