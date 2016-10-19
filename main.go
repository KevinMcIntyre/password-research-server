package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/KevinMcIntyre/password-research-server/jobs"
	"github.com/KevinMcIntyre/password-research-server/models"
	"github.com/KevinMcIntyre/password-research-server/services"
	"github.com/KevinMcIntyre/password-research-server/utils"
	"github.com/codegangsta/negroni"
	"github.com/flosch/pongo2"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

var db *sql.DB = setupDatabase()
var config Config = ReadConfig("./app-config.toml")

func uploadImageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var (
		status int
		err    error
	)

	defer func() {
		if nil != err {
			http.Error(w, err.Error(), status)
		}
	}()

	// parse request
	const _24K = (1 << 20) * 24
	if err = r.ParseMultipartForm(_24K); nil != err {
		fmt.Println("Err 1")
		status = http.StatusInternalServerError
		return
	}

	uploadResponse := new(models.ImageUploadResponse)

	uploadResponse.SubjectId, err = strconv.Atoi(r.MultipartForm.Value["subjectId"][0])
	uploadResponse.CollectionId, err = strconv.Atoi(r.MultipartForm.Value["collectionId"][0])

	for _, fileHeaders := range r.MultipartForm.File {
		for _, header := range fileHeaders {
			// open uploaded
			var inFile multipart.File
			if inFile, err = header.Open(); nil != err {
				fmt.Println("Err 2")
				status = http.StatusInternalServerError
				return
			}

			// open destination
			var outFile *os.File
			fileName := "./upload/" + header.Filename
			if outFile, err = os.Create(fileName); nil != err {
				fmt.Println("Err 3")
				status = http.StatusInternalServerError
				return
			}

			// 32K buffer copy
			if _, err = io.Copy(outFile, inFile); nil != err {
				fmt.Println("Err 4")
				status = http.StatusInternalServerError
				return
			}

			utils.ResizeImage(fileName)

			imageBytes := utils.ImageFileToByteArray(fileName)

			fileType := http.DetectContentType(imageBytes)

			imageAlias := models.SaveUpload(db, true, uploadResponse.SubjectId, uploadResponse.CollectionId, imageBytes, fileType)
			uploadResponse.Aliases = append(uploadResponse.Aliases, imageAlias)

			os.Remove(fileName)
		}
	}

	jsonResponse, err := json.Marshal(uploadResponse)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func discardUploadImageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var discardRequest models.ImageDiscardRequest
	err := decoder.Decode(&discardRequest)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	models.DiscardImage(db, discardRequest.SubjectId, discardRequest.CollectionId, discardRequest.ImageAlias)

	var jsonArray []string

	jsonArray = append(jsonArray, "swag")

	jsonResponse, err := json.Marshal(jsonArray)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func uploadPreviewHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var imageBytes string
	var imageType string
	err := db.QueryRow("SELECT image, image_type FROM uploaded_images WHERE alias=$1", ps.ByName("alias")).Scan(&imageBytes, &imageType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buffer := bytes.NewBufferString(imageBytes)
	w.Header().Set("Content-Type", imageType)
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		fmt.Println("unable to write image.")
	}
}

func saveImageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)

	var saveRequest models.SaveImageRequest
	err := decoder.Decode(&saveRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	alias := models.SaveImage(db, saveRequest.SubjectId, saveRequest.CollectionId, saveRequest.ImageAlias)

	var jsonArray []string

	jsonArray = append(jsonArray, alias)

	jsonResponse, err := json.Marshal(jsonArray)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func getUserImagesHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)

	var passImageRequest models.UserPassImageRequest
	err := decoder.Decode(&passImageRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	passImages := models.GetPassImages(db, passImageRequest.SubjectId, passImageRequest.CollectionId)

	jsonResponse, _ := json.Marshal(passImages)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func getImageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var imageBytes string
	var imageType string
	err := db.QueryRow("SELECT image, image_type FROM saved_images WHERE alias=$1", ps.ByName("alias")).Scan(&imageBytes, &imageType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buffer := bytes.NewBufferString(imageBytes)
	w.Header().Set("Content-Type", imageType)
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		fmt.Println("unable to write image.")
	}
}

func getRandomImageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var imageBytes string
	var imageType string
	err := db.QueryRow("SELECT image, image_type FROM random_stage_images WHERE alias=$1", ps.ByName("alias")).Scan(&imageBytes, &imageType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buffer := bytes.NewBufferString(imageBytes)
	w.Header().Set("Content-Type", imageType)
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		fmt.Println("unable to write image.")
	}
}

func newSubjectHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)

	var newSubjectRequest models.NewSubjectRequest
	err := decoder.Decode(&newSubjectRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	validationResults, birthday := newSubjectRequest.Validate()

	var jsonResponse []byte
	var jsonMap map[string]interface{}

	if len(validationResults) == 0 {
		newSubjectId := models.SaveNewSubject(db, newSubjectRequest, birthday)
		jsonMap = map[string]interface{}{"id": newSubjectId}
	} else {
		jsonMap = map[string]interface{}{"errors": validationResults}
	}

	jsonResponse, _ = json.Marshal(jsonMap)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func testHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

func errorHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

func getSubjectListHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	subjectList, err := models.GetSubjectList(db)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(subjectList)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func getCollectionListHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	collectionList := models.GetCollectionList(db)

	jsonResponse, err := json.Marshal(collectionList)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func getSubjectHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	subjectId, err := strconv.Atoi(ps.ByName("id"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	subjectProfile := models.GetSubjectProfileById(db, subjectId)

	jsonResponse, err := json.Marshal(subjectProfile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func newCollectionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)

	var newCollectionRequest models.NewCollectionRequest
	err := decoder.Decode(&newCollectionRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	validationResults := newCollectionRequest.Validate()

	var jsonResponse []byte
	var jsonMap map[string]interface{}

	if len(validationResults) == 0 {
		newCollectionId := newCollectionRequest.Save(db)
		jsonMap = map[string]interface{}{"id": newCollectionId}
	} else {
		jsonMap = map[string]interface{}{"errors": validationResults}
	}

	jsonResponse, _ = json.Marshal(jsonMap)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func saveSubjectPasswordHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	subjectId, err := strconv.Atoi(r.FormValue("subjectId"))
	if err != nil {
		fmt.Println(err)
	}

	password := r.FormValue("password")

	score, err := strconv.Atoi(r.FormValue("strength"))
	if err != nil {
		fmt.Println(err)
	}

	entropy, err := strconv.ParseFloat(r.FormValue("entropy"), 64)
	if err != nil {
		fmt.Println(err)
	}

	//if passes password validation
	models.SaveSubjectPassword(db, subjectId, password, score, entropy)

	jsonResponse, err := json.Marshal("hi")
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func saveSubjectPinHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	subjectId, err := strconv.Atoi(r.FormValue("subjectId"))
	if err != nil {
		fmt.Println(err)
	}

	pinNumber := r.FormValue("pinNumber")

	//if passes pinNumber validation
	models.SaveSubjectPin(db, subjectId, pinNumber)

	jsonResponse, err := json.Marshal("hi")
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func getNewConfigId() int {
	var configId int
	db.QueryRow("INSERT INTO test_configs(creation_date) VALUES($1) returning id;", time.Now()).Scan(&configId)
	return configId
}

func randomStageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	stages, _ := strconv.Atoi(r.FormValue("stages"))
	rows, _ := strconv.Atoi(r.FormValue("rows"))
	columns, _ := strconv.Atoi(r.FormValue("columns"))

	randomImages := *services.GetRandomImgurImages(rows * columns * stages)
	configID := getNewConfigId()

	transaction, err := db.Begin()
	if err != nil {
		log.Println("Error starting db transaction", err)
	}

	var wg sync.WaitGroup
	wg.Add(rows * columns * stages)
	i := 1
	for i <= stages {
		row := 1
		column := 1
		for _, image := range randomImages[((i - 1) * (rows * columns)):(i * rows * columns)] {
			go image.Save(&wg, transaction, configID, i, row, column)
			if column == columns {
				column = 1
				row++
			} else {
				column++
			}
		}
		i++
	}
	wg.Wait()
	err = transaction.Commit()

	if err != nil {
		log.Println("Error committing db transaction", err)
	}

	randomImagesByConfig := models.GetRandomImagesByConfigId(db, configID)
	matrixMap := models.GetMatrixMap(db, randomImagesByConfig)

	matrixResponse := new(models.ImageMatrixResponse)
	matrixResponse.Id = configID
	matrixResponse.Matrix = *matrixMap

	jsonResponse, err := json.Marshal(matrixResponse)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func randomImageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	configId, _ := strconv.Atoi(r.FormValue("configId"))
	replacingAlias := r.FormValue("alias")
	randomImage := *services.GetRandomImgurImage()
	newAlias := randomImage.Replace(db, configId, replacingAlias)

	jsonResponse, err := json.Marshal(newAlias)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func replaceStageImageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	configId, _ := strconv.Atoi(r.FormValue("configId"))
	collectionId, _ := strconv.Atoi(r.FormValue("collectionId"))
	selectedAlias := r.FormValue("selectedAlias")
	selectedStageString := r.FormValue("selectedStage")
	selectedRowString := r.FormValue("selectedRow")
	selectedColumnString := r.FormValue("selectedColumn")
	replacementAlias := r.FormValue("replacementAlias")
	replacementType := r.FormValue("replacementType")

	selectedStage, _ := strconv.Atoi(selectedStageString)
	selectedRow, _ := strconv.Atoi(selectedRowString)
	selectedColumn, _ := strconv.Atoi(selectedColumnString)

	replacement := models.ReplaceRandomImage(db, configId, collectionId, selectedAlias, selectedStage, selectedRow, selectedColumn, replacementAlias, replacementType)

	jsonResponse, err := json.Marshal(replacement)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func saveConfigHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var configSaveRequest models.Config
	err := decoder.Decode(&configSaveRequest)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = configSaveRequest.Save(db)
	var jsonArray []string

	if err != nil {
		jsonArray = append(jsonArray, err.Error())
	} else {
		jsonArray = append(jsonArray, "swag")
	}

	jsonResponse, err := json.Marshal(jsonArray)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func getConfigListHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	configList, err := models.GetConfigList(db)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse, err := json.Marshal(configList)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func postConfigHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	configId, _ := strconv.Atoi(r.FormValue("configId"))

	config := models.GetConfigById(db, configId)

	jsonResponse, err := json.Marshal(config)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func getConfigHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	configID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	configInfo, err := models.GetTestConfigInfoByConfigId(db, configID)

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(configInfo)

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func configImageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var imageBytes string
	var imageType string
	err := db.QueryRow("SELECT image, image_type FROM test_config_stage_images WHERE alias=$1", ps.ByName("alias")).Scan(&imageBytes, &imageType)
	if err != nil {
		// We need this for when user images are assigned during the trial setup in order to preview the matrix
		err = db.QueryRow("SELECT image, image_type FROM saved_images WHERE alias=$1", ps.ByName("alias")).Scan(&imageBytes, &imageType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	buffer := bytes.NewBufferString(imageBytes)
	w.Header().Set("Content-Type", imageType)
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		fmt.Println("unable to write image.")
	}
}

func testImageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var imageBytes string
	var imageType string
	err := db.QueryRow("SELECT image, image_type FROM image_trial_images WHERE alias=$1", ps.ByName("alias")).Scan(&imageBytes, &imageType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buffer := bytes.NewBufferString(imageBytes)
	w.Header().Set("Content-Type", imageType)
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		fmt.Println("unable to write image.")
	}
}

func getSubjectPassImages(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	subjectID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	subjectImages, err := models.GetSubjectPassImages(db, subjectID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(subjectImages)

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func testSettingImageSubmitHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var trialParams models.ImageTrialRequest
	err := decoder.Decode(&trialParams)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	trialID, err := trialParams.Save(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	trialInfo := models.GetImageTrialInfoById(db, trialID)

	jsonResponse, err := json.Marshal(trialInfo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func testSettingPasswordSubmitHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var trialParams models.PasswordTrialRequest
	err := decoder.Decode(&trialParams)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	trialID, err := trialParams.Save(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	trialInfo := models.GetPasswordTrialInfoById(db, trialID)

	jsonResponse, err := json.Marshal(trialInfo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func trialListHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	trialList, err := models.GetTrialList(db)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(trialList)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func trialStartHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	trialStartRequest := struct {
		TrialID      int  `json:"trialId"`
		IsImageTrial bool `json:"isImageTrial"`
	}{}
	err := decoder.Decode(&trialStartRequest)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if trialStartRequest.IsImageTrial {
		startImageTrialHandler(w, trialStartRequest.TrialID)
	} else {
		startPasswordTrialHandler(w, trialStartRequest.TrialID)
	}
}

func trialStartTimeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	request := struct {
		TrialID       int    `json:"trialId"`
		IsImageTrial  bool   `json:"isImageTrial"`
		UnixTimestamp string `json:"unixTimestamp"`
	}{}
	err := decoder.Decode(&request)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	timeStamp, err := utils.MsToTime(request.UnixTimestamp)
	if err != nil {
		timeStamp = time.Now()
		err = nil
	}

	if request.IsImageTrial {
		_, err = db.Query(`
			UPDATE image_trial_stage_results
			SET start_time = $2
			WHERE trial_id = $1 AND stage_number = 1;
		`, request.TrialID, timeStamp)
	} else {
		_, err = db.Query(`
			UPDATE password_trials
			SET start_time = $2
			WHERE id = $1;
		`, request.TrialID, timeStamp)
	}

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(struct {
		Success bool `json:"success"`
	}{true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func startImageTrialHandler(w http.ResponseWriter, trialId int) {
	trialData, err := models.GetImageTrial(db, trialId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(*trialData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func startPasswordTrialHandler(w http.ResponseWriter, trialID int) {
	trialData, err := models.GetPasswordTrial(db, trialID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(*trialData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func trialImageSubmitHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var trialSubmission models.ImageTrialSubmission
	err := decoder.Decode(&trialSubmission)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := trialSubmission.Save(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func trialPasswordSubmitHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var trialSubmission models.PasswordTrialSubmission
	err := decoder.Decode(&trialSubmission)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := trialSubmission.Save(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func subjectTrialsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	subjectIDString := ps.ByName("id")
	subjectID, err := strconv.Atoi(subjectIDString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	subjectTrialResults, err := models.GetSubjectTrialResults(db, subjectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(subjectTrialResults)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func trialDetailHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	request := struct {
		TrialID      int  `json:"trialId"`
		IsImageTrial bool `json:"isImageTrial"`
	}{}
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var jsonResponse []byte

	if request.IsImageTrial {
		trialDetails, err := models.GetImageTrialDetails(db, request.TrialID)
		if err == nil {
			jsonResponse, err = json.Marshal(trialDetails)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		trialDetails, err := models.GetPasswordTrialDetails(db, request.TrialID)
		if err == nil {
			jsonResponse, err = json.Marshal(trialDetails)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func exportHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	timeString := time.Now().Format("01_02_2006_150405")
	pathToZip, err := jobs.ExportData(db, timeString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	zippedBytes, err := ioutil.ReadFile(pathToZip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Length", strconv.Itoa(len(zippedBytes)))
	w.Header().Set("Content-Disposition", "attachment; filename=\""+timeString+".zip"+"\"")
	w.Write(zippedBytes)
}

func indexHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var index = pongo2.Must(pongo2.FromFile("./public/dist/index.html"))
	err := index.ExecuteWriter(getContextMap(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getContextMap() pongo2.Context {
	contextMap := make(pongo2.Context)
	return contextMap
}

type notFound struct{}

func (n notFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var index = pongo2.Must(pongo2.FromFile("./public/dist/index.html"))
	err := index.ExecuteWriter(getContextMap(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	defer db.Close()

	router := httprouter.New()
	router.GET("/", indexHandler)

	router.POST("/subject/new", newSubjectHandler)
	router.GET("/subject/profile/:id", getSubjectHandler)
	router.GET("/subject/images/:id", getSubjectPassImages)
	router.GET("/subject/list", getSubjectListHandler)
	router.GET("/subject/trials/:id", subjectTrialsHandler)
	router.POST("/subject/save/password", saveSubjectPasswordHandler)
	router.POST("/subject/save/pin", saveSubjectPinHandler)

	router.GET("/collections/list", getCollectionListHandler)
	router.POST("/collections/new", newCollectionHandler)

	router.POST("/upload", uploadImageHandler)
	router.POST("/upload/discard", discardUploadImageHandler)
	router.GET("/upload/preview/:alias", uploadPreviewHandler)

	router.POST("/save/image", saveImageHandler)

	router.GET("/image/:alias", getImageHandler)
	router.POST("/image/replace", replaceStageImageHandler)
	router.POST("/images", getUserImagesHandler)

	router.POST("/random/stages", randomStageHandler)
	router.POST("/random/image", randomImageHandler)
	router.GET("/random/image/:alias", getRandomImageHandler)

	router.GET("/config/:id", getConfigHandler)
	router.POST("/config", postConfigHandler)
	router.GET("/configs/image/:alias", configImageHandler)
	router.POST("/config/save", saveConfigHandler)
	router.GET("/configs/list", getConfigListHandler)

	router.GET("/test/image/:alias", testImageHandler)
	router.POST("/test/settings/image/submit", testSettingImageSubmitHandler)
	router.POST("/test/settings/password/submit", testSettingPasswordSubmitHandler)
	router.GET("/trial/list", trialListHandler)
	router.POST("/trial/start", trialStartHandler)
	router.POST("/trial/details", trialDetailHandler)
	router.POST("/trial/set-start-time", trialStartTimeHandler)
	router.POST("/trial/submit-image", trialImageSubmitHandler)
	router.POST("/trial/submit-password", trialPasswordSubmitHandler)
	router.GET("/export", exportHandler)

	router.RedirectTrailingSlash = true

	notFound := notFound{}
	router.NotFound = notFound
	router.MethodNotAllowed = notFound

	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
	)

	n.Use(negroni.NewStatic(http.Dir("./public")))

	// Add CORS support (Cross Origin Resource Sharing)
	handler := cors.New(cors.Options{
		AllowedOrigins: []string{config.CrossOriginAddress},
		AllowedHeaders: []string{"Origin", "Accept", "Content-Type", "Authorization", "Access-Control-Allow-Origin"},
	}).Handler(router)

	n.UseHandler(handler)

	n.Run(":" + config.ServerPort)
}

func setupDatabase() *sql.DB {
	db_url := fmt.Sprintf("user=%s password=%s dbname=%s  sslmode=%s",
		config.DBUser, config.DBPassword, config.DBName, config.DBSSLMode)

	db, err := sql.Open("postgres", db_url)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return db
}

type Config struct {
	CrossOriginAddress string
	ServerPort         string
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
}

func ReadConfig(configfile string) Config {
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	var config Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}
	return config
}
