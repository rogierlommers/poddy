package poddy

import (
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
	//files := FileList()

	c := gopod.ChannelFactory("My personal channel", "http://RubyDeveloper.com/", "My Blog", "http://example.com/image.png")
	c.SetPubDate(time.Now().UTC())
	c.SetiTunesExplicit("No")

	c.AddItem(&gopod.Item{
		Title:       "Stack Overflow",
		Link:        "http://stackoverflow.com",
		Description: "Stack Overflow",
		PubDate:     time.Now().UTC().Format(time.RFC1123),
	})

	// Example: Using an item's methods
	t := "My title"
	l := "http://linkedin.com"
	i := &gopod.Item{
		Title:         t,
		TunesSubtitle: t,
		Link:          l,
		Description:   "My LinkedIn",
		TunesDuration: "600",
		TunesSummary:  "I asked myself that question more than a decade ago and it changed my...",
		Guid:          l,
		Creator:       "Daniel's Channel",
	}
	i.SetEnclosure("http://example.com/sound.mp3", "600", "audio/mpeg")
	i.SetPubDate(time.Now().Unix())
	c.AddItem(i)

	//	for _, file := range files {
	//		// link := fmt.Sprintf("%s/download/%s", "http://poddy.lommers.org", file.Name)
	//
	//	}

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
