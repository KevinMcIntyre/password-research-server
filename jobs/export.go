package jobs

import (
	"archive/zip"
	"database/sql"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ExportData(db *sql.DB, timeString string) (string, error) {
	dirName := "export/" + timeString
	os.Mkdir(dirName, os.ModePerm)

	subjects, err := getAllSubjectData(db)
	if err != nil {
		return "", err
	}
	err = createCSV(dirName, "subjects.csv", &subjects)
	if err != nil {
		return "", err
	}

	passwordTrials, err := getAllPasswordTrialData(db, false)
	if err != nil {
		return "", err
	}
	err = createCSV(dirName, "password_trials.csv", &passwordTrials)
	if err != nil {
		return "", err
	}

	passwordSubmissions, err := getAllPasswordSubmissions(db)
	if err != nil {
		return "", err
	}
	err = createCSV(dirName, "password_submissions.csv", &passwordSubmissions)
	if err != nil {
		return "", err
	}

	pinTrials, err := getAllPasswordTrialData(db, true)
	if err != nil {
		return "", err
	}
	err = createCSV(dirName, "pin_trials.csv", &pinTrials)
	if err != nil {
		return "", err
	}

	pinSubmissions, err := getAllPinSubmissions(db)
	if err != nil {
		return "", err
	}
	err = createCSV(dirName, "pin_submissions.csv", &pinSubmissions)
	if err != nil {
		return "", err
	}

	imageTrialConfigs, err := getAllImageConfigData(db)
	if err != nil {
		return "", err
	}
	err = createCSV(dirName, "image_trial_configs.csv", &imageTrialConfigs)
	if err != nil {
		return "", err
	}

	imageTrials, err := getAllImageTrialData(db)
	if err != nil {
		return "", err
	}
	err = createCSV(dirName, "image_trials.csv", &imageTrials)
	if err != nil {
		return "", err
	}

	imageTrialStages, err := getAllImageTrialStageData(db)
	if err != nil {
		return "", err
	}
	err = createCSV(dirName, "image_trials_stages.csv", &imageTrialStages)
	if err != nil {
		return "", err
	}

	subjectImages, err := getSubjectImages(db)
	if err != nil {
		return "", err
	}
	go saveSubjectImages(dirName, subjectImages)

	configImages, err := getConfigImages(db)
	if err != nil {
		return "", err
	}
	go saveConfigImages(dirName, configImages)

	trialUserImages, err := getTrialUserImages(db)
	if err != nil {
		return "", err
	}
	go saveTrialUserImages(dirName, trialUserImages)

	err = zipItUp(dirName, dirName+".zip")
	if err != nil {
		return "", err
	}

	return dirName + ".zip", nil
}

func zipItUp(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}
