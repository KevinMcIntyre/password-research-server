package models

import "database/sql"

type Trial struct {
	SubjectID            int                                     `json:"subjectId"`
	ConfigID             int                                     `json:"configId"`
	Stages               int                                     `json:"stages"`
	Rows                 int                                     `json:"rows"`
	Columns              int                                     `json:"columns"`
	ImageMaybeNotPresent bool                                    `json:"imageMaybeNotPresent"`
	Matrix               map[string]map[string]map[string]string `json:"matrix"`
	UserPassImages       []UserPassImage                         `json:"passImages"`
}

type UserPassImage struct {
	StageNumber  int    `json:"stage"`
	RowNumber    int    `json:"row"`
	ColumnNumber int    `json:"column"`
	ImageAlias   string `json:"image"`
}

func (request Trial) Save(db *sql.DB) (int, error) {
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
			SET alias = $1, is_user_image = true
			WHERE stage_number = $2, row_number = $3, column_number = $4;
		`, userPassImage.ImageAlias,
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
