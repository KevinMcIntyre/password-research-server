package jobs

import (
	"database/sql"
	"os"
	"time"
)

func ExportData(db *sql.DB) error {
	dirName := "export/" + time.Now().Format("01022006150405")
	os.Mkdir(dirName, os.ModePerm)

	subjects, err := getAllSubjectData(db)
	if err != nil {
		return err
	}
	err = createCSV(dirName, "subjects.csv", &subjects)
	if err != nil {
		return err
	}

	passwordTrials, err := getAllPasswordTrialData(db, false)
	if err != nil {
		return err
	}
	err = createCSV(dirName, "password_trials.csv", &passwordTrials)
	if err != nil {
		return err
	}

	passwordSubmissions, err := getAllPasswordSubmissions(db)
	if err != nil {
		return err
	}
	err = createCSV(dirName, "password_submissions.csv", &passwordSubmissions)
	if err != nil {
		return err
	}

	pinTrials, err := getAllPasswordTrialData(db, true)
	if err != nil {
		return err
	}
	err = createCSV(dirName, "pin_trials.csv", &pinTrials)
	if err != nil {
		return err
	}

	pinSubmissions, err := getAllPinSubmissions(db)
	if err != nil {
		return err
	}
	err = createCSV(dirName, "pin_submissions.csv", &pinSubmissions)
	if err != nil {
		return err
	}

	imageTrialConfigs, err := getAllImageConfigData(db)
	if err != nil {
		return err
	}
	err = createCSV(dirName, "image_trial_configs.csv", &imageTrialConfigs)
	if err != nil {
		return err
	}

	imageTrials, err := getAllImageTrialData(db)
	if err != nil {
		return err
	}
	err = createCSV(dirName, "image_trials.csv", &imageTrials)
	if err != nil {
		return err
	}

	imageTrialStages, err := getAllImageTrialStageData(db)
	if err != nil {
		return err
	}
	err = createCSV(dirName, "image_trials_stages.csv", &imageTrialStages)
	if err != nil {
		return err
	}

	subjectImages, err := getSubjectImages(db)
	if err != nil {
		return err
	}
	go saveSubjectImages(dirName, subjectImages)

	configImages, err := getConfigImages(db)
	if err != nil {
		return err
	}
	go saveConfigImages(dirName, configImages)

	trialUserImages, err := getTrialUserImages(db)
	if err != nil {
		return err
	}
	go saveTrialUserImages(dirName, trialUserImages)

	return nil
}
