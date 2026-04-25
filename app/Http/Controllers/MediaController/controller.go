package mediacontroller

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"govibe/app/Http/Response"
	"govibe/configs"

	"cloud.google.com/go/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

type MediaController struct{}

func New() *MediaController {
	return &MediaController{}
}

func (ctl *MediaController) Upload(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil || fileHeader == nil {
		return fiber.NewError(fiber.StatusBadRequest, "missing file")
	}

	cfg, cfgErr := configs.LoadFirebaseConfig()
	if cfgErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, cfgErr.Error())
	}

	folder := strings.TrimSpace(c.FormValue("folder"))
	if folder == "" {
		folder = "uploads"
	}
	folder = prefixedObjectPath(cfg.UploadPrefix, folder)

	originalName := sanitizeFilename(fileHeader.Filename)
	if originalName == "" {
		originalName = "file"
	}

	// Keep uploads stable and collision-free.
	now := time.Now().UTC()
	objectName := fmt.Sprintf("%s/%s-%s-%s", folder, now.Format("20060102T150405Z"), uuid.NewString(), originalName)

	src, openErr := fileHeader.Open()
	if openErr != nil {
		return fiber.NewError(fiber.StatusBadRequest, "failed to open file")
	}
	defer func() { _ = src.Close() }()

	head := make([]byte, 512)
	n, readErr := src.Read(head)
	if readErr != nil && readErr != io.EOF {
		return fiber.NewError(fiber.StatusBadRequest, "failed to read file")
	}
	if n == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "empty file")
	}

	contentType := http.DetectContentType(head[:n])
	if !isAllowedImageContentType(contentType) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "only JPG, PNG, WEBP, or GIF images are allowed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	reader := io.MultiReader(bytes.NewReader(head[:n]), src)
	downloadURL, uploadErr := uploadToFirebaseStorage(ctx, cfg, objectName, contentType, reader)
	if uploadErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, uploadErr.Error())
	}

	return response.Created(c, "uploaded", fiber.Map{
		"bucket":       cfg.StorageBucket,
		"name":         objectName,
		"filename":     originalName,
		"size":         fileHeader.Size,
		"content_type": contentType,
		"download_url": downloadURL,
	})
}

func isAllowedImageContentType(contentType string) bool {
	switch strings.ToLower(strings.TrimSpace(contentType)) {
	case "image/jpeg", "image/png", "image/webp", "image/gif":
		return true
	default:
		return false
	}
}

func prefixedObjectPath(prefix, folder string) string {
	folder = sanitizeObjectPath(folder)
	prefix = sanitizeObjectPath(prefix)
	if prefix == "" || folder == prefix || strings.HasPrefix(folder, prefix+"/") {
		return folder
	}
	return prefix + "/" + folder
}

func uploadToFirebaseStorage(ctx context.Context, cfg configs.FirebaseConfig, objectName, contentType string, r io.Reader) (string, error) {
	opts := make([]option.ClientOption, 0, 1)
	if cfg.CredentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(cfg.CredentialsFile))
	} else {
		opts = append(opts, option.WithCredentialsJSON([]byte(cfg.CredentialsJSONRaw)))
	}

	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return "", err
	}
	defer func() { _ = client.Close() }()

	token := uuid.NewString()

	w := client.Bucket(cfg.StorageBucket).Object(objectName).NewWriter(ctx)
	if strings.TrimSpace(contentType) != "" {
		w.ContentType = contentType
	}
	w.Metadata = map[string]string{
		// Enables token-based download URLs for Firebase Storage.
		"firebaseStorageDownloadTokens": token,
	}

	if _, err := io.Copy(w, r); err != nil {
		_ = w.Close()
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}

	escaped := pathEscapeObjectName(objectName)
	downloadURL := fmt.Sprintf(
		"https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&token=%s",
		cfg.StorageBucket,
		escaped,
		token,
	)
	return downloadURL, nil
}

func sanitizeFilename(name string) string {
	name = strings.TrimSpace(name)
	name = path.Base(strings.ReplaceAll(name, "\\", "/"))
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.Trim(name, ".-_")
	if len(name) > 140 {
		name = name[:140]
	}
	return name
}

func sanitizeObjectPath(p string) string {
	p = strings.TrimSpace(p)
	p = strings.ReplaceAll(p, "\\", "/")
	p = strings.Trim(p, "/")
	p = strings.ReplaceAll(p, "..", "")
	p = strings.ReplaceAll(p, " ", "-")
	p = strings.Trim(p, ".-_")
	if p == "" {
		return "uploads"
	}
	return p
}

func pathEscapeObjectName(name string) string {
	// Firebase download URLs need path-style escaping; using path segments helps keep "/" intact.
	parts := strings.Split(name, "/")
	for i := range parts {
		parts[i] = url.PathEscape(parts[i])
	}
	return strings.Join(parts, "%2F")
}
