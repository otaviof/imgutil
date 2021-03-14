package oci

//
// TODO: * DRY over local.IDIdentifier;
//

type IDIdentifier struct {
	ImageID string
}

func (i IDIdentifier) String() string {
	return i.ImageID
}
