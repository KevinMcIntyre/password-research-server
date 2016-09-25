package jobs

import (
	"database/sql"
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

type Subject struct {
	ID               string `csv:"ID"`
	EMail            string `csv:"EMAIL"`
	Password         string `csv:"PASSWORD"`
	PasswordEntropy  string `csv:"PASSWORD_ENTROPY"`
	PasswordStrength string `csv:"PASSWORD_STRENGTH"`
	PinNumber        string `csv:"PIN_NUMBER"`
	FirstName        string `csv:"FIRST_NAME"`
	LastName         string `csv:"LAST_NAME"`
	Birthday         string `csv:"BIRTH_DATE"`
	CreationDate     string `csv:"CREATION_DATE"`
}

type PasswordTrial struct {
	ID              string `csv:"ID"`
	SubjectID       string `csv:"EMAIL"`
	AttemptsAllowed string `csv:"ATTEMPTS_ALLOWED"`
	PassedAuth      string `csv:"PASSED_AUTH"`
	StartTime       string `csv:"START_TIME"`
	EndTime         string `csv:"END_TIME"`
	Notes           string `csv:"NOTES"`
	CreationDate    string `csv:"CREATION_DATE"`
}

type PasswordSubmission struct {
	ID              string `csv:"ID"`
	PasswordTrialID string `csv:"PASSWORD_TRIAL_ID"`
	PasswordEntered string `csv:"PASSWORD_ENTERED"`
	AttemptNumber   string `csv:"ATTEMPT_NUMBER"`
	PassedAuth      string `csv:"PASSED_AUTH"`
	SubmissionTime  string `csv:"SUBMISSION_TIME"`
}

type PinSubmission struct {
	ID             string `csv:"ID"`
	PinTrialID     string `csv:"PIN_TRIAL_ID"`
	PinEntered     string `csv:"PIN_ENTERED"`
	AttemptNumber  string `csv:"ATTEMPT_NUMBER"`
	PassedAuth     string `csv:"PASSED_AUTH"`
	SubmissionTime string `csv:"SUBMISSION_TIME"`
}

type ImageTrialConfig struct {
	ID                   string `csv:"ID"`
	Name                 string `csv:"NAME"`
	Rows                 string `csv:"ROWS"`
	Columns              string `csv:"COLUMNS"`
	Stages               string `csv:"STAGES"`
	ImageMayNotBePresent string `csv:"IMAGE_MAY_NOT_BE_PRESENT"`
	CreationDate         string `csv:"CREATION_DATE"`
}

type ImageTrial struct {
	ID           string `csv:"ID"`
	SubjectID    string `csv:"SUBJECT_ID"`
	ConfigID     string `csv:"CONFIG_ID"`
	PassedAuth   string `csv:"PASSED_AUTH"`
	Notes        string `csv:"NOTES"`
	CreationDate string `csv:"CREATION_DATE"`
}

type ImageTrialStage struct {
	ID            string `csv:"ID"`
	TrialID       string `csv:"TRIAL_ID"`
	StageNumber   string `csv:"STAGE_NUMBER"`
	SelectedImage string `csv:"SELECTED_IMAGE"`
	CorrectImages string `csv:"CORRECT_IMAGES"`
	PassedAuth    string `csv:"PASSED_AUTH"`
	StartTime     string `csv:"START_TIME"`
	EndTime       string `csv:"END_TIME"`
}

func getAllImageTrialStageData(db *sql.DB) ([]ImageTrialStage, error) {
	rows, err := db.Query(`
	SELECT
	stages.id::VARCHAR,
	stages.trial_id::VARCHAR,
	stages.stage_number::VARCHAR,
	CASE WHEN selected_trial_image_id != 0
		THEN concat('"', concat('Row: ', images.row_number::VARCHAR, ' Column: ', images.column_number::VARCHAR), '"')
		ELSE '"None"'
	END AS selected_image,
	CASE WHEN correct_images IS NOT NULL
		THEN correct_images.correct_images
		ELSE '"None"'
	END AS correct_images,
	stages.passed_auth::VARCHAR,
	stages.start_time::VARCHAR,
	stages.end_time::VARCHAR
	FROM image_trial_stage_results stages
	JOIN image_trial_images images ON images.id = stages.selected_trial_image_id
	FULL JOIN (
		SELECT trial_id, stage_number, concat('"', string_agg(concat('Row: ', row_number::VARCHAR, ' Column: ', column_number::VARCHAR), ', '), '"') AS correct_images
		FROM image_trial_images WHERE is_user_image = true
		GROUP BY trial_id, stage_number
	) AS correct_images ON concat('trial', correct_images.trial_id, 'stage', correct_images.stage_number) = concat('trial', stages.trial_id, 'stage', stages.stage_number)
	WHERE stages.trial_id IN (
		WITH counts AS (
			SELECT
			stages.trial_id,
			count(stages.end_time) AS count
			FROM image_trial_stage_results stages
			GROUP BY stages.trial_id
		)
		SELECT
		counts.trial_id
		FROM counts
		JOIN image_trials trials ON trials.id = counts.trial_id
		JOIN test_configs configs ON configs.id = trials.test_config_id
		WHERE counts.count = configs.stage_count
	)
	ORDER BY trial_id ASC, stage_number ASC
	`)
	if err != nil {
		return nil, err
	}

	var trialStages []ImageTrialStage
	for rows.Next() {
		var stage ImageTrialStage
		err := rows.Scan(
			&stage.ID,
			&stage.TrialID,
			&stage.StageNumber,
			&stage.SelectedImage,
			&stage.CorrectImages,
			&stage.PassedAuth,
			&stage.StartTime,
			&stage.EndTime,
		)
		if err != nil {
			return nil, err
		}
		trialStages = append(trialStages, stage)
	}

	return trialStages, nil
}

func getAllImageTrialData(db *sql.DB) ([]ImageTrial, error) {
	rows, err := db.Query(`
	WITH image_trial_data AS (
		SELECT
		trials.id,
		trials.subject_id,
		trials.test_config_id,
		trials.notes,
		CASE WHEN image_stage_success_counts.successful_auths = tc.stage_count
			THEN TRUE
			ELSE FALSE
		END AS successful_auth,
		CASE WHEN completed_stages.count = tc.stage_count
			THEN TRUE
			ELSE FALSE
		END AS trial_complete,
		trials.creation_date
		FROM image_trials trials
		JOIN test_configs tc ON tc.id = trials.test_config_id
		JOIN (
			SELECT trial_id, count(passed_auth) AS successful_auths
			FROM image_trial_stage_results
			WHERE passed_auth = true
			GROUP BY trial_id
		) AS image_stage_success_counts ON image_stage_success_counts.trial_id = trials.id
		JOIN (
			SELECT trial_id, count(end_time) AS count
			FROM image_trial_stage_results
			GROUP BY trial_id
		) AS completed_stages ON completed_stages.trial_id = trials.id
	)
	SELECT
	id,
	subject_id,
	test_config_id,
	CASE WHEN notes IS NULL
		THEN ''
		ELSE notes
	END AS notes,
	successful_auth,
	creation_date
	FROM image_trial_data WHERE trial_complete = true
	`)
	if err != nil {
		return nil, err
	}

	var trials []ImageTrial
	for rows.Next() {
		var trial ImageTrial
		err := rows.Scan(
			&trial.ID,
			&trial.SubjectID,
			&trial.ConfigID,
			&trial.Notes,
			&trial.PassedAuth,
			&trial.CreationDate,
		)
		if err != nil {
			return nil, err
		}
		trials = append(trials, trial)
	}

	return trials, nil
}

func getAllImageConfigData(db *sql.DB) ([]ImageTrialConfig, error) {
	rows, err := db.Query(`SELECT * FROM test_configs WHERE name IS NOT NULL`)
	if err != nil {
		return nil, err
	}

	var configs []ImageTrialConfig
	for rows.Next() {
		var config ImageTrialConfig
		err := rows.Scan(
			&config.ID,
			&config.Name,
			&config.Rows,
			&config.Columns,
			&config.Stages,
			&config.ImageMayNotBePresent,
			&config.CreationDate,
		)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}

func getAllSubjectData(db *sql.DB) ([]Subject, error) {
	rows, err := db.Query(`
	SELECT
	CASE WHEN id IS NULL
		THEN ''
		ELSE id::varchar
	END AS id,
	CASE WHEN email IS NULL
		THEN ''
		ELSE email
	END AS email,
	CASE WHEN password IS NULL
		THEN ''
		ELSE password
	END AS password,
	CASE WHEN password_entropy IS NULL
		THEN ''
		ELSE password_entropy::varchar
	END AS password_strength,
	CASE WHEN password_strength IS NULL
		THEN ''
		ELSE password_strength::varchar
	END AS password_strength,
	CASE WHEN pin_number IS NULL
		THEN ''
		ELSE pin_number
	END AS pin_number,
	CASE WHEN first_name IS NULL
		THEN ''
		ELSE first_name
	END AS first_name,
	CASE WHEN last_name IS NULL
		THEN ''
		ELSE last_name
	END AS last_name,
	CASE WHEN birth_date IS NULL
		THEN ''
		ELSE birth_date::varchar
	END AS birth_date,
	CASE WHEN creation_date IS NULL
		THEN ''
		ELSE creation_date::varchar
	END AS creation_date
	FROM subjects
	WHERE id != 0
	`)
	if err != nil {
		return nil, err
	}

	var subjects []Subject
	for rows.Next() {
		var subject Subject
		err := rows.Scan(
			&subject.ID,
			&subject.EMail,
			&subject.Password,
			&subject.PasswordEntropy,
			&subject.PasswordStrength,
			&subject.PinNumber,
			&subject.FirstName,
			&subject.LastName,
			&subject.Birthday,
			&subject.CreationDate,
		)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}

	return subjects, nil
}

func getAllPasswordTrialData(db *sql.DB, isPinNumber bool) ([]PasswordTrial, error) {
	var whereclauseParam string
	if isPinNumber {
		whereclauseParam = `'pin'`
	} else {
		whereclauseParam = `'password'`
	}
	rows, err := db.Query(`
	SELECT
	CASE WHEN id IS NULL
		THEN ''
		ELSE id::varchar
	END AS id,
	CASE WHEN subject_id IS NULL
		THEN ''
		ELSE subject_id::varchar
	END AS subject_id,
	CASE WHEN attempts_allowed IS NULL
		THEN ''
		ELSE attempts_allowed::varchar
	END AS attempts_allowed,
	passed_auth::varchar,
	CASE WHEN start_time IS NULL
		THEN ''
		ELSE start_time::varchar
	END AS start_time,
	CASE WHEN end_time IS NULL
		THEN ''
		ELSE end_time::varchar
	END AS end_time,
	CASE WHEN notes IS NULL
		THEN ''
		ELSE notes
	END AS notes,
	CASE WHEN creation_date IS NULL
		THEN ''
		ELSE creation_date::varchar
	END AS creation_date
	FROM password_trials
	WHERE passed_auth IS NOT NULL AND trial_type =` + whereclauseParam)
	if err != nil {
		return nil, err
	}

	var trials []PasswordTrial
	for rows.Next() {
		var trial PasswordTrial
		err := rows.Scan(
			&trial.ID,
			&trial.SubjectID,
			&trial.AttemptsAllowed,
			&trial.PassedAuth,
			&trial.StartTime,
			&trial.EndTime,
			&trial.Notes,
			&trial.CreationDate,
		)
		if err != nil {
			return nil, err
		}
		trials = append(trials, trial)
	}

	return trials, nil
}

func getAllPasswordSubmissions(db *sql.DB) ([]PasswordSubmission, error) {
	rows, err := db.Query(getSubmissionQuery(false))
	if err != nil {
		return nil, err
	}

	var submissions []PasswordSubmission
	for rows.Next() {
		var submission PasswordSubmission
		err := rows.Scan(
			&submission.ID,
			&submission.PasswordTrialID,
			&submission.PasswordEntered,
			&submission.AttemptNumber,
			&submission.PassedAuth,
			&submission.SubmissionTime,
		)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, submission)
	}
	return submissions, nil
}

func getAllPinSubmissions(db *sql.DB) ([]PinSubmission, error) {
	rows, err := db.Query(getSubmissionQuery(true))
	if err != nil {
		return nil, err
	}

	var submissions []PinSubmission
	for rows.Next() {
		var submission PinSubmission
		err := rows.Scan(
			&submission.ID,
			&submission.PinTrialID,
			&submission.PinEntered,
			&submission.AttemptNumber,
			&submission.PassedAuth,
			&submission.SubmissionTime,
		)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, submission)
	}
	return submissions, nil
}

