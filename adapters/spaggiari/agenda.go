package spaggiari

import (
	"time"
	"encoding/json"
)

// func (c SpaggiariAdapter) ListAgenda(since, until time.Time) ([]AgendaEntry, error) {
// 	identity, err := c.IdentityProvider.Get()
// 	if err != nil {
// 		return []AgendaEntry{}, err
// 	}

// 	_since := since.Format("20060102")
// 	_until := until.Format("20060102")

// 	url := baseUrl + "/students/" + identity.ID + "/agenda/all/" + _since + "/" + _until
// 	log.Trace(string(url))

// 	req, err := c.newRequest("GET", url, nil, identity)
// 	if err != nil {
// 		return []AgendaEntry{}, err
// 	}

// 	resp, err := c.client.Do(req)
// 	if err != nil {
// 		return []AgendaEntry{}, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != 200 {
// 		return []AgendaEntry{}, fmt.Errorf("failed to fetch agenda entries, status_code: %d", resp.StatusCode)
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return []AgendaEntry{}, err
// 	}
// 	log.Trace(string(body))

// 	envelope := map[string][]AgendaEntry{}

// 	err = json.Unmarshal(body, &envelope)
// 	if err != nil {
// 		return []AgendaEntry{}, err
// 	}

// 	return envelope["agenda"], nil
// }

type AgendaReceiver interface {
	List(since, until time.Time) ([]AgendaEntry, error)
}

type agendaReceiver struct {
	Client           SpaggiariClient
	IdentityProvider Provider
}

func (r agendaReceiver) List(since, until time.Time) ([]AgendaEntry, error) {
	identity, err := r.IdentityProvider.Get()
	if err != nil {
		return []AgendaEntry{}, err
	}

	_since := since.Format("20060102")
	_until := until.Format("20060102")

	url := baseUrl + "/students/" + identity.ID + "/agenda/all/" + _since + "/" + _until

	items := []AgendaEntry{}

	err = r.Client.Get(url, func(body []byte) error {
		envelope := map[string][]AgendaEntry{}

		err = json.Unmarshal(body, &envelope)
		if err != nil {
			return err
		}

		items = envelope["agenda"]

		return nil
	})

	return items, nil
}
