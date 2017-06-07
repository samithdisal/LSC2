package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rs/xid"
	"malhora.info/lew/crw"
)

func serv(host string) {
	log.Printf("Listening on %s", host)
	http.HandleFunc("/job/submit", jobSubmit)
	http.ListenAndServe(host, nil)
}

func jobSubmit(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	url := params["url"][0]
	typ := params["typ"][0]
	id := scheduleJob(url, typ)
	log.Printf("Getting %s as %s under job id %s", url, typ, id)
	w.WriteHeader(200)
	w.Write([]byte(id))
}

func scheduleJob(url string, typ string) string {
	id := xid.New().String()
	go _jobExec(url, typ, id)
	return id
}

func _jobExec(url string, typ string, id string) {

	os.MkdirAll("."+string(filepath.Separator)+id, 0777)

	basedir := fmt.Sprintf("%s/", id)

	if typ == "a" {
		crw.GetAuthor(url, basedir)
	} else {
		crw.GetPub(url, basedir)
	}
}
