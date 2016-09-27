package jobs

import (
	"bytes"
	"database/sql"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
	"time"

	"github.com/KevinMcIntyre/password-research-server/utils"
)

type Image interface {
	Save() error
}

type TrialUserImage struct {
	TrialID   string
	TrialDate time.Time
	Stage     string
	Row       string
	Column    string
	Image     []byte
	ImageType string
}

type SubjectImage struct {
	ImageID     string
	SubjectID   string
	SubjectName string
	Image       []byte
	ImageType   string
}

type ConfigImage struct {
	ConfigID   string
	ConfigName string
	Stage      string
	Row        string
	Column     string
	Image      []byte
	ImageType  string
}

func (img SubjectImage) Save(path string) error {
	imgDecode, _, err := image.Decode(bytes.NewReader(img.Image))
	if err != nil {
		return err
	}

	fileExtension := strings.TrimPrefix(img.ImageType, "image/")
	out, err := os.Create(path + "/" + img.ImageID + "." + fileExtension)
	defer out.Close()
	if err != nil {
		return err
	}

	if fileExtension == "png" {
		err = png.Encode(out, imgDecode)
	} else {
		var opt jpeg.Options
		opt.Quality = 100
		err = jpeg.Encode(out, imgDecode, &opt)
	}

	if err != nil {
		return err
	}

	return nil
}

func (img TrialUserImage) Save(path string) error {
	imgDecode, _, err := image.Decode(bytes.NewReader(img.Image))
	if err != nil {
		return err
	}

	fileExtension := strings.TrimPrefix(img.ImageType, "image/")
	out, err := os.Create(path + "/" + fmt.Sprintf("stage-%s_row-%s_column-%s", img.Stage, img.Row, img.Column) + "." + fileExtension)
	defer out.Close()
	if err != nil {
		return err
	}

	if fileExtension == "png" {
		err = png.Encode(out, imgDecode)
	} else {
		var opt jpeg.Options
		opt.Quality = 100
		err = jpeg.Encode(out, imgDecode, &opt)
	}

	if err != nil {
		return err
	}

	return nil
}

func (img ConfigImage) Save(path string) error {
	imgDecode, _, err := image.Decode(bytes.NewReader(img.Image))
	if err != nil {
		return err
	}

	fileExtension := strings.TrimPrefix(img.ImageType, "image/")
	out, err := os.Create(path + "/" + fmt.Sprintf("stage-%s_row-%s_column-%s", img.Stage, img.Row, img.Column) + "." + fileExtension)
	defer out.Close()
	if err != nil {
		return err
	}

	if fileExtension == "png" {
		err = png.Encode(out, imgDecode)
	} else {
		var opt jpeg.Options
		opt.Quality = 100
		err = jpeg.Encode(out, imgDecode, &opt)
	}

	if err != nil {
		return err
	}

	return nil
}

func getSubjectImages(db *sql.DB) ([]SubjectImage, error) {
	rows, err := db.Query(`
	SELECT
	images.id AS image_id,
	subjects.id AS subject_id,
	concat(subjects.first_name, ' ', subjects.last_name) AS subject_name,
	images.image,
	images.image_type
	FROM saved_images images
	JOIN subjects ON subjects.id = images.subject_id
	WHERE subject_id != 0
	ORDER BY subject_id ASC, image_id ASC
	`)
	if err != nil {
		return nil, err
	}

	var images []SubjectImage
	for rows.Next() {
		var image SubjectImage
		err := rows.Scan(
			&image.ImageID,
			&image.SubjectID,
			&image.SubjectName,
			&image.Image,
			&image.ImageType,
		)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	return images, nil
}

func getConfigImages(db *sql.DB) ([]ConfigImage, error) {
	rows, err := db.Query(`
	SELECT
	configs.id AS config_id,
	configs.name AS config_name,
	stages.stage_number,
	images.row_number,
	images.column_number,
	images.image,
	images.image_type
	FROM test_config_stage_images images
	JOIN test_config_stages stages ON stages.id = images.stage_id
	JOIN test_configs configs ON  configs.id = stages.test_config_id
	`)
	if err != nil {
		return nil, err
	}

	var images []ConfigImage
	for rows.Next() {
		var image ConfigImage
		err := rows.Scan(
			&image.ConfigID,
			&image.ConfigName,
			&image.Stage,
			&image.Row,
			&image.Column,
			&image.Image,
			&image.ImageType,
		)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	return images, nil
}

func getTrialUserImages(db *sql.DB) ([]TrialUserImage, error) {
	rows, err := db.Query(`
	SELECT
	trials.id,
	trials.creation_date,
	images.stage_number,
	images.row_number,
	images.column_number,
	images.image,
	images.image_type
	FROM image_trial_images images
	JOIN image_trials trials ON trials.id = images.trial_id
	WHERE images.is_user_image = true
	`)
	if err != nil {
		return nil, err
	}

	var images []TrialUserImage
	for rows.Next() {
		var image TrialUserImage
		err := rows.Scan(
			&image.TrialID,
			&image.TrialDate,
			&image.Stage,
			&image.Row,
			&image.Column,
			&image.Image,
			&image.ImageType,
		)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	return images, nil
}

func saveSubjectImages(dirName string, subjectImages []SubjectImage) error {
	os.Mkdir(dirName+"/subject_images", os.ModePerm)
	for _, img := range subjectImages {
		imgDirName := dirName + "/subject_images/" + img.SubjectID + "_" + img.SubjectName
		directoryExists, err := utils.FileExists(imgDirName)
		if err != nil {
			return err
		}
		if !directoryExists {
			os.Mkdir(imgDirName, os.ModePerm)
		}
		go img.Save(imgDirName)
	}
	return nil
}

func saveConfigImages(dirName string, configImages []ConfigImage) error {
	os.Mkdir(dirName+"/config_images", os.ModePerm)
	for _, img := range configImages {
		imgDirName := dirName + "/config_images/" + img.ConfigID + "_" + img.ConfigName
		directoryExists, err := utils.FileExists(imgDirName)
		if err != nil {
			return err
		}
		if !directoryExists {
			os.Mkdir(imgDirName, os.ModePerm)
		}
		go img.Save(imgDirName)
	}
	return nil
}

func saveTrialUserImages(dirName string, configImages []TrialUserImage) error {
	os.Mkdir(dirName+"/trial_pass_images", os.ModePerm)
	for _, img := range configImages {
		imgDirName := dirName + "/trial_pass_images/" + img.TrialID + "_" + img.TrialDate.Format("01-02-2006_15-04-05")
		directoryExists, err := utils.FileExists(imgDirName)
		if err != nil {
			return err
		}
		if !directoryExists {
			os.Mkdir(imgDirName, os.ModePerm)
		}
		go img.Save(imgDirName)
	}
	return nil
}
