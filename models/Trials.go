package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/KevinMcIntyre/password-research-server/utils"
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

type PasswordTrial struct {
	ID              int    `json:"id"`
	TrialType       string `json:"trialType"`
	SubjectName     string `json:"subjectName"`
	AttemptsAllowed int    `json:"attemptsAllowed"`
}

type TrialInfo struct {
	ID              int       `json:"id"`
	SubjectName     string    `json:"subjectName"`
	TrialType       string    `json:"trialType"`
	AttemptsAllowed int       `json:"attemptsAllowed"`
	ConfigName      string    `json:"configName"`
	CreationDate    time.Time `json:"creationDate"`
}

type PasswordTrialRequest struct {
	SubjectID       int    `json:"subjectId"`
	TrialType       string `json:"trialType"`
	AllowedAttempts int    `json:"allowedAttempts"`
}

type ImageTrialRequest struct {
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

type ImageTrialSubmission struct {
	TrialID       int    `json:"trialId"`
	StageNumber   int    `json:"stage"`
	ImageAlias    string `json:"imageAlias"`
	UnixTimestamp string `json:"unixTimestamp"`
}

type PasswordTrialSubmission struct {
	TrialID       int    `json:"trialId"`
	Password      string `json:"password"`
	UnixTimestamp string `json:"unixTimestamp"`
}

type TrialSubmissionResponse struct {
	TrialComplete            bool `json:"trialComplete"`
	SuccessfulAuthentication bool `json:"successfulAuth"`
}

type ImageTrialDetail struct {
	StageNumber        int      `json:"stage"`
	SelectedImageAlias string   `json:"selectedAlias"`
	CorrectImageAlias  []string `json:"correctAlias"`
	Success            bool     `json:"success"`
	TimeSpentInSeconds float64  `json:"timeSpentInSeconds"`
}

type PasswordTrialDetail struct {
	AttemptNumber      int     `json:"attemptNumber"`
	PasswordEntered    string  `json:"passwordEntered"`
	CorrectPassword    string  `json:"correctPassword"`
	Success            bool    `json:"success"`
	TimeSpentInSeconds float64 `json:"timeSpentInSeconds"`
}

func GetImageTrialDetails(db *sql.DB, trialId int) ([]ImageTrialDetail, error) {
	rows, err := db.Query(`
		SELECT
		stage_results.stage_number,
		CASE WHEN images.alias IS NULL
			THEN ''
			ELSE images.alias
		END AS selected_image,
		CASE WHEN correct_images.user_images IS NULL
			THEN '{}'
			ELSE correct_images.user_images
		END AS user_images,
		stage_results.passed_auth AS success,
		stage_results.start_time,
		stage_results.end_time
		FROM image_trial_stage_results stage_results
		JOIN image_trial_images images ON images.id = stage_results.selected_trial_image_id
		FULL JOIN (
			SELECT trial_id, stage_number, array_agg(alias) AS user_images
			FROM image_trial_images
			WHERE trial_id = $1 AND is_user_image = true
			GROUP BY trial_id, stage_number
		) AS correct_images ON concat(correct_images.trial_id, correct_images.stage_number) = concat(stage_results.trial_id, stage_results.stage_number)
		WHERE stage_results.trial_id = $1
		ORDER BY stage_number ASC
	`, trialId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var trialDetails []ImageTrialDetail
	for rows.Next() {
		trialDetail := new(ImageTrialDetail)
		startTime := new(time.Time)
		endTime := new(time.Time)
		correctImageAliasArrayString := new(string)
		if err := rows.Scan(
			&trialDetail.StageNumber,
			&trialDetail.SelectedImageAlias,
			correctImageAliasArrayString,
			&trialDetail.Success,
			&startTime,
			&endTime); err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		testDuration := endTime.Sub(*startTime)
		correctImageAliasArray, err := utils.ParsePGArray(*correctImageAliasArrayString)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		trialDetail.CorrectImageAlias = correctImageAliasArray
		trialDetail.TimeSpentInSeconds = utils.RoundUp(testDuration.Seconds(), 2)

		trialDetails = append(trialDetails, *trialDetail)
	}
	return trialDetails, nil
}

func GetPasswordTrialDetails(db *sql.DB, trialId int) ([]PasswordTrialDetail, error) {
	rows, err := db.Query(`
	SELECT
	submitted.attempt_number,
	submitted.password_entered,
	CASE WHEN trial.trial_type = 'pin'
		THEN subjects.pin_number
		ELSE subjects.password
	END AS correct_password,
	CASE WHEN trial.trial_type = 'pin'
		THEN
		CASE WHEN submitted.password_entered = subjects.pin_number
			THEN true
			ELSE false
		END
		ELSE
		CASE WHEN submitted.password_entered = subjects.password
			THEN true
			ELSE false
		END
	END AS success,
	start_times.start_time,
	submitted.submission_time
	FROM passwords_submitted submitted
	JOIN password_trials trial ON trial.id = submitted.trial_id
	JOIN subjects ON subjects.id = trial.subject_id
	JOIN (
		SELECT
			attempt_number,
			CASE WHEN attempt_number = 1
			THEN trial.start_time
			ELSE lag(submission_time) OVER submission_window
			END AS start_time
		FROM passwords_submitted submitted
		JOIN password_trials trial ON trial.id = submitted.trial_id
		WHERE submitted.trial_id = $1
		WINDOW submission_window as (partition by trial_id order by attempt_number)
	) AS start_times ON start_times.attempt_number = submitted.attempt_number
	WHERE submitted.trial_id = $1
	ORDER BY submitted.attempt_number
	`, trialId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var trialDetails []PasswordTrialDetail
	for rows.Next() {
		trialDetail := new(PasswordTrialDetail)
		startTime := new(time.Time)
		endTime := new(time.Time)
		if err := rows.Scan(
			&trialDetail.AttemptNumber,
			&trialDetail.PasswordEntered,
			&trialDetail.CorrectPassword,
			&trialDetail.Success,
			&startTime,
			&endTime); err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		testDuration := endTime.Sub(*startTime)
		trialDetail.TimeSpentInSeconds = utils.RoundUp(testDuration.Seconds(), 2)

		trialDetails = append(trialDetails, *trialDetail)
	}
	return trialDetails, nil
}

func (submission PasswordTrialSubmission) Save(db *sql.DB) (*TrialSubmissionResponse, error) {
	var response TrialSubmissionResponse
	timeStamp, err := utils.MsToTime(submission.UnixTimestamp)
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(`
	SELECT *
	FROM submit_password_submission($1, $2, $3)
	AS f(trial_complete bool, successful_auth bool)
	`, submission.TrialID, submission.Password, timeStamp)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&response.TrialComplete, &response.SuccessfulAuthentication); err != nil {
			return nil, err
		}
	}

	return &response, nil
}

func (submission ImageTrialSubmission) Save(db *sql.DB) (*TrialSubmissionResponse, error) {
	var response TrialSubmissionResponse
	timeStamp, err := utils.MsToTime(submission.UnixTimestamp)
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(`
	SELECT *
	FROM submit_image_selection($1, $2, $3, $4)
	AS f(trial_complete bool, successful_auth bool)
	`, submission.TrialID, submission.StageNumber, submission.ImageAlias, timeStamp)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&response.TrialComplete, &response.SuccessfulAuthentication); err != nil {
			return nil, err
		}
	}

	return &response, nil
}

