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
	ErrNotAllowedMimeType = errors.New("mime type is not allowed")
	ErrNoMimeType         = errors.New("no Content-Type set")
)

func LoadFileAndArchive(srv *TaskService, id model.TaskID, url string) {
	file := []byte{}
	filename := ""

	fn := func() (stop bool, err error) {
		client := http.Client{Timeout: time.Duration(srv.Cfg.HttpClientTimeout * 1_000_000)}
		filename, file, err = LoadFile(client, srv.Cfg.AllowedMIMETypes, url)
		if err == ErrNotAllowedMimeType {
			return true, ErrNotAllowedMimeType
		}
		if err == ErrNoMimeType {
			return true, ErrNoMimeType
		}
		if err != nil {
			return false, err
		}
		return true, nil
	}

	err := Retry(fn, time.Duration(srv.Cfg.RetryWaitTime*1_000_000), srv.Cfg.MaxRetryAmount)

	if err == ErrNotAllowedMimeType {
		srv.Storage.ChangeStatus(id, url, model.NotAllowedType)
		archiveFinished, _ := srv.Storage.EmptyWriteToArchive(id)
		srv.Logger.Warn("Trying to load file with not allowed type", "url", url)
		if archiveFinished {
			srv.Semaphore.Release(1)
		}
		return
	}

	if err == ErrFailedAllRetries || err == ErrNoMimeType {
		srv.Storage.ChangeStatus(id, url, model.FailedToLoad)
		archiveFinished, _ := srv.Storage.EmptyWriteToArchive(id)
		srv.Logger.Warn("Failed to load file", "url", url)
		if archiveFinished {
			srv.Semaphore.Release(1)
		}
		return
	}

	archiveFinished, err := srv.Storage.WriteToArchive(id, filename, file)

	if err != nil {
		srv.Logger.Error("Failed to write loaded file into archive", "url", url, "error", err.Error())
		srv.Storage.ChangeStatus(id, url, model.FailedToArchive)
		return
	}

	if archiveFinished {
		srv.Semaphore.Release(1)
	}

	srv.Logger.Info("Suceesfully loaded and wrote file", "url", url)
	srv.Storage.ChangeStatus(id, url, model.Archived)
}

// fn should return false, if retries should continue, and true if they should stop,
// and error either err from fn func if retries stopped forcefully or ErrFailedAllRetries
func Retry(fn func() (stop bool, err error), wait time.Duration, maxRetries int) error {
	for range maxRetries {
		stop, err := fn()
		if stop {
			return err
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
	if respMIMEType == "" {
		return ErrNoMimeType
	}
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
