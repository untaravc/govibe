package configs

import (
	"errors"
	"os"
	"strings"
)

type FirebaseConfig struct {
	StorageBucket      string
	CredentialsFile    string
	CredentialsJSONRaw string
}

func LoadFirebaseConfig() (FirebaseConfig, error) {
	bucket := strings.TrimSpace(os.Getenv("FIREBASE_STORAGE_BUCKET"))
	credsFile := strings.TrimSpace(os.Getenv("FIREBASE_CREDENTIALS_FILE"))
	credsJSON := strings.TrimSpace(os.Getenv("FIREBASE_CREDENTIALS_JSON"))

	if bucket == "" {
		return FirebaseConfig{}, errors.New("missing FIREBASE_STORAGE_BUCKET")
	}
	if credsFile == "" && credsJSON == "" {
		return FirebaseConfig{}, errors.New("missing FIREBASE_CREDENTIALS_FILE or FIREBASE_CREDENTIALS_JSON")
	}

	return FirebaseConfig{
		StorageBucket:      bucket,
		CredentialsFile:    credsFile,
		CredentialsJSONRaw: credsJSON,
	}, nil
}
