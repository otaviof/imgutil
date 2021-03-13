module github.com/buildpacks/imgutil

require (
	github.com/containers/buildah v1.19.6
	github.com/containers/image/v5 v5.10.4
	github.com/containers/storage v1.27.0
	github.com/docker/docker v17.12.0-ce-rc1.0.20201020191947-73dc6a680cdd+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/google/go-cmp v0.5.4
	github.com/google/go-containerregistry v0.4.0
	github.com/pkg/errors v0.9.1
	github.com/sclevine/spec v1.4.0
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
)

// replace golang.org/x/sys => golang.org/x/sys v0.0.0-20200523222454-059865788121
replace github.com/docker/docker => github.com/docker/docker v1.4.2-0.20190924003213-a8608b5b67c7

go 1.14
