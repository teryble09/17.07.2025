package service

import (
	"bytes"
	"errors"
	"mime"
	"net/http"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/teryble09/17.07.2025/internal/archiver/model"
)

var (
	ErrFailedAllRetries   = errors.New("failed all retries")
	ErrNotAllowedMimeType = errors.New("mime type is nopt allowed")
)

func LoadFileAndArchive(srv *TaskService, id model.TaskID, url string) {

}

func Retry(fn func() error, wait time.Duration, maxRetries int) error {
	for range maxRetries {
		err := fn()
		if err == nil {
			return nil
		}
		time.Sleep(wait)
	}
	return ErrFailedAllRetries
}

func checkMIMEType(client http.Client, allowedMIMETypes []string, url string) error {
	resp, err := client.Head(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respMIMEType := resp.Header.Get("Content-Type")
	if slices.Contains(allowedMIMETypes, respMIMEType) {
		return nil
	}
	return ErrNotAllowedMimeType
}

func LoadFile(client http.Client, allowedMIMETypes []string, url string) (filename string, file []byte, err error) {
	err = checkMIMEType(client, allowedMIMETypes, url)
	if err != nil {
		return "", nil, err
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	extensions, err := mime.ExtensionsByType(resp.Header.Get("Content-Type"))
	if err != nil {
		return "", nil, err
	}

	filename = uuid.NewString() + extensions[0]
	buf := bytes.NewBuffer([]byte{})
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", nil, err
	}

	return filename, buf.Bytes(), nil
}
