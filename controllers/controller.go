package controllers

import (
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

	c.HTML(http.StatusOK, "rolls.html", gin.H{
		"filmRoll": filmRoll.Result,
		"types":    typeids,
	})
}
