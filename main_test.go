package main

import (
	"os"
	"path"
	"testing"
)

func TestIsNewAndStoreVersion(t *testing.T) {
	storage := path.Join(os.TempDir(), "mikrocheck")
	ver := "1.2.3"

	os.Remove(storage)

	if !isNew(ver, storage) {
		t.Error("should be true if storage doesn't exist yet")
	}

	storeVersion("1.2.4", storage)

	if !isNew(ver, storage) {
		t.Error("should be true if versions aren't equal")
	}
	if isNew("1.2.4", storage) {
		t.Error("should be false if versions are equal")
	}

	// should change existed version
	storeVersion(ver, storage)

	if isNew(ver, storage) {
		t.Error("should be false if versions are equal")
	}
}

func TestGetLastVersion(t *testing.T) {
	ver, info := getLastVersion()

	if len(ver) == 0 {
		t.Error("Version shouldn't be empty")
	}

	if len(info) == 0 {
		t.Error("Info shouldn't be empty")
	}
}
