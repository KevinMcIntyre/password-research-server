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

func (request Trial) Save(db *sql.DB) int {
	var trialID int
	db.QueryRow(`SELECT create_image_trial($1, $2, $3);`,
		request.SubjectID,
		request.ConfigID,
		request.Stages).Scan(&trialID)

	return trialID
}
