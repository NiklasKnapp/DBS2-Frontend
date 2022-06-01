package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"de.stuttgart.hft/DBS2-Frontend/models"
	"de.stuttgart.hft/DBS2-Frontend/utils"
	"github.com/gin-gonic/gin"
)

var (
	host string = "http://localhost:8080"
)

func RedirectToIndex(c *gin.Context) {
	c.Redirect(http.StatusFound, "/rolls")
}

func OpenPhotos(c *gin.Context) {
	// resp, err := http.Get(host + "/filmroll/")
	// if err != nil {
	// 	log.Println("No response from request: ", err)
	// }
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Println("Could not read response: ", err)
	// }

	// var filmRoll models.FilmRollResponse
	// if err := json.Unmarshal(body, &filmRoll); err != nil {
	// 	log.Println("Can not unmarshal JSON: ", err)
	// }
	// fmt.Printf("%#v\n", filmRoll.Result)
	c.HTML(http.StatusOK, "photos.html", gin.H{
		"title": "Main website",
	})
}

func OpenRolls(c *gin.Context) {
	filmRoll := &models.FilmRollResponse{}
	err := utils.GetJson(host+"/filmroll/", filmRoll)
	if err != nil {
		log.Println(err)
	}

	typeids := make(map[int]string)
	for _, e := range filmRoll.Result {
		if _, ok := typeids[(e.Type_id)]; ok {
			continue
		} else {
			rollType := &models.RollTypeResponse{}
			path := "/rolltype/" + strconv.Itoa(e.Type_id)
			err := utils.GetJson(host+path, rollType)
			if err != nil {
				log.Println(err)
			}
			typeids[(e.Type_id)] = rollType.Result.StockName
		}
	}

	allRollTypes := &models.MultipleRollTypeResponse{}
	err = utils.GetJson(host+"/rolltype/", allRollTypes)
	if err != nil {
		log.Println(err)
	}

	// fmt.Printf("%#v\n", allRollTypes.Result)

	c.HTML(http.StatusOK, "rolls.html", gin.H{
		"filmRoll":     filmRoll.Result,
		"types":        typeids,
		"allRollTypes": allRollTypes.Result,
	})
}

func OpenRollById(c *gin.Context) {

}

func ShowPhoto(c *gin.Context) {
	uuid := c.Params.ByName("uuid")
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

	c.HTML(http.StatusOK, "rollById.html", gin.H{
		"photo": photoLink,
	})

	// c.Data(http.StatusOK, "image/jpeg", photo)

	// c.HTML(http.StatusOK, "rollById.html", gin.H{
	// 	"photo": resp,
	// })

	// uuid := c.Params.ByName("uuid")
	// resp, _ := http.Get(host + "/photodata/" + uuid)

	// photo, _ := ioutil.ReadAll(resp.Body)
	// resp.Body.Close()

	// fmt.Printf("%#v\n", photo)

	// var base64Photo string
	// mimeType := http.DetectContentType(photo)
	// switch mimeType {
	// case "image/jpeg":
	// 	base64Photo += "data:image/jpeg;base64,"
	// case "image/png":
	// 	base64Photo += "data:image/png;base64,"
	// }
	// base64Photo += base64.StdEncoding.EncodeToString(photo)
	// fmt.Printf("%#v\n", base64Photo)

	// c.HTML(http.StatusOK, "rollById.html", gin.H{
	// 	"photo": photo,
	// })
}

func CreateRoll(c *gin.Context) {
	filmRequest := &models.FilmRollRequest{}
	filmRequest.Title = c.PostForm("title")
	filmRequest.Description = c.PostForm("description")
	stockNameValue, _ := strconv.Atoi(c.Request.FormValue("stockName"))
	filmRequest.Type_Id = stockNameValue
	jsonValues, _ := json.Marshal(filmRequest)

	http.Post(host+"/filmroll/", "application/json", bytes.NewBuffer(jsonValues))
	c.Redirect(http.StatusFound, "/rolls")
}
