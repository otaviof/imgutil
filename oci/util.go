package oci

import (
	"fmt"

	"github.com/containers/buildah"
	is "github.com/containers/image/v5/storage"
	"github.com/containers/storage"
	"github.com/containers/storage/pkg/unshare"
)

// bootstrapContainerStorage get a storage.Store instance.
func bootstrapContainerStorage() (storage.Store, error) {
	opts, err := storage.DefaultStoreOptions(unshare.IsRootless(), unshare.GetRootlessUID())
	if err != nil {
		return nil, err
	}
	store, err := storage.GetStore(opts)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, fmt.Errorf("unable to instantiate storage")
	}
	is.Transport.SetStore(store)
	return store, nil
}

// ReInit executes buildah steps to initialize unshare in a fork/exec fashion.
func ReInit() {
	if buildah.InitReexec() {
		return
	}
	unshare.MaybeReexecUsingUserNamespace(true)
}
