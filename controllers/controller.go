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
	"net/url"
	"os"
	"strconv"
	"fmt"

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
	photos_album := &models.PA{}
	photos_album.Photo_id, _ = strconv.Atoi(c.PostForm("photo_id"))
	photos_album.Album_id, _ = strconv.Atoi(c.PostForm("album_id"))
	jsonValues, _ := json.Marshal(photos_album)

	http.Post(host+"/rating/", "application/json", bytes.NewBuffer(jsonValues))

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

	albumResponse := &models.AlbumResponse{}
	errr := utils.GetJson(host+"/albums/", albumResponse)
	if errr != nil {
		log.Println(errr)
	}

	albums := make(map[int]string)
	for _, e := range albumResponse.Result {
		albums[e.Album_id] = e.Title
	}

	//Create map for uuids and base64-templates -> receive each photo individually from server
	photoData := make(map[int]template.URL)
	for _, e := range allPhotos.Result {
		photoData[e.PhotoId] = utils.GetPhotoData(host, e.Uuid)
	}

	fmt.Printf("%#v\n", albums)

	c.HTML(http.StatusOK, "photos.html", gin.H{
		"photos":        photoData,
		"allRollTypes":  allRollTypes.Result,
		"manufacturers": manufacturersId,
		"albums":        albums,
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

	rating := &models.Rating{}
	rating.Photo_id = c.PostForm("photoId")
	rating.Rating = c.PostForm("rating")
	log.Printf(c.PostForm("photoId"))
	jsonValues, _ := json.Marshal(rating)

	http.Post(host+"/rating/", "application/json", bytes.NewBuffer(jsonValues))
	
	filmRoll := &models.FilmRollResponse{}
	err := utils.GetJson(host+"/filmroll/", filmRoll)
	if err != nil {
		log.Println(err)
	}

	images := make(map[string]template.URL)
	for _, e := range filmRoll.Result {
		uuid := e.Uuid

		images[(e.Uuid)] = utils.GetPhotoData(host, uuid)
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
		"images":        images,
	})
}

func OpenAlbums(c *gin.Context) {
	album := &models.AlbumResponse{}
	err := utils.GetJson(host+"/album/", album)
	if err != nil {
		log.Println(err)
	}

	images := make(map[string]template.URL)
	for _, e := range album.Result {

		
		uuid := e.Uuid

		images[(e.Uuid)] = utils.GetPhotoData(host, uuid)
	}

	c.HTML(http.StatusOK, "albums.html", gin.H{
		"album":      album.Result,
		"images":     images,
	})
}

func OpenRollById(c *gin.Context) {
	//Insert rating into DB
	log.Println(c.Params)



	rating := &models.RatingRaw{}
	rating.Photo_id = c.PostForm("photo_id")
	rating.Rating = c.PostForm("rating")
	jsonValues, _ := json.Marshal(rating)

	http.Post(host+"/rating/", "application/json", bytes.NewBuffer(jsonValues))



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

	ratings := make(map[int]float32)
	for _, e := range photosResponse.Result {
		ratings[e.PhotoId] = e.Rating
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
		"ratings":   ratings,
	})
}

func OpenAlbumById(c *gin.Context) {
	//Insert rating into DB
	log.Println(c.Params)
	http.PostForm(host+"/rating/", url.Values{"photoId": {"185"}, "rating": {"3"}})

	//Call backend and map response to struct
	photosResponse := &models.AlbumPhotosResponse{}
	albumId := c.Params.ByName("id")
	err := utils.GetJson(host+"/photo/album/"+albumId, photosResponse)
	if err != nil {
		log.Println(err)
	}

	//Create map for uuids and base64-templates -> receive each photo individually from server
	photoData := make(map[int]template.URL)
	for _, e := range photosResponse.Result {
		photoData[e.PhotoId] = utils.GetPhotoData(host, e.Uuid)
	}

	ratings := make(map[int]float32)
	for _, e := range photosResponse.Result {
		ratings[e.PhotoId] = e.Rating
	}

	//Get FilmRoll Title
	album := &models.SingleAlbumResponse{}
	err = utils.GetJson(host+"/album/"+albumId, album)
	if err != nil {
		log.Println(err)
	}

	c.HTML(http.StatusOK, "albumById.html", gin.H{
		"photos":    photoData,
		"albumTitle": album.Result,
		"ratings":   ratings,
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
	body, _ := ioutil.ReadAll(c.Request.Body)
    println(string(body))
	filmRequest := &models.FilmRollRequest{}
	filmRequest.Title = c.PostForm("title")
	filmRequest.Description = c.PostForm("description")
	stockNameValue, _ := strconv.Atoi(c.Request.FormValue("stockName"))
	filmRequest.Type_Id = stockNameValue
	jsonValues, _ := json.Marshal(filmRequest)

	http.Post(host+"/filmroll/", "application/json", bytes.NewBuffer(jsonValues))
	c.Redirect(http.StatusFound, "/rolls")
}

func CreateAlbum(c *gin.Context) {
	filmRequest := &models.AlbumRequest{}
	filmRequest.Title = c.PostForm("title")
	filmRequest.Description = c.PostForm("description")
	jsonValues, _ := json.Marshal(filmRequest)

	http.Post(host+"/album/", "application/json", bytes.NewBuffer(jsonValues))
	c.Redirect(http.StatusFound, "/albums")
}

func CreateRating(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
    println(string(body))
	rating := &models.Rating{}
	rating.Photo_id = c.PostForm("photoId")
	rating.Rating = c.PostForm("rating")
	log.Printf(c.PostForm("photoId"))
	jsonValues, _ := json.Marshal(rating)

	http.Post(host+"/rating/", "application/json", bytes.NewBuffer(jsonValues))
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

func DeleteAlbum(c *gin.Context) {
	//Delete album
	req, err := http.NewRequest("DELETE", host+"/album/"+c.Params.ByName("id"), nil)
	if err != nil {
		log.Println("Could not create Delete Album Request: ", err)
		return
	}
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Println("Could not send Delete Album Request: ", err)
		return
	}

	c.Redirect(http.StatusFound, "/albums")
}