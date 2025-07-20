package service

import (
	"bytes"
	"errors"
	"mime"
	"net/http"
	"regexp"
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
	file := []byte{}
	filename := ""

	fn := func() error {
		var err error
		client := http.Client{Timeout: 5 * time.Second}
		filename, file, err = LoadFile(client, srv.Cfg.AllowedMIMETypes, url)
		if err == ErrNotAllowedMimeType {
			srv.Logger.Warn("Not allowed type", "utl", url)
			srv.Storage.ChangeStatus(id, url, model.NotAllowedType)
			return ErrNotAllowedMimeType
		}
		if err != nil {
			return err
		}
		return nil
	}

	err := Retry(fn, time.Duration(srv.Cfg.RetryWaitTime*1_000_000), srv.Cfg.MaxRetryAmount)

	if err == ErrFailedAllRetries {
		srv.Logger.Warn("Failed to load file", "url", url)
		srv.Storage.ChangeStatus(id, url, model.FailedToLoad)
		return
	}

	archiveFinished, err := srv.Storage.WriteToArchive(id, filename, file)

	if err != nil {
		srv.Logger.Error("Failed to write loaded file", "url", url, "error", err.Error())
		srv.Storage.ChangeStatus(id, url, model.FailedToLoad)
		return
	}

	if archiveFinished {
		srv.Semaphore.Release(1)
	}

	srv.Logger.Info("Suceesfully loaded and wrote file", "url", url)
	srv.Storage.ChangeStatus(id, url, model.Archived)

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

	filename = sanitizeFilename(uuid.NewString()) + extensions[0]
	buf := bytes.NewBuffer([]byte{})
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", nil, err
	}

	return filename, buf.Bytes(), nil
}

func sanitizeFilename(name string) string {
	re := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)
	name = re.ReplaceAllString(name, "_")

	if len(name) > 200 {
		name = name[:200]
	}
	return name
}
