package models

import "time"

type ImageTrial struct {
	SubjectID    int
	ConfigID     int
	StartTime    time.Time
	EndTime      time.Time
	ImageMatrix  map[string]interface{}
	Notes        string
	CreationDate time.Time
	Stages       []ImageTrialStage
}

type ImageTrialStage struct {
	SuccessfulAuthentication bool
	ConfigStageID            int
	SelectedTestImageID      int
	CorrectSavedImageID      int
	StartTime                time.Time
	EndTime                  time.Time
}
