package controllers

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
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
	//Call backend and map response to struct
	photosResponse := &models.FilmRollPhotosResponse{}
	rollId := c.Params.ByName("id")
	err := utils.GetJson(host+"/photo/roll/"+rollId, photosResponse)
	if err != nil {
		log.Println(err)
	}

	//Create map for uuids and base64-templates -> receive each photo individually from server
	photoData := make(map[string]template.URL)
	for _, e := range photosResponse.Result {
		photoData[e.Uuid] = utils.GetPhotoData(host, e.Uuid)
	}

	//Get FilmRoll Title
	filmRoll := &models.SingleFilmRollResponse{}
	err = utils.GetJson(host+"/filmroll/"+rollId, filmRoll)
	if err != nil {
		log.Println(err)
	}
	rollType := &models.RollTypeResponse{}
	path := "/rolltype/" + strconv.Itoa(filmRoll.Result.Type_id)
	err = utils.GetJson(host+path, rollType)
	if err != nil {
		log.Println(err)
	}

	c.HTML(http.StatusOK, "rollById.html", gin.H{
		"photos":    photoData,
		"rollTitle": filmRoll.Result,
		"rollType":  rollType.Result,
	})
}

func ShowPhoto(c *gin.Context) {
	uuid := c.Params.ByName("uuid")
	photoLink := utils.GetPhotoData(host, uuid)

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

func UploadPhotos(c *gin.Context) {
	//Create Buffer
	buf := new(bytes.Buffer)
	bw := multipart.NewWriter(buf)

	formfiles, _ := c.MultipartForm()
	files := formfiles.File["myPhotos"]
	for _, file := range files {
		//Get file from from
		inHeader := file
		in, err := file.Open()
		if err != nil {
			log.Println("Couldn't open FormFile: ", err)
		}
		defer in.Close()

		//Create temp file
		temp, err := ioutil.TempFile("tmp/img/", "*"+inHeader.Filename)
		if err != nil {
			log.Println("Couldn't open parsed File for upload: ", err)
		}
		_, err = io.Copy(temp, in)
		if err != nil {
			log.Println("Couldn't write to temp file: ", err)
		}

		//Reopen tempfile
		output, err := os.Open(temp.Name())
		if err != nil {
			log.Println("Couldn't open parsed File for upload: ", err)
		}

		//Write file to form
		fw, _ := bw.CreateFormFile("files", output.Name())
		io.Copy(fw, output)
	}

	//Write rollId to form once
	tw, _ := bw.CreateFormField("rollId")
	tw.Write([]byte(c.PostForm("rollId")))

	//Close and send form to backend
	bw.Close()
	http.Post(host+"/photo", bw.FormDataContentType(), buf)
	defer os.RemoveAll("tmp/img/")
	c.Redirect(http.StatusFound, "/roll/"+c.PostForm("rollId"))

	// //Get file from form
	// in, inHeader, err := c.Request.FormFile("myPhotos")
	// if err != nil {
	// 	log.Println("Couldn't parse FormFile: ", err)
	// }
	// defer in.Close()

	// //Create tempfile
	// temp, err := ioutil.TempFile("tmp/", "*"+inHeader.Filename)
	// if err != nil {
	// 	log.Println("Couldn't open parsed File for upload: ", err)
	// }
	// _, err = io.Copy(temp, in)
	// if err != nil {
	// 	log.Println("Couldn't write to temp file: ", err)
	// }

	// //Reopen tempfile
	// output, err := os.Open(temp.Name())
	// if err != nil {
	// 	log.Println("Couldn't open parsed File for upload: ", err)
	// }

	// tw, _ := bw.CreateFormField("rollId")
	// tw.Write([]byte(c.PostForm("rollId")))

	// fw, _ := bw.CreateFormFile("files", output.Name())
	// io.Copy(fw, output)

	// bw.Close()
	// http.Post(host+"/photo", bw.FormDataContentType(), buf)
	// c.Redirect(http.StatusFound, "/roll/"+c.PostForm("rollId"))
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
