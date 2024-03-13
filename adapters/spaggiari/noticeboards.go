package spaggiari

import "encoding/json"

type NoticeboardsReceiver interface {
	List() ([]Noticeboard, error)
}

type noticeboardsReceiver struct {
	Client           SpaggiariClient
	IdentityProvider Provider
}

func (r noticeboardsReceiver) List() ([]Noticeboard, error) {
	identity, err := r.IdentityProvider.Get()
	if err != nil {
		return []Noticeboard{}, err
	}

	url := baseUrl + "/students/" + identity.ID + "/noticeboard"

	items := []Noticeboard{}

	err = r.Client.Get(url, func(body []byte) error {
		envelope := map[string][]Noticeboard{}

		err := json.Unmarshal(body, &envelope)
		if err != nil {
			return err
		}

		items = envelope["items"]

		return nil
	})
	if err != nil {
		return []Noticeboard{}, err
	}

	return items, nil
}