func createCSV(directory string, fileName string, structArray interface{}) error {
	file, err := os.OpenFile(directory+"/"+fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)

	if err != nil {
		return err
	}
	err = gocsv.MarshalFile(structArray, file)
	if err != nil {
		return err
	}
	return nil
}

func CreateCSVFiles(db *sql.DB) error {
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

	return nil
}

func getSubmissionQuery(isPinNumber bool) string {
	var whereclauseParam string
	if isPinNumber {
		whereclauseParam = `'pin'`
	} else {
		whereclauseParam = `'password'`
	}
	return `
		SELECT
		CASE WHEN submission.id IS NULL
			THEN ''
			ELSE submission.id::varchar
		END AS id,
		CASE WHEN submission.trial_id IS NULL
			THEN ''
			ELSE submission.trial_id::varchar
		END AS trial_id,
		CASE WHEN submission.password_entered IS NULL
			THEN ''
			ELSE submission.password_entered
		END AS password_entered,
		CASE WHEN submission.attempt_number IS NULL
			THEN ''
			ELSE submission.attempt_number::varchar
		END AS attempt_number,
		CASE WHEN trial_passwords.subject_password = submission.password_entered
			THEN true
			ELSE false
		END AS passed_auth,
		CASE WHEN submission.submission_time IS NULL
			THEN ''
			ELSE submission.submission_time::varchar
		END AS submission_time
		FROM passwords_submitted submission
		JOIN password_trials trial ON trial.id = submission.trial_id
		JOIN (
			SELECT
			trials.id AS trial_id,
			CASE WHEN trials.trial_type = 'pin'
				THEN subjects.pin_number
				ELSE subjects.password
			END AS subject_password
			FROM password_trials trials
			JOIN subjects ON subjects.id = trials.subject_id
		) AS trial_passwords ON trial_passwords.trial_id = trial.id
		WHERE trial.trial_type = ` + whereclauseParam
}
