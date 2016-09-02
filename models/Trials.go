package models

import "database/sql"

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

	db.QueryRow(`SELECT create_image_trial($1, $2, $3);`,
		request.SubjectID,
		request.ConfigID,
		request.Stages).Scan(&trialID)

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
