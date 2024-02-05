package ttp

import (
	"crypto/x509/pkix"
	"errors"
)

func findCertExtensions(extensions []pkix.Extension, label []int) (*pkix.Extension, error) {
	for _, ext := range extensions {
		if ext.Id.Equal(label) {
			return &ext, nil
		}
	}

	return nil, errors.New("extension not found")
}
