package spaggiari

import (
	"encoding/json"
	"time"
)

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

	entries := []AgendaEntry{}

	err = r.Client.Get(url, func(body []byte) error {
		envelope := map[string][]AgendaEntry{}

		err = json.Unmarshal(body, &envelope)
		if err != nil {
			return err
		}

		entries = envelope["agenda"]

		return nil
	})

	return entries, nil
}
