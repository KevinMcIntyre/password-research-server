package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/KevinMcIntyre/password-research-server/utils"
)

type ImgurImage struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Nsfw    bool   `json: "nsfw"`
	Link    string `json:"link"`
	IsAlbum bool   `json:"is_album"`
}

type ImgurResponse struct {
	Images  []ImgurImage `json:"data"`
	Success bool         `json:"success"`
	Status  int          `json:"status"`
}

func (imgur ImgurImage) Save(waitgroup *sync.WaitGroup, transaction *sql.Tx, testConfigId int, stage int, row int, column int) {
	defer waitgroup.Done()
	filePath := imgurImageToFile(imgur)
	utils.ResizeImage(*filePath)
	imageBytes := utils.ImageFileToByteArray(*filePath)
	_, err := transaction.Exec("INSERT INTO random_stage_images (image, image_type, test_config_id, stage_number, row_number, column_number, alias, creation_date) VALUES($1, $2, $3, $4, $5, $6, replace(md5(random()::text || clock_timestamp()::text), '-'::text, ''::text)::varchar(60), $7);",
		imageBytes, imgur.Type, testConfigId, stage, row, column, time.Now())
	if err != nil {
		log.Println("Error saving random image", err)
	}
	os.Remove(*filePath)
}

func (imgur ImgurImage) Replace(db *sql.DB, testConfigId int, replacingAlias string) string {
	fmt.Println("ConfigId", testConfigId)
	fmt.Println("Replacing alias", replacingAlias)
	filePath := imgurImageToFile(imgur)
	utils.ResizeImage(*filePath)
	imageBytes := utils.ImageFileToByteArray(*filePath)
	var alias string
	db.QueryRow("INSERT INTO random_stage_images (image, image_type, test_config_id, replacement_alias, alias, creation_date, stage_number, column_number, row_number) VALUES($1, $2, $3, $4, replace(md5(random()::text || clock_timestamp()::text), '-'::text, ''::text)::varchar(60), $5, -1, -1, -1) returning alias;",
		imageBytes, imgur.Type, testConfigId, replacingAlias, time.Now()).Scan(&alias)
	os.Remove(*filePath)
	return alias
}

func imgurImageToFile(imgur ImgurImage) *string {
	response, e := http.Get(imgur.Link)
	if e != nil {
		log.Println("Error converting imgur image []Byte", e)
		return nil
	}

	defer response.Body.Close()

	//open a file for writing
	filePath := "./imgur/tmp/" + imgur.Id + "." + strings.Split(imgur.Type, "/")[1]
	file, err := os.Create(filePath)

	if err != nil {
		log.Println(err)
		return nil
	}
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	file.Close()

	return &filePath
}

func FilterImgurImages(images []ImgurImage) []ImgurImage {
	var filteredImages []ImgurImage
	i := 1
	for _, image := range images {
		if !image.IsAlbum && !image.Nsfw && (strings.Contains(image.Type, "jpeg") || strings.Contains(image.Type, "jpg") || strings.Contains(image.Type, "png")) {
			if image.Height >= 200 && image.Width >= 200 && (image.Width/image.Height) < 2 && (image.Height/image.Width) < 2 {
				filteredImages = append(filteredImages, image)
			}
		}
		i++
	}

	return filteredImages
}

func getGallery(galleryNumber int) string {
	galleries := [...]string{
		"plants",
		"wallpaper",
		"wallpapers",
		"aww",
		"art",
		"photos",
		"AmateurPhotography",
		"earthporn",
		"AmateurEarthPorn",
		"BackgroundArt",
		"EyeCandy",
		"ITookAPicture",
		"Nature",
		"cats",
		"dogs",
		"food",
		"foodporn",
		"cars",
		"carporn",
		"woahdude",
		"motorcycles",
		"caferacers",
		"drawings",
		"oldschoolcool",
		"rateme",
		"babies",
		"rateme",
		"architecture",
	}
	return galleries[galleryNumber]
}

func GetRandomImgurImagesFromGallery(galleryNumber int, pageNumber int) *[]ImgurImage {
	client := &http.Client{}
	page := strconv.Itoa(pageNumber)
	requestUrl := "https://api.imgur.com/3/gallery/r/" + getGallery(galleryNumber) + "/top/all/" + page
	log.Println("Attempting request to: " + requestUrl)
	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		log.Println(err.Error())
		log.Println("Request URL: " + requestUrl)
	}
	request.Header.Add("Authorization", "Client-ID e0d1f81b6d72d22")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	response, err := client.Do(request)
	defer response.Body.Close()
	if err != nil {
		log.Println(err.Error())
	}
	if response.StatusCode != 200 {
		log.Println("Request URL: " + requestUrl)
		log.Println("Received a " + response.Status + " response while attempting to get imgur images")
		return nil
	} else {
		imgurResponse, err := parseRequestToImgurResponse(response.Body)
		if err != nil {
			log.Println("An error occured while checking the imgur response: " + err.Error())
			return nil
		}

		filteredImages := FilterImgurImages(imgurResponse.Images)
		return &filteredImages
	}
}

func parseRequestToImgurResponse(body io.ReadCloser) (*ImgurResponse, error) {
	defer body.Close()
	responseJson, err := ioutil.ReadAll(body)

	if err != nil {
		return nil, err
	}

	var imgurResponse ImgurResponse
	if err := json.Unmarshal(responseJson, &imgurResponse); err != nil {
		return nil, err
	}

	return &imgurResponse, nil
}

func GetRandomImgurImages(imagesRequested int) *[]ImgurImage {
	rand.Seed(int64(time.Now().Nanosecond()))
	var images []ImgurImage
	var chosenGalleries []int
	j := 0

	for len(images) < imagesRequested {
		for j < 10 {
			pageNumber := (rand.Intn(2) + 1)
			galleryNumber := (rand.Intn(27) + 1)
			for utils.IntInSlice(galleryNumber, chosenGalleries) {
				galleryNumber = (rand.Intn(27) + 1)
			}
			chosenGalleries = append(chosenGalleries, galleryNumber)

			images = append(images, *GetRandomImgurImagesFromGallery(galleryNumber, pageNumber)...)
			j++
		}
	}

	// shuffle images
	for p, _ := range images {
		v := rand.Intn(p + 1)
		images[p], images[v] = images[v], images[p]
	}

	// return only the amount of images needed
	images = images[:imagesRequested]

	return &images
}

func GetRandomImgurImage() *ImgurImage {
	rand.Seed(int64(time.Now().Nanosecond()))
	var images []ImgurImage

	for len(images) < 1 {
		pageNumber := (rand.Intn(2) + 1)
		galleryNumber := (rand.Intn(27) + 1)
		images = append(images, *GetRandomImgurImagesFromGallery(galleryNumber, pageNumber)...)
	}

	// shuffle images
	for p, _ := range images {
		v := rand.Intn(p + 1)
		images[p], images[v] = images[v], images[p]
	}

	return &images[0]
}
