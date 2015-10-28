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
	name string
}

func uploadPodcast(r *http.Request) (uploadFile UploadFile, err error) {
	// the FormFile function takes in the POST input id file
	file, header, err := r.FormFile("file")
	if err != nil {
		return fmt.Errorf("error uploading file: %s:", err)
	}

	defer file.Close()

	target := filepath.Join(common.Storage, "uploadedfile")
	log.Debug("addpodcast", "target", target)
	out, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("Unable to create the file for writing. Check your write access privilege")
	}

	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		return fmt.Errorf("error saving file: %s", err)
	}
	log.Debug("addpodcast", "saved", header.Filename)

	uploadFile.name = header.Filename
	return uploadFile
}

func verifyUpload() {

	// open the uploaded file
	file, err := os.Open("./img.png")

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
