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
	t.Run("SetLabel/Label", func(t *testing.T) {
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

		if err := img.Save(); err != nil {
			t.Fatalf("Save: err='%v'", err)
		}
	})

	t.Run("SetEnv/Env", func(t *testing.T) {
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
