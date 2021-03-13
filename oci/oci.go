package oci

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/buildpacks/imgutil"
	"github.com/containers/buildah"
	"github.com/containers/buildah/imagebuildah"
	is "github.com/containers/image/v5/storage"
	"github.com/containers/image/v5/types"
	"github.com/containers/storage"
)

// OCI implements imgutil.Image interface, using buildah container-manager to handle local images.
type OCI struct {
	ctx      context.Context  // shared context
	from     string           // from image-tag
	repoName string           // target image-tag
	store    storage.Store    // container storage instance
	builder  *buildah.Builder // builder instance
}

// systemContext global system context instance.
var systemContext = &types.SystemContext{}

func (o *OCI) AddLayer(path string) error {
	panic(fmt.Sprintf("[NOT-IMPLEMENTED] AddLayer(path='%s')", path))
	return nil
}

func (o *OCI) AddLayerWithDiffID(path, diffID string) error {
	panic(fmt.Sprintf("[NOT-IMPLEMENTED] AddLayerWithDiffID(path='%s',diffID='%s')", path, diffID))
	return nil
}

func (o *OCI) Architecture() (string, error) {
	panic("[NOT-IMPLEMENTED] Architecture()")
	return "", nil
}

func (o *OCI) CreatedAt() (time.Time, error) {
	panic("[NOT-IMPLEMENTED] CreatedAt()")
	return time.Now(), nil
}

func (o *OCI) Delete() error {
	panic("[NOT-IMPLEMENTED] Delete()")
	return nil
}

func (o *OCI) Env(key string) (string, error) {
	panic(fmt.Sprintf("[NOT-IMPLEMENTED] Env(key='%s')", key))
	return "", nil
}

func (o *OCI) Found() bool {
	panic("[NOT-IMPLEMENTED] Found()")
	return false
}

func (o *OCI) GetLayer(diffID string) (io.ReadCloser, error) {
	panic(fmt.Sprintf("[NOT-IMPLEMENTED] GetLayer(diffID='%s')", diffID))
	return nil, nil
}

func (o *OCI) Identifier() (imgutil.Identifier, error) {
	panic("[NOT-IMPLEMENTED] Identifier()")
	return nil, nil
}

// Label returns the given lable name (key), or empty when not found.
func (o *OCI) Label(key string) (string, error) {
	// TODO: it should return a specific error when informed label is not found;
	labels := o.builder.Labels()
	return labels[key], nil
}

func (o *OCI) Labels() (map[string]string, error) {
	panic("[NOT-IMPLEMENTED] Labels()")
	return nil, nil
}

func (o *OCI) Name() string {
	panic("[NOT-IMPLEMENTED] Name()")
	return ""
}

func (o *OCI) OS() (string, error) {
	panic("[NOT-IMPLEMENTED] OS()")
	return "", nil
}

func (o *OCI) OSVersion() (string, error) {
	panic("[NOT-IMPLEMENTED] OSVersion()")
	return "", nil
}

func (o *OCI) Rebase(baseTopLayer string, baseImage imgutil.Image) error {
	panic(fmt.Sprintf("[NOT-IMPLEMENTED] Rebase(baseTopLayer='%s',baseImage'%#v')", baseTopLayer, baseImage))
	return nil
}

func (o *OCI) RemoveLabel(l string) error {
	panic(fmt.Sprintf("[NOT-IMPLEMENTED] RemoveLabel(l='%s')", l))
	return nil
}

func (o *OCI) Rename(name string) {
	panic(fmt.Sprintf("[NOT-IMPLEMENTED] Rename(name='%s')", name))
}

func (o *OCI) ReuseLayer(diffID string) error {
	panic(fmt.Sprintf("[NOT-IMPLEMENTED] ReuseLayer(diffID='%s')", diffID))
	return nil
}

func (o *OCI) commit(name string) error {
	imageRef, err := is.Transport.ParseStoreReference(o.store, name)
	if err != nil {
		return err
	}

	id, ref, _, err := o.builder.Commit(o.ctx, imageRef, buildah.CommitOptions{
		SystemContext: systemContext,
		Compression:   imagebuildah.Gzip,
	})
	log.Printf("id='%s', ref='%s'", id, ref)
	return err
}

func (o *OCI) Save(additionalNames ...string) error {
	names := append(additionalNames, o.repoName)
	for _, name := range names {
		if err := o.commit(name); err != nil {
			return err
		}
	}
	return nil
}

func (o *OCI) SetArchitecture(arch string) error {
	panic(fmt.Sprintf("[NOT-IMPLEMENTED] SetArchitecture(arch='%s')", arch))
	return nil
}

func (o *OCI) SetCmd(cmd ...string) error {
	panic(fmt.Sprintf("[NOT-IMPLEMENTED] SetCmd(cmd='%#v')", cmd))
	return nil
}

func (o *OCI) SetEntrypoint(entrypoint ...string) error {
	panic(fmt.Sprintf("[NOT-IMPLEMENTED] SetEntrypoint(entrypoint='%#v')", entrypoint))
	return nil
}

func (o *OCI) SetEnv(k, v string) error {
	panic(fmt.Sprintf("[NOT-IMPLEMENTED] SetEnv(k='%s', v='%s')", k, v))
	return nil
}

// SetLabel writes the label key-value pair on the working image.
func (o *OCI) SetLabel(k, v string) error {
	o.builder.SetLabel(k, v)
	return nil
}

func (o *OCI) SetOS(name string) error {
	panic(fmt.Sprintf("[NOT-IMPLEMENTED] SetOS(name='%s')", name))
	return nil
}

func (o *OCI) SetOSVersion(version string) error {
	panic(fmt.Sprintf("[NOT-IMPLEMENTED] SetOSVersion(version='%s')", version))
	return nil
}

func (o *OCI) SetWorkingDir(dir string) error {
	panic(fmt.Sprintf("[NOT-IMPLEMENTED] SetWorkingDir(dir='%s')", dir))
	return nil
}

func (o *OCI) TopLayer() (string, error) {
	panic("[NOT-IMPLEMENTED] TopLayer()")
	return "", nil
}

func (o *OCI) Entrypoint() ([]string, error) {
	panic("[NOT-IMPLEMENTED] Entrypoint()")
	return []string{}, nil
}

func (o *OCI) ManifestSize() (int64, error) {
	panic("[NOT-IMPLEMENTED] ManifestSize()")
	return 0, nil
}

func (o *OCI) bootstrap() error {
	var err error
	if o.store, err = bootstrapContainerStorage(); err != nil {
		return err
	}
	o.builder, err = buildah.NewBuilder(o.ctx, o.store, buildah.BuilderOptions{
		CommonBuildOpts:  &buildah.CommonBuildOptions{},
		ConfigureNetwork: buildah.NetworkDefault,
		Format:           buildah.OCIv1ImageManifest,
		FromImage:        o.from,
		Isolation:        buildah.IsolationChroot,
		ReportWriter:     os.Stderr,
		SystemContext:    systemContext,
	})
	return err
}

// TODO:
// * re-organize methods, make sure the sequence is well organized;
// * rename "NewOCI" to "NewImage", following the "local" and "remote" implementations;

func NewImage(ctx context.Context, repoName, from string) (*OCI, error) {
	ReInit()
	img := &OCI{
		ctx:      ctx,
		repoName: repoName,
		from:     from,
	}
	return img, img.bootstrap()
}