func (request PasswordTrialRequest) Save(db *sql.DB) (int, error) {
	var trialID int

	rows, err := db.Query(`
		INSERT INTO password_trials(subject_id, trial_type, attempts_allowed, creation_date)
		VALUES($1, $2, $3, $4)
		RETURNING id;
	`,
		request.SubjectID,
		request.TrialType,
		request.AllowedAttempts,
		time.Now())
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&trialID); err != nil {
			return 0, err
		}
	}

	if err != nil {
		return 0, err
	}
	return trialID, nil
}

func GetPasswordTrialInfoById(db *sql.DB, trialID int) TrialInfo {
	trialInfo := new(TrialInfo)
	db.QueryRow(`
		SELECT 
		pt.id,
		s.first_name || ' ' || s.last_name AS subject_name,
		initcap(pt.trial_type) AS trial_type,
		pt.attempts_allowed,
		'N/A' as config_name,
		pt.creation_date
		FROM password_trials pt
		JOIN subjects s ON s.id = pt.subject_id 
		WHERE pt.id = $1
	`, trialID).Scan(&trialInfo.ID,
		&trialInfo.SubjectName,
		&trialInfo.TrialType,
		&trialInfo.AttemptsAllowed,
		&trialInfo.ConfigName,
		&trialInfo.CreationDate)

	return *trialInfo
}

