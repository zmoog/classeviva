package spaggiari_test

import (
	"io/ioutil"
	"testing"

	"github.com/zmoog/classeviva/adapters/spaggiari"
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
