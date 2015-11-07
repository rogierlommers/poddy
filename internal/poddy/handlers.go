package poddy

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/dustin/go-humanize"
	"github.com/rogierlommers/poddy/internal/common"
	log "gopkg.in/inconshreveable/log15.v2"
)

var staticBox *rice.Box

func IndexPage(w http.ResponseWriter, r *http.Request) {
	//	podcastList, err := listPodcasts()
	//	if err != nil {
	//		log.Error(err.Error())
	//	}
	//	spew.Dump(podcastList)

	renderObject := map[string]interface{}{
		"IsLandingPage": "true",
		"buildversion":  common.BuildDate,
	}
	displayPage(w, r, renderObject)
}

func AddPodcast(w http.ResponseWriter, r *http.Request) {
	var humanSize string
	uploadedFile, err := uploadPodcast(r)

	if err != nil {
		uploadedFile.failed = true
		log.Warn("error uploading/saving podcast", "message", err)
	} else {
		humanSize = humanize.Bytes(uint64(uploadedFile.size))
		log.Info("file succesfully uploaded", "filename", uploadedFile.name, "size", humanSize)
	}

	renderObject := map[string]interface{}{
		"IsConfirmationPage": "true",
		"failed":             uploadedFile.failed,
		"name":               uploadedFile.name,
		"size":               humanSize,
		"errormessage":       uploadedFile.errormessage,
	}
	displayPage(w, r, renderObject)
}

func DisplayPodcastFeed(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(common.Storage)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}

	w.Write([]byte("rogk"))
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