func (request ImageTrialRequest) Save(db *sql.DB) (int, error) {
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
			WHERE trial_id = $6 AND stage_number = $3 AND row_number = $4 AND column_number = $5;
		`, userPassImage.ImageAlias,
			request.SubjectID,
			userPassImage.StageNumber,
			userPassImage.RowNumber,
			userPassImage.ColumnNumber,
			trialID)
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
		WITH trials AS (
			SELECT
			it.id AS id,
			s.first_name || ' ' || s.last_name AS subject_name,
			'Pass-Image' AS trial_type,
			1 AS attempts_allowed,
			tc.name As config_name,
			it.creation_date AS creation_date
			FROM image_trials it
			JOIN subjects s ON s.id = it.subject_id 
			JOIN test_configs tc ON tc.id = it.test_config_id
			WHERE it.id NOT IN (
				SELECT DISTINCT trial_id
				FROM image_trial_stage_results
				WHERE start_time IS NOT NULL
			)

			UNION

			SELECT 
			pt.id AS id,
			s.first_name || ' ' || s.last_name AS subject_name,
			initcap(pt.trial_type) AS trial_type,
			pt.attempts_allowed AS attempts_allowed,
			'N/A' AS config_name,
			pt.creation_date AS creation_date
			FROM password_trials pt
			JOIN subjects s ON s.id = pt.subject_id 
			WHERE pt.start_time IS NULL AND pt.end_time IS NULL AND pt.passed_auth IS NULL
		)
		SELECT * FROM trials ORDER BY creation_date ASC
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
			&trialInfo.AttemptsAllowed,
			&trialInfo.ConfigName,
			&trialInfo.CreationDate); err != nil {
			return nil, err
		}
		trials = append(trials, *trialInfo)
	}

	return trials, nil
}

func GetImageTrialInfoById(db *sql.DB, trialID int) TrialInfo {
	trialInfo := new(TrialInfo)
	db.QueryRow(`
		SELECT 
		it.id,
		s.first_name || ' ' || s.last_name AS subject_name,
		'Pass-Image' as trial_type,
		1 as attempts_allowed,
		tc.name as config_name,
		it.creation_date
		FROM image_trials it
		JOIN subjects s ON s.id = it.subject_id 
		JOIN test_configs tc ON tc.id = it.test_config_id
		WHERE it.id = $1
	`, trialID).Scan(&trialInfo.ID,
		&trialInfo.SubjectName,
		&trialInfo.TrialType,
		&trialInfo.AttemptsAllowed,
		&trialInfo.ConfigName,
		&trialInfo.CreationDate)

	return *trialInfo
}

func GetPasswordTrial(db *sql.DB, trialID int) (*PasswordTrial, error) {
	passwordTrial := new(PasswordTrial)

	err := db.QueryRow(`
		SELECT 
		pt.id,
		s.first_name AS subject_name,
		initcap(pt.trial_type) AS trial_type,
		pt.attempts_allowed
		FROM password_trials pt
		JOIN subjects s ON s.id = pt.subject_id 
		WHERE pt.id = $1
	`, trialID).Scan(&passwordTrial.ID,
		&passwordTrial.SubjectName,
		&passwordTrial.TrialType,
		&passwordTrial.AttemptsAllowed)

	return passwordTrial, err
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
