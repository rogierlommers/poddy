package poddy

import (
	"html/template"
	"net/http"
	"time"

	"fmt"

	"github.com/GeertJohan/go.rice"
	"github.com/dustin/go-humanize"
	"github.com/gorilla/feeds"
	"github.com/rogierlommers/poddy/internal/common"
	log "gopkg.in/inconshreveable/log15.v2"
)

var staticBox *rice.Box

func IndexPage(w http.ResponseWriter, r *http.Request) {
	renderObject := map[string]interface{}{
		"IsLandingPage": "true",
		"buildversion":  common.BuildDate,
		"Filelist":      FileList(),
	}
	displayPage(w, r, renderObject)
}

func AddPodcast(w http.ResponseWriter, r *http.Request) {
	uploadedFile, err := uploadPodcast(r)
	if err != nil {
		log.Warn("error uploading/saving podcast", "message", err)
		uploadedFile.Failed = true
	} else {
		log.Info("file succesfully uploaded", "filename", uploadedFile.Name, "size", humanize.Bytes(uint64(uploadedFile.Size)))
	}

	renderObject := map[string]interface{}{
		"IsConfirmationPage": "true",
		"failed":             uploadedFile.Failed,
		"name":               uploadedFile.Name,
		"size":               uploadedFile.Size,
	}
	displayPage(w, r, renderObject)
}

func Feed(w http.ResponseWriter, r *http.Request) {
	files := FileList()
	now := time.Now()

	feed := &feeds.Feed{
		Title:       "my poddy feed",
		Link:        &feeds.Link{Href: "http://poddy.lommers.org"},
		Description: "My saved podcasts",
		Author:      &feeds.Author{"dummy", "dummy"},
		Created:     now,
	}

	for _, file := range files {
		link := fmt.Sprintf("%s/download/%s", "http://poddy.lommers.org", file.Name)
		newItem := feeds.Item{
			Title: file.Name,
			Link:  &feeds.Link{Href: link},
		}
		feed.Add(&newItem)
	}

	rss, err := feed.ToRss()
	if err != nil {
		log.Error("error generation RSS feed", "message", err)
		return
	}
	w.Write([]byte(rss))

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
