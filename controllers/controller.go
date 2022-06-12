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
	//Get Roll Type names
	allRollTypes := &models.MultipleRollTypeResponse{}
	err := utils.GetJson(host+"/rolltype/", allRollTypes)
	if err != nil {
		log.Println(err)
	}

	//Get manufacturer names
	manufacturersId := make(map[int]string)
	for _, e := range allRollTypes.Result {
		if _, ok := manufacturersId[(e.Type_id)]; ok {
			continue
		} else {
			manufacturer := &models.ManufacturerResponse{}
			err := utils.GetJson(host+"/manufacturer/"+strconv.Itoa(e.M_id), manufacturer)
			if err != nil {
				log.Println(err)
			}
			manufacturersId[e.Type_id] = manufacturer.Result.Name
		}
	}

	//set allPhotos to empty
	allPhotos := &models.FilmRollPhotosResponse{}
	err = utils.GetJson(host+"/photo/type/-2", allPhotos)
	if err != nil {
		log.Println(err)
	}

	//Create map for uuids and base64-templates -> receive each photo individually from server
	photoData := make(map[int]template.URL)
	for _, e := range allPhotos.Result {
		photoData[e.PhotoId] = utils.GetPhotoData(host, e.Uuid)
	}

	c.HTML(http.StatusOK, "photos.html", gin.H{
		"photos":        photoData,
		"allRollTypes":  allRollTypes.Result,
		"manufacturers": manufacturersId,
	})
}

func OpenPhotosByTypeId(c *gin.Context) {
	typeId := c.Request.FormValue("stockName")

	//Get Roll Type names
	allRollTypes := &models.MultipleRollTypeResponse{}
	err := utils.GetJson(host+"/rolltype/", allRollTypes)
	if err != nil {
		log.Println(err)
	}

	//Get manufacturer names
	manufacturersId := make(map[int]string)
	for _, e := range allRollTypes.Result {
		if _, ok := manufacturersId[(e.Type_id)]; ok {
			continue
		} else {
			manufacturer := &models.ManufacturerResponse{}
			err := utils.GetJson(host+"/manufacturer/"+strconv.Itoa(e.M_id), manufacturer)
			if err != nil {
				log.Println(err)
			}
			manufacturersId[e.Type_id] = manufacturer.Result.Name
		}
	}

	//Get photos from database
	allPhotos := &models.FilmRollPhotosResponse{}
	if typeId == "-2" {
		//Select Roll Type -> return empty
		err = utils.GetJson(host+"/photo/type/-2", allPhotos)
		if err != nil {
			log.Println(err)
		}
	} else if typeId == "-1" {
		//Select all -> return all
		err = utils.GetJson(host+"/photo/", allPhotos)
		if err != nil {
			log.Println(err)
		}
	} else {
		//Select specific -> return specific type id
		err = utils.GetJson(host+"/photo/type/"+typeId, allPhotos)
		if err != nil {
			log.Println(err)
		}
	}

	//Create map for uuids and base64-templates -> receive each photo individually from server
	photoData := make(map[int]template.URL)
	rollIdMap := make(map[int]int)
	for _, e := range allPhotos.Result {
		photoData[e.PhotoId] = utils.GetPhotoData(host, e.Uuid)
		rollIdMap[e.PhotoId] = e.RollId
	}

	//fmt.Printf("%#v\n", rollIdMap)

	c.HTML(http.StatusOK, "photos.html", gin.H{
		"photos":        photoData,
		"allRollTypes":  allRollTypes.Result,
		"manufacturers": manufacturersId,
		"rollIdMap":     rollIdMap,
	})
}

func OpenRolls(c *gin.Context) {
	filmRoll := &models.FilmRollResponse{}
	err := utils.GetJson(host+"/filmroll/", filmRoll)
	if err != nil {
		log.Println(err)
	}

	//Get Roll Types and stock name
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

	//Get Roll Type names
	allRollTypes := &models.MultipleRollTypeResponse{}
	err = utils.GetJson(host+"/rolltype/", allRollTypes)
	if err != nil {
		log.Println(err)
	}

	//Get manufacturer names
	manufacturersId := make(map[int]string)
	for _, e := range allRollTypes.Result {
		if _, ok := manufacturersId[(e.Type_id)]; ok {
			continue
		} else {
			manufacturer := &models.ManufacturerResponse{}
			err := utils.GetJson(host+"/manufacturer/"+strconv.Itoa(e.M_id), manufacturer)
			if err != nil {
				log.Println(err)
			}
			manufacturersId[e.Type_id] = manufacturer.Result.Name
		}
	}

	c.HTML(http.StatusOK, "rolls.html", gin.H{
		"filmRoll":      filmRoll.Result,
		"types":         typeids,
		"allRollTypes":  allRollTypes.Result,
		"manufacturers": manufacturersId,
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
	photoData := make(map[int]template.URL)
	for _, e := range photosResponse.Result {
		photoData[e.PhotoId] = utils.GetPhotoData(host, e.Uuid)
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
}

func UploadPhotos(c *gin.Context) {
	//Create Buffer
	buf := new(bytes.Buffer)
	bw := multipart.NewWriter(buf)

	os.Mkdir("tmp/img/", 0777)

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

func DeleteSinglePhoto(c *gin.Context) {
	req, err := http.NewRequest("DELETE", host+"/photo/"+c.Params.ByName("id"), nil)
	if err != nil {
		log.Println("Could not create Delete Photo Request: ", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Could not send Delete Photo Request: ", err)
		return
	}
	defer resp.Body.Close()

	response := &models.PhotoDeletedResponse{}
	json.NewDecoder(resp.Body).Decode(response)
	rollId := response.Result.RollId

	c.Redirect(http.StatusFound, "/roll/"+strconv.Itoa(rollId))
}

func DeletePhotoFromPool(c *gin.Context) {
	req, err := http.NewRequest("DELETE", host+"/photo/"+c.Params.ByName("id"), nil)
	if err != nil {
		log.Println("Could not create Delete Photo Request: ", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Could not send Delete Photo Request: ", err)
		return
	}
	defer resp.Body.Close()

	c.Redirect(http.StatusFound, "/photos")
}

func DeleteRollAndPhotos(c *gin.Context) {
	photosInRoll := &models.FilmRollPhotosResponse{}
	utils.GetJson(host+"/photo/roll/"+c.Params.ByName("id"), photosInRoll)

	//Delete photos in roll
	for _, e := range photosInRoll.Result {
		req, err := http.NewRequest("DELETE", host+"/photo/"+strconv.Itoa(e.PhotoId), nil)
		if err != nil {
			log.Println("Could not create Delete Photo Request: ", err)
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Could not send Delete Photo Request: ", err)
			return
		}
		defer resp.Body.Close()
	}

	//Delete roll
	req, err := http.NewRequest("DELETE", host+"/filmroll/"+c.Params.ByName("id"), nil)
	if err != nil {
		log.Println("Could not create Delete Roll Request: ", err)
		return
	}
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Println("Could not send Delete Roll Request: ", err)
		return
	}

	c.Redirect(http.StatusFound, "/rolls")
}
