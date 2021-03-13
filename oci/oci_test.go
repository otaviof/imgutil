package oci

import (
	"context"
	"testing"

	"github.com/buildpacks/imgutil"
)

const busyboxTag = "quay.io/quay/busybox:latest"

func newImageTest(t *testing.T, repoName, from string) imgutil.Image {
	var img imgutil.Image
	var err error

	if img, err = NewImage(context.TODO(), repoName, from); err != nil {
		t.Fatalf("err='%v'", err)
	}
	return img
}

func TestOCI(t *testing.T) {
	t.Run("Name / OS / Architecture", func(t *testing.T) {
		img := newImageTest(t, "new-image", busyboxTag)

		name := img.Name()
		if name != "new-image" {
			t.Fatalf("Name: invalid image name '%s'", name)
		}

		imgOS, _ := img.OS()
		if imgOS == "" {
			t.Fatal("OS: expected to return a given OS string")
		}

		// TODO: check why it's returning empty;
		// arch, _ := img.Architecture()
		// if arch != "" {
		// 	t.Fatal("Architecture: expected to return a given architecture string")
		// }
	})

	t.Run("SetLabel / Label / Labels / RemoveLabel", func(t *testing.T) {
		img := newImageTest(t, "new-image", busyboxTag)

		err := img.SetLabel("key", "value")
		if err != nil {
			t.Fatalf("SetLabel: err='%v'", err)
		}

		value, err := img.Label("key")
		if err != nil {
			t.Fatalf("Label: err='%v'", err)
		}
		if value != "value" {
			t.Fatal("Label: label value is not correct")
		}

		labels, _ := img.Labels()
		value, ok := labels["key"]
		if value != "value" || !ok {
			t.Fatal("Labels: expect to not retrieve key")
		}

		_ = img.RemoveLabel("key")
		labels, _ = img.Labels()
		_, ok = labels["key"]
		if ok {
			t.Fatal("UnsetLabel: expect not to find lable 'key' anymore")
		}

		if err := img.Save(); err != nil {
			t.Fatalf("Save: err='%v'", err)
		}
	})

	t.Run("SetEnv / Env", func(t *testing.T) {
		img := newImageTest(t, "new-image", busyboxTag)

		_ = img.SetEnv("key", "value")

		value, _ := img.Env("key")
		if value != "value" {
			t.Fatal("Label: label value is not correct")
		}

		if err := img.Save(); err != nil {
			t.Fatalf("Save: err='%v'", err)
		}
	})
}
