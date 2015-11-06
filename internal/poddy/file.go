package poddy

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rogierlommers/poddy/internal/common"
	log "gopkg.in/inconshreveable/log15.v2"
)

type UploadFile struct {
	name     string
	size     int64
	filetype string
	failed   bool
}

func uploadPodcast(r *http.Request) (uploadFile UploadFile, err error) {
	// the FormFile function takes in the POST input id file
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

	err = verifyUpload(target)

	fileInfo, _ := out.Stat()
	uploadFile.name = header.Filename
	uploadFile.size = fileInfo.Size()
	return uploadFile, nil
}

func verifyUpload(target string) {
	log.Debug("checked filetype", "filename", target)

	// open the uploaded file
	file, err := os.Open(target)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
	buff := make([]byte, 512)
	_, err = file.Read(buff)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	filetype := http.DetectContentType(buff)

	fmt.Println(filetype)

	switch filetype {
	case "image/jpeg", "image/jpg":
		fmt.Println(filetype)

	case "image/gif":
		fmt.Println(filetype)

	case "image/png":
		fmt.Println(filetype)

	case "application/pdf":
		fmt.Println(filetype)
	default:
		fmt.Println("unknown file type uploaded")
	}

}
