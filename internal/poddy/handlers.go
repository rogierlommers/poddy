package poddy

import (
	"html/template"
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/dustin/go-humanize"
	"github.com/rogierlommers/poddy/internal/common"
	log "gopkg.in/inconshreveable/log15.v2"
)

var staticBox *rice.Box

func IndexPage(w http.ResponseWriter, r *http.Request) {
	renderObject := map[string]interface{}{
		"IsLandingPage": "true",
		"buildversion":  common.BuildDate,
	}
	displayPage(w, r, renderObject)
}

func AddPodcast(w http.ResponseWriter, r *http.Request) {
	uploadedFile, err := uploadPodcast(r)
	if err != nil {
		log.Warn("error uploading/saving podcast", "message", err)
		uploadedFile.failed = true
	} else {
		log.Info("file succesfully uploaded", "filename", uploadedFile.name, "size", humanize.Bytes(uint64(uploadedFile.size)))
	}

	renderObject := map[string]interface{}{
		"IsConfirmationPage": "true",
		"failed":             uploadedFile.failed,
		"name":               uploadedFile.name,
		"size":               uploadedFile.size,
	}
	displayPage(w, r, renderObject)
}

func CreateStaticBox() {
	// create rice.box with static files
	staticBox = rice.MustFindBox("../../static")

	// static files are being exposed through /static endpoint
	staticFileServer := http.StripPrefix("/static/", http.FileServer(staticBox.HTTPBox()))
	http.Handle("/static/", staticFileServer)
}

func displayPage(w http.ResponseWriter, r *http.Request, dynamicData interface{}) {
	templateString, err := staticBox.String("index.tmpl")
	if err != nil {
		log.Crit("template", "error", err)
	}

	tmplMessage, err := template.New("messsage").Parse(templateString)
	if err != nil {
		log.Crit("template", "error", err)
	}

	tmplMessage.Execute(w, dynamicData)
}
