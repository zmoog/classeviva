package spaggiari

import (
	"encoding/json"
	"fmt"
)

type NoticeboardsReceiver interface {
	List() ([]Noticeboard, error)
	DownloadAttachment(publicationID int, attachmentSequenceNumber int) ([]byte, error)
	SetAsRead(eventCode string, publicationID int) error
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

func (r noticeboardsReceiver) SetAsRead(eventCode string, publicationID int) error {

	identity, err := r.IdentityProvider.Get()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/students/%s/noticeboard/read/%s/%d", baseUrl, identity.ID, eventCode, publicationID)

	return r.Client.Post(url, func(body []byte) error {
		// I'm not interested in the response body if the request was successful.
		return nil
	})
}

func (r noticeboardsReceiver) DownloadAttachment(publicationID int, attachmentSequenceNumber int) ([]byte, error) {
	identity, err := r.IdentityProvider.Get()
	if err != nil {
		return []byte{}, err
	}

	url := fmt.Sprintf("%s/students/%s/noticeboard/attach/CF/%d/%d", baseUrl, identity.ID, publicationID, attachmentSequenceNumber)
	document := []byte{}

	err = r.Client.Get(url, func(body []byte) error {
		document = body
		return nil
	})
	if err != nil {
		return []byte{}, err
	}

	return document, nil
}
