package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/forj-oss/goforjj"
)

// RESTReaderLimit is the maximun supported REST API request size.
const RESTReaderLimit = 65535

// PluginData response object creator
func newPluginData() *goforjj.PluginData {
	r := goforjj.PluginData{
		Repos:   make(map[string]goforjj.PluginRepo),
		Options: make(map[string]goforjj.PluginOption),
		Services: goforjj.PluginService{
			Urls: make(map[string]string),
		},
	}
	return &r
}

// Function to detect header content-type matching
// return true if match
func contentTypeMatch(header http.Header, match string) (string, bool) {
	contentType, found := header["Content-Type"]
	if !found {
		return "Header 'Content-Type' missing.", false
	}

	for _, v := range contentType {
		if v == match {
			return "", true
		}
	}
	return strings.Join(contentType, ", "), false
}

func panicIfError(w http.ResponseWriter, err error, message string, pars ...interface{}) {
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if message != "" {
			err = fmt.Errorf("%s %s", fmt.Errorf(message, pars...), err)
		}
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
}

func requestResponse(w http.ResponseWriter, data *goforjj.PluginData, code int) {
	if data.ErrorMessage != "" {
		if code == 0 {
			code = 422 // unprocessable entity
		}
		log.Print("HTTP ERROR: ", code, " - ", data.ErrorMessage)
	} else {
		code = 200
	}
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}

}

func doPluginAction(w http.ResponseWriter, r *http.Request,
	requestUnmarshal func([]byte) error,
	requestDo func(*goforjj.PluginData) int) {

	data := newPluginData()
	var errCode int

	// Respond to the request in json format except if fatal
	defer requestResponse(w, data, errCode)

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, RESTReaderLimit))

	if err != nil {
		data.Errorf("Unable to read request buffer. %s", err)
		return
	}

	if contentType, found := contentTypeMatch(r.Header, "application/json"); !found {
		data.Errorf("Invalid request payload format. Must be 'application/json'. Got %s", contentType)
		return
	}

	err = requestUnmarshal(body)
	if err != nil {
		if err.Error() == "unexpected end of JSON input" && len(body) == RESTReaderLimit {
			err = fmt.Errorf("%s. The plugin reach REST API request buffer limit (%d)", err, RESTReaderLimit)
		}
		data.Errorf("Unable to decode json. %s", err)
		return
	}

	errCode = requestDo(data)

	// defer will respond, with data and errCode
}

// request Handlers

// Create handler
func Create(w http.ResponseWriter, r *http.Request) {
	var reqData CreateReq

	doPluginAction(w, r,
		func(body []byte) error {
			return json.Unmarshal(body, &reqData)
		},
		func(data *goforjj.PluginData) (errCode int) {
			errCode = DoCreate(r, &reqData, data)

			reqData.Objects.SaveMaintainOptions(data)
			return
		})
}

// Update handler
func Update(w http.ResponseWriter, r *http.Request) {
	var reqData UpdateReq

	doPluginAction(w, r,
		func(body []byte) error {
			return json.Unmarshal(body, &reqData)
		},
		func(data *goforjj.PluginData) (errCode int) {
			errCode = DoUpdate(r, &reqData, data)

			reqData.Objects.SaveMaintainOptions(data)
			return
		})
}

// Maintain handler
func Maintain(w http.ResponseWriter, r *http.Request) {
	var reqData MaintainReq

	doPluginAction(w, r,
		func(body []byte) error {
			return json.Unmarshal(body, &reqData)
		},
		func(data *goforjj.PluginData) int {
			return DoMaintain(r, &reqData, data)
		})
}

// Index Handler
//
func Index(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "FORJJ - jenkins driver for FORJJ. It is Implemented as a REST API.")
}

// Quit handler
func Quit(w http.ResponseWriter, _ *http.Request) {
	goforjj.DefaultQuit(w, "")
}
