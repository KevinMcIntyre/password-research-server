package models

import (
	"database/sql"
	"time"
)

type ImageTrial struct {
	ID                   int                                     `json:"id"`
	SubjectName          string                                  `json:"subjectName"`
	Stages               int                                     `json:"stages"`
	Rows                 int                                     `json:"rows"`
	Columns              int                                     `json:"columns"`
	ImageMayNotBePresent bool                                    `json:"imageMayNotBePresent"`
	Matrix               map[string]map[string]map[string]string `json:"matrix"`
}

type TrialInfo struct {
	ID           int       `json:"id"`
	SubjectName  string    `json:"subjectName"`
	TrialType    string    `json:"trialType"`
	ConfigName   string    `json:"configName"`
	CreationDate time.Time `json:"creationDate"`
}

type TrialRequest struct {
	SubjectID      int             `json:"subjectId"`
	ConfigID       int             `json:"configId"`
	Stages         int             `json:"stages"`
	UserPassImages []UserPassImage `json:"userPassImages"`
}

type UserPassImage struct {
	StageNumber  int    `json:"stage"`
	RowNumber    int    `json:"row"`
	ColumnNumber int    `json:"column"`
	ImageAlias   string `json:"alias"`
}

func (request TrialRequest) Save(db *sql.DB) (int, error) {
	var trialID int

	rows, err := db.Query(`SELECT create_image_trial($1, $2, $3);`,
		request.SubjectID,
		request.ConfigID,
		request.Stages)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&trialID); err != nil {
			return 0, err
		}
	}

	transaction, err := db.Begin()
	if err != nil {
		return 0, err
	}
	for _, userPassImage := range request.UserPassImages {
		_, err = transaction.Exec(`
			UPDATE image_trial_images
			SET 
			image = subject_image.image, 
			image_type = subject_image.image_type, 
			alias = subject_image.alias, 
			is_user_image = true
			FROM (
				SELECT
				image,
				image_type,
				replace(md5(random() :: TEXT || clock_timestamp() :: TEXT), '-' :: TEXT, '' :: TEXT) :: VARCHAR(60) AS alias
				FROM saved_images WHERE alias = $1 AND subject_id = $2
			) as subject_image
			WHERE stage_number = $3 AND row_number = $4 AND column_number = $5;
		`, userPassImage.ImageAlias,
			request.SubjectID,
			userPassImage.StageNumber,
			userPassImage.RowNumber,
			userPassImage.ColumnNumber)
		if err != nil {
			return 0, err
		}
	}

	err = transaction.Commit()

	if err != nil {
		return 0, err
	}
	return trialID, nil
}

func GetTrialList(db *sql.DB) ([]TrialInfo, error) {
	rows, err := db.Query(`
		SELECT 
		it.id,
		s.first_name || ' ' || s.last_name AS subject_name,
		'Pass-Image' as trial_type,
		tc.name as config_name,
		it.creation_date
		FROM image_trials it
		JOIN subjects s ON s.id = it.subject_id 
		JOIN test_configs tc ON tc.id = it.test_config_id
		WHERE it.id NOT IN (
			SELECT DISTINCT trial_id
			FROM image_trial_stage_results
			WHERE start_time IS NOT NULL
		)
		ORDER BY it.creation_date ASC
	`)
	if err != nil {
		return nil, err
	}

	var trials []TrialInfo
	defer rows.Close()
	for rows.Next() {
		trialInfo := new(TrialInfo)
		if err := rows.Scan(&trialInfo.ID,
			&trialInfo.SubjectName,
			&trialInfo.TrialType,
			&trialInfo.ConfigName,
			&trialInfo.CreationDate); err != nil {
			return nil, err
		}
		trials = append(trials, *trialInfo)
	}

	return trials, nil
}

func GetTrialInfoById(db *sql.DB, trialID int) TrialInfo {
	trialInfo := new(TrialInfo)
	db.QueryRow(`
		SELECT 
		it.id,
		s.first_name || ' ' || s.last_name AS subject_name,
		'Pass-Image' as trial_type,
		tc.name as config_name,
		it.creation_date
		FROM image_trials it
		JOIN subjects s ON s.id = it.subject_id 
		JOIN test_configs tc ON tc.id = it.test_config_id
		WHERE it.id = $1
	`, trialID).Scan(&trialInfo.ID,
		&trialInfo.SubjectName,
		&trialInfo.TrialType,
		&trialInfo.ConfigName,
		&trialInfo.CreationDate)

	return *trialInfo
}

func GetImageTrial(db *sql.DB, trialId int) (*ImageTrial, error) {
	var imageTrial ImageTrial
	db.QueryRow(`
		SELECT
		it.id,
		s.first_name AS subject_name,
		tc.stage_count,
		tc.rows_in_matrix,
		tc.cols_in_matrix,
		tc.image_may_not_be_present
		FROM image_trials it
		JOIN subjects s ON s.id = it.subject_id
		JOIN test_configs tc ON tc.id = it.test_config_id
		WHERE it.id = $1
	`, trialId).Scan(&imageTrial.ID,
		&imageTrial.SubjectName,
		&imageTrial.Stages,
		&imageTrial.Rows,
		&imageTrial.Columns,
		&imageTrial.ImageMayNotBePresent)
	trialImages, err := getTrialImages(db, trialId)
	if err != nil {
		return nil, err
	}

	imageMatrix := GetMatrixMap(db, trialImages)
	imageTrial.Matrix = *imageMatrix

	return &imageTrial, nil
}

func getTrialImages(db *sql.DB, trialID int) (*[]*MatrixImage, error) {
	var images []*MatrixImage
	rows, err := db.Query(`
		SELECT image.alias, image.stage_number, image.row_number, image.column_number
		FROM image_trial_images image
		WHERE image.trial_id = $1
		ORDER BY image.stage_number ASC, image.row_number, image.column_number ASC
	`, trialID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		image := new(MatrixImage)
		if err := rows.Scan(&image.Alias, &image.Stage, &image.Row, &image.Column); err != nil {
			return nil, err
		}
		images = append(images, image)
	}
	return &images, nil
}
