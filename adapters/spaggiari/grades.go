package spaggiari

import "encoding/json"

type GradesReceiver interface {
	List() ([]Grade, error)
}

type gradeReceiver struct {
	Client           SpaggiariClient
	IdentityProvider Provider
}

func (r gradeReceiver) List() ([]Grade, error) {
	identity, err := r.IdentityProvider.Get()
	if err != nil {
		return []Grade{}, err
	}

	url := baseUrl + "/students/" + identity.ID + "/grades"

	grades := []Grade{}

	err = r.Client.Get(url, func(body []byte) error {
		envelope := map[string][]Grade{}

		err := json.Unmarshal(body, &envelope)
		if err != nil {
			return err
		}

		grades = envelope["grades"]

		return nil
	})
	if err != nil {
		return []Grade{}, err
	}

	return grades, nil
}
