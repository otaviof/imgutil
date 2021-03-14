package oci

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/buildpacks/imgutil"
	"github.com/containers/buildah"
	"github.com/containers/buildah/imagebuildah"
	is "github.com/containers/image/v5/storage"
	"github.com/containers/image/v5/types"
	"github.com/containers/storage"
)

//
// TODO:
//	* create a `doc.go` file;
//

// Image implements imgutil.Image interface, using buildah container-manager to handle local images.
type Image struct {
	ctx      context.Context  // shared context
	from     string           // from image-tag
	repoName string           // target image-tag
	store    storage.Store    // container storage instance
	builder  *buildah.Builder // builder instance
}

// systemContext global system context instance.
var systemContext = &types.SystemContext{}

// Name returns the current image repository name, in short the name.
func (i *Image) Name() string {
	return i.repoName
}

// Rename renames the current image repository name.
func (i *Image) Rename(name string) {
	i.repoName = name
}

// OS returns the OS string.
func (i *Image) OS() (string, error) {
	return i.builder.OS(), nil
}

// SetOS sets the os name.
func (i *Image) SetOS(name string) error {
	i.builder.SetOS(name)
	return nil
}

// OSVersion returns the OSVersion string.
func (i *Image) OSVersion() (string, error) {
	return i.builder.Docker.OSVersion, nil
}

// SetOSVersion operating system version is not currently supported.
func (i *Image) SetOSVersion(version string) error {
	// TODO: check the impact of not setting os-version;
	return nil
}

// Architecture returns the architecture string.
func (i *Image) Architecture() (string, error) {
	return i.builder.Architecture(), nil
}

// SetArchitecture sets the image arachitecture.
func (i *Image) SetArchitecture(arch string) error {
	i.builder.SetArchitecture(arch)
	return nil
}

// CreatedAt returns the time of image creation.
func (i *Image) CreatedAt() (time.Time, error) {
	return i.builder.Docker.Created, nil
}

// SetCmd set a new cmd slice.
func (i *Image) SetCmd(cmd ...string) error {
	i.builder.SetCmd(cmd)
	return nil
}

// Entrypoint returns the current entrypoint.
func (i *Image) Entrypoint() ([]string, error) {
	return i.builder.Entrypoint(), nil
}

// SetEntrypoint set a new entrypoint slice.
func (i *Image) SetEntrypoint(entrypoint ...string) error {
	i.builder.SetEntrypoint(entrypoint)
	return nil
}

// SetWorkingDir set the working-directory.
func (i *Image) SetWorkingDir(dir string) error {
	i.builder.SetWorkDir(dir)
	return nil
}

// ManifestSize returns the actual manifest size.
func (i *Image) ManifestSize() (int64, error) {
	return i.builder.Docker.Size, nil
}

// Env retrieve environment variable value, or empty in case of not found.
func (i *Image) Env(key string) (string, error) {
	// TODO: refactor to reuse the same logic that's on ``local.Image.Env()` (DRY);
	for _, envVar := range i.builder.Env() {
		parts := strings.Split(envVar, "=")
		if key == parts[0] {
			return parts[1], nil
		}
	}
	return "", nil
}

// SetEnv set a environment variable key/value.
func (i *Image) SetEnv(k, v string) error {
	i.builder.SetEnv(k, v)
	return nil
}

// Identifier returns the unique image identifier.
func (i *Image) Identifier() (imgutil.Identifier, error) {
	return IDIdentifier{
		ImageID: i.builder.ContainerID,
	}, nil
}

// Found make sure the image containers a unique identifier.
func (i *Image) Found() bool {
	return i.builder.ContainerID != ""
}

// Label returns the given lable name (key), or empty when not found.
func (i *Image) Label(key string) (string, error) {
	// TODO: it should return a specific error when informed label is not found;
	labels := i.builder.Labels()
	return labels[key], nil
}

// Labels returns all labels present in the working container.
func (i *Image) Labels() (map[string]string, error) {
	return i.builder.Labels(), nil
}

// SetLabel writes the label key-value pair on the working image.
func (i *Image) SetLabel(k, v string) error {
	i.builder.SetLabel(k, v)
	return nil
}

// RemoveLabel unset a given label name from working container.
func (i *Image) RemoveLabel(key string) error {
	i.builder.UnsetLabel(key)
	return nil
}

func (i *Image) AddLayer(path string) error {
	panic(fmt.Sprintf("[NOT-IMPLEMENTED] AddLayer(path='%s')", path))
	return nil
}

func (i *Image) AddLayerWithDiffID(path, diffID string) error {
	panic(fmt.Sprintf("[NOT-IMPLEMENTED] AddLayerWithDiffID(path='%s',diffID='%s')", path, diffID))
	return nil
}

func (i *Image) TopLayer() (string, error) {
	panic("[NOT-IMPLEMENTED] TopLayer()")
	return "", nil
}

func (i *Image) GetLayer(diffID string) (io.ReadCloser, error) {
	panic(fmt.Sprintf("[NOT-IMPLEMENTED] GetLayer(diffID='%s')", diffID))
	return nil, nil
}

func (i *Image) Rebase(baseTopLayer string, baseImage imgutil.Image) error {
	panic(fmt.Sprintf("[NOT-IMPLEMENTED] Rebase(baseTopLayer='%s',baseImage'%#v')", baseTopLayer, baseImage))
	return nil
}

func (i *Image) ReuseLayer(diffID string) error {
	panic(fmt.Sprintf("[NOT-IMPLEMENTED] ReuseLayer(diffID='%s')", diffID))
	return nil
}

// Delete removes the working container.
func (i *Image) Delete() error {
	return i.builder.Delete()
}

// commit perform the commit action for informed image name.
func (i *Image) commit(name string) error {
	imageRef, err := is.Transport.ParseStoreReference(i.store, name)
	if err != nil {
		return err
	}

	id, ref, _, err := i.builder.Commit(i.ctx, imageRef, buildah.CommitOptions{
		SystemContext: systemContext,
		Compression:   imagebuildah.Gzip,
	})
	log.Printf("id='%s', ref='%s'", id, ref)
	return err
}

// Save commit working container to persistent image storage, using default and additional names.
func (i *Image) Save(additionalNames ...string) error {
	names := append(additionalNames, i.repoName)
	for _, name := range names {
		if err := i.commit(name); err != nil {
			return err
		}
	}
	return nil
}

// bootstrap perform the initial steps to bootstrap the container manager.
func (i *Image) bootstrap() error {
	var err error
	if i.store, err = bootstrapContainerStorage(); err != nil {
		return err
	}
	i.builder, err = buildah.NewBuilder(i.ctx, i.store, buildah.BuilderOptions{
		CommonBuildOpts:  &buildah.CommonBuildOptions{},
		ConfigureNetwork: buildah.NetworkDefault,
		Format:           buildah.OCIv1ImageManifest,
		FromImage:        i.from,
		Isolation:        buildah.IsolationChroot,
		ReportWriter:     os.Stderr,
		SystemContext:    systemContext,
	})
	return err
}

// NewImage instantiate a new OCI image, executing re-init and other actions to bootstrap the
// container manager and its dependencies.
func NewImage(ctx context.Context, repoName, from string) (*Image, error) {
	ReInit()
	i := &Image{
		ctx:      ctx,
		repoName: repoName,
		from:     from,
	}
	return i, i.bootstrap()
}
