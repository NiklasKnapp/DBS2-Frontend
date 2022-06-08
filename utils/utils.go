package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func GetJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return fmt.Errorf("no response from request: %v", err)
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func GetPhotoData(host string, uuid string) template.URL {
	resp, _ := http.Get(host + "/photodata/" + uuid)

	photo, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var base64Photo string
	mimeType := http.DetectContentType(photo)
	switch mimeType {
	case "image/jpeg":
		base64Photo += "data:image/jpeg;base64,"
	case "image/png":
		base64Photo += "data:image/png;base64,"
	}
	base64Photo += base64.StdEncoding.EncodeToString(photo)
	photoLink := template.URL(base64Photo)
	return photoLink
}
