package spaggiari_test

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zmoog/classeviva/adapters/spaggiari"
	"github.com/zmoog/classeviva/mocks"
)

func TestInMemoryLoaderStorer(t *testing.T) {
	t.Run("Load from empty an store", func(t *testing.T) {
		ls := spaggiari.InMemoryLoaderStorer{}

		_, exists, err := ls.Load()
		if err != nil {
			t.Error(err)
		}

		if exists {
			t.Errorf("Expected false, got: [%t]", exists)
		}
	})

	t.Run("Check load and store", func(t *testing.T) {
		ls := spaggiari.InMemoryLoaderStorer{}

		err := ls.Store(spaggiari.Identity{Ident: "ident"})
		if err != nil {
			t.Error(err)
		}

		actual, exists, err := ls.Load()
		if err != nil {
			t.Error(err)
		}

		if !exists {
			t.Errorf("Expected true, got: [%t]", exists)
		}

		if actual.Ident != "ident" {
			t.Errorf("Expected 'ident', got: [%s]", actual.Ident)
		}
	})
}

func TestFilesystemLoaderStorer(t *testing.T) {
	t.Run("Load from empty an store", func(t *testing.T) {
		path, err := ioutil.TempDir("", "")
		if err != nil {
			t.Error(err)
		}

		ls := spaggiari.FilesystemLoaderStorer{
			Path: path,
		}

		_, exists, err := ls.Load()
		if err != nil {
			t.Error(err)
		}

		if exists {
			t.Errorf("Expected false, got: [%t]", exists)
		}
	})

	t.Run("Store and load the identity", func(t *testing.T) {
		path, err := ioutil.TempDir("", "")
		if err != nil {
			t.Error(err)
		}
		ls := spaggiari.FilesystemLoaderStorer{Path: path}

		err = ls.Store(spaggiari.Identity{Ident: "ident"})
		if err != nil {
			t.Error(err)
		}

		actual, exists, err := ls.Load()
		if err != nil {
			t.Error(err)
		}

		if !exists {
			t.Errorf("Expected true, got: [%t]", exists)
		}

		if actual.Ident != "ident" {
			t.Errorf("Expected 'ident', got: [%s]", actual.Ident)
		}
	})
}

func TestIdentityProvider(t *testing.T) {

	t.Run("Get existent valid identity", func(t *testing.T) {
		fetcher := &mocks.Fetcher{}
		loaderStorer := &mocks.LoaderStorer{}

		identityProvider := spaggiari.IdentityProvider{
			Fetcher:      fetcher,
			LoaderStorer: loaderStorer,
		}

		loaderStorer.On("Load").Return(
			spaggiari.Identity{
				Ident:     "123456",
				ID:        "123456",
				FirstName: "John",
				LastName:  "Doe",
				Token:     "123",
				Release:   time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
				Expire:    time.Now().Add(1 * time.Hour).Format(time.RFC3339),
			}, true, nil)
		loaderStorer.AssertNotCalled(t, "Store")
		fetcher.AssertNotCalled(t, "Load")

		i, err := identityProvider.Get()

		fetcher.AssertExpectations(t)
		loaderStorer.AssertExpectations(t)

		assert.Nil(t, err)
		assert.Equal(t, "123456", i.Ident)
	})

	t.Run("Get existent expired identity", func(t *testing.T) {
		fetcher := &mocks.Fetcher{}
		loaderStorer := &mocks.LoaderStorer{}

		identityProvider := spaggiari.IdentityProvider{
			Fetcher:      fetcher,
			LoaderStorer: loaderStorer,
		}

		identity := spaggiari.Identity{
			Ident:     "123456",
			ID:        "123456",
			FirstName: "John",
			LastName:  "Doe",
			Token:     "123",
			Release:   time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
			Expire:    time.Now().Add(1 * time.Hour).Format(time.RFC3339),
		}
		loaderStorer.On("Load").Return(
			spaggiari.Identity{
				Ident:     "123456",
				ID:        "123456",
				FirstName: "John",
				LastName:  "Doe",
				Token:     "123",
				Release:   "2022-04-23T07:53:55+02:00",
				Expire:    "2022-04-23T09:23:55+02:00",
			}, true, nil,
		)
		fetcher.On("Fetch").Return(
			identity, nil,
		)
		loaderStorer.On("Store", identity).Return(nil)

		i, err := identityProvider.Get()

		fetcher.AssertExpectations(t)
		loaderStorer.AssertExpectations(t)

		assert.Nil(t, err)
		assert.Equal(t, "123456", i.Ident)
	})

	t.Run("Get NON existent identity", func(t *testing.T) {
		fetcher := &mocks.Fetcher{}
		loaderStorer := &mocks.LoaderStorer{}

		identityProvider := spaggiari.IdentityProvider{
			Fetcher:      fetcher,
			LoaderStorer: loaderStorer,
		}

		loaderStorer.On("Load").Return(spaggiari.Identity{}, false, nil)
		loaderStorer.On("Store", mock.AnythingOfType("spaggiari.Identity")).Return(nil)
		fetcher.On("Fetch").Return(
			spaggiari.Identity{
				Ident:     "123456",
				ID:        "123456",
				FirstName: "John",
				LastName:  "Doe",
				Token:     "123",
				Release:   "2022-04-24T07:53:55+02:00",
				Expire:    "2022-04-24T09:23:55+02:00",
			}, nil,
		)

		i, err := identityProvider.Get()

		fetcher.AssertExpectations(t)
		loaderStorer.AssertExpectations(t)

		assert.Nil(t, err)
		assert.Equal(t, "123456", i.Ident)
	})
}