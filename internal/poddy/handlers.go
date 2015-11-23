package poddy

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/andjosh/gopod"
	"github.com/dustin/go-humanize"
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
	log.Debug("crawler passde by", "user-agent", r.Header.Get("User-Agent"))
	files := FileList()
	c := gopod.ChannelFactory("Poddy!", common.Self, "Personal poddy feed...", fmt.Sprintf("%s/static/images/poddy.png", common.Self))
	c.SetPubDate(time.Now().UTC())
	c.SetiTunesExplicit("No")

	for _, file := range files {
		link := fmt.Sprintf("%s/download/%s", "http://poddy.lommers.org", file.Name)
		i := &gopod.Item{
			Title:         file.Name,
			TunesSubtitle: file.Name,
			Link:          link,
			Description:   file.Name,
			Guid:          link,
			Creator:       "Rogier",
		}

		i.SetEnclosure(link, "1", file.Filetype)
		i.SetPubDate(time.Now().Unix())
		c.AddItem(i)

	}

	feed := c.Publish()
	w.Write([]byte(feed))

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
