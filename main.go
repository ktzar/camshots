package main

import (
	_ "encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

var addr = flag.Int("port", 8082, "port to listen to")

const MAX_RECORDS = 48

func main() {
	flag.Parse()

	http.Handle("/shot", http.HandlerFunc(shot))
	http.Handle("/last", http.HandlerFunc(last))
	http.Handle("/latest", http.HandlerFunc(latest))

	port := strconv.Itoa(*addr)
	fmt.Println("Listening in " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		return
	}
}

func shot(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("cache-control", "max-age=31536000")
	w.Header().Set("content-type", "image/jpeg")
	file := req.FormValue("file")

	r, err := regexp.Compile("[a-z0-9_]+\\.jpg")
	if err != nil {
		fmt.Println(err)
	}
	matched := r.MatchString(file)
	if !matched {
		http.Error(w,
			"Invalid file "+file,
			http.StatusNotFound,
		)
		return
	}
	dumpImage(w, file)
}

func latest(w http.ResponseWriter, req *http.Request) {
	addImageHeader(w)

	shots := getShots()
	dumpImage(w, shots[len(shots)-1])
}

func last(w http.ResponseWriter, req *http.Request) {
	addHtmlHeaders(w)
	shots := getShots()
	w.Write([]byte("<html>"))

	// List in reverse
	for i := len(shots) - 1; i > len(shots)-MAX_RECORDS; i-- {
		url := "/shot?file=" + shots[i]
		line := []byte(
			"<a href='" + url + "'>" +
				"<img width='320px' src='" + url + "' title='" + shots[i] + "'/></a>")
		w.Write(line)
	}
	w.Write([]byte("</html>"))
}

func getShots() []string {
	files, err := ioutil.ReadDir("../home_shots/shots")
	if err != nil {
		log.Fatal(err)
	}

	filenames := []string{}

	for _, file := range files {
		filenames = append(filenames, file.Name())
	}

	return filenames
}

func addImageHeader(w http.ResponseWriter) {
	w.Header().Set("cache-control", "max-age=31536000")
	w.Header().Set("content-type", "image/jpeg")
}

func addHtmlHeaders(w http.ResponseWriter) {
	w.Header().Set("cache-control", "private, max-age=0, no-cache")
	w.Header().Set("content-type", "text/html")
}

func dumpImage(w http.ResponseWriter, file string) {
	filename := "../home_shots/shots/" + file
	shot_file, err := ioutil.ReadFile(filename)

	if err != nil {
		http.Error(w,
			"File not found",
			http.StatusNotFound,
		)
		return
	}

	w.Write(shot_file)
}
