package configs

import (
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type FirebaseConfig struct {
	StorageBucket      string
	UploadPrefix       string
	CredentialsFile    string
	CredentialsJSONRaw string
}

func LoadFirebaseConfig() (FirebaseConfig, error) {
	_ = godotenv.Load()

	bucket := strings.TrimSpace(os.Getenv("FIREBASE_STORAGE_BUCKET"))
	uploadPrefix := strings.TrimSpace(os.Getenv("FIREBASE_UPLOAD_PREFIX"))
	credsFile := strings.TrimSpace(os.Getenv("FIREBASE_CREDENTIALS_FILE"))
	credsJSON := strings.TrimSpace(os.Getenv("FIREBASE_CREDENTIALS_JSON"))

	if uploadPrefix == "" {
		uploadPrefix = "Govibe"
	}
	if bucket == "" {
		return FirebaseConfig{}, errors.New("missing FIREBASE_STORAGE_BUCKET; set it to your Firebase Storage bucket name")
	}
	if credsFile == "" && credsJSON == "" && fileExists("firebase-adminsdk.json") {
		credsFile = "firebase-adminsdk.json"
	}
	if credsFile == "" && credsJSON == "" {
		return FirebaseConfig{}, errors.New("missing FIREBASE_CREDENTIALS_FILE or FIREBASE_CREDENTIALS_JSON")
	}

	return FirebaseConfig{
		StorageBucket:      bucket,
		UploadPrefix:       uploadPrefix,
		CredentialsFile:    credsFile,
		CredentialsJSONRaw: credsJSON,
	}, nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
