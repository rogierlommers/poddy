package poddy

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"time"

	"github.com/rogierlommers/poddy/internal/common"
	log "gopkg.in/inconshreveable/log15.v2"
)

// UploadFile describes files in storage
type UploadFile struct {
	Name     string
	Size     int64
	Filetype string
	Failed   bool
	Added    time.Time
}

func uploadPodcast(r *http.Request) (uploadFile UploadFile, err error) {
	file, header, err := r.FormFile("file")
	if err != nil {
		return uploadFile, err
	}
	defer file.Close()

	target := filepath.Join(common.Storage, header.Filename)
	out, err := os.Create(target)
	if err != nil {
		return uploadFile, err
	}
	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		return uploadFile, err
	}

	filetype, invalid := isLegalFileFormat(target)
	if invalid {
		log.Error("invalid filetype detected", "file", target, "filetype", filetype)
		deleteError := os.Remove(target)
		if deleteError != nil {
			log.Error("could not delete invalid file", "file", target, "filetype", filetype)
		}
	}

	fileInfo, _ := out.Stat()
	uploadFile.Name = header.Filename
	uploadFile.Size = fileInfo.Size()
	return uploadFile, nil
}

func isLegalFileFormat(target string) (filetype string, invalid bool) {
	file, err := os.Open(target)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
	buff := make([]byte, 512)
	_, err = file.Read(buff)

	if err != nil {
		log.Error("unable to determine filetype", "file", target, "message", err)
		return "unknown", true
	}

	filetype = http.DetectContentType(buff)
	if filetype != "application/octet-stream" {
		return filetype, true
	}
	return filetype, false
}

// FileList returns slice of UploadedFiles
func FileList() []UploadFile {
	items := []UploadFile{}
	log.Debug("opening dir", "dirname", common.Storage)

	list, err := ioutil.ReadDir(common.Storage)
	if err != nil {
		log.Error("storage error", "error", err)
	}

	for _, file := range list {
		newFile := UploadFile{}
		if file.Mode().IsRegular() {
			newFile.Name = file.Name()
			newFile.Size = file.Size()
		}
		items = append(items, newFile)
	}
	return items
}
