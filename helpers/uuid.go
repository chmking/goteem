package helpers

import (
	"encoding/base64"

	"github.com/google/uuid"
)

// MustUUID generates a URL safe base64 encoded UUID and panics on failure.
func MustUUID() string {
	uuid, err := UUID()
	if err != nil {
		panic(err)
	}
	return uuid
}

// UUID generates a URL safe base64 encoded UUID.
func UUID() (string, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	enc := base64.RawURLEncoding.EncodeToString([]byte(uid.String()))
	return enc, nil
}
