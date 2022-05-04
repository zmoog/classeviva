package mocks

import "net/http"

type MockDoType func(req *http.Request) (*http.Response, error)

// inspired by https://levelup.gitconnected.com/mocking-outbound-http-calls-in-golang-9e5a044c2555
type MockClient struct {
	MockDo MockDoType
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.MockDo(req)
}
