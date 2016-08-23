package models

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type Config struct {
	SubjectID            int                                     `json:"subjectId"`
	ConfigID             int                                     `json:"configId"`
	Name                 string                                  `json:"name"`
	Stages               int                                     `json:"stages"`
	Rows                 int                                     `json:"rows"`
	Columns              int                                     `json:"columns"`
	ImageMaybeNotPresent bool                                    `json:"imageMaybeNotPresent"`
	Matrix               map[string]map[string]map[string]string `json:"matrix"`
}

type TestConfigInfo struct {
	ConfigID             int         `json:"configId"`
	Name                 string      `json:"name"`
	Stages               int         `json:"stages"`
	Rows                 int         `json:"rows"`
	Columns              int         `json:"columns"`
	ImageMaybeNotPresent bool        `json:"imageMaybeNotPresent"`
	UserImages           []ImageInfo `json:"userImages"`
}

type ImageInfo struct {
	StageNumber  int `json:"stage"`
	RowNumber    int `json:"row"`
	ColumnNumber int `json:"column"`
}

type ConfigLabelAndId struct {
	Id    int
	Label string
}

func GetTestConfigInfoByConfigId(db *sql.DB, configId int) (TestConfigInfo, error) {
	var configInfo TestConfigInfo
	var err error
	db.QueryRow(`
		SELECT
			id,
			name,
			rows_in_matrix,
			cols_in_matrix,
			stage_count,
			image_may_not_be_present
		FROM test_configs
		WHERE id = $1
		`, configId).Scan(
		&configInfo.ConfigID,
		&configInfo.Name,
		&configInfo.Rows,
		&configInfo.Columns,
		&configInfo.Stages,
		&configInfo.ImageMaybeNotPresent)
	rows, err := db.Query(`
		SELECT
		stages.stage_number,
		images.row_number,
		images.column_number
		FROM test_config_stage_images images
		JOIN test_config_stages stages ON stages.id = images.stage_id
		JOIN test_configs configs ON configs.id = stages.test_config_id
		WHERE configs.id = $1 AND images.alias = 'user-img'
	`, configId)

	defer rows.Close()

	var userImages []ImageInfo

	for rows.Next() {
		var userImage ImageInfo
		err = rows.Scan(
			&userImage.StageNumber,
			&userImage.RowNumber,
			&userImage.ColumnNumber)
		userImages = append(userImages, userImage)
	}

	configInfo.UserImages = userImages

	return configInfo, err
}

func (request Config) SaveAsConfig(db *sql.DB) error {
	_, err := db.Exec(`
    UPDATE test_configs
    SET name = $1,
        rows_in_matrix = $2,
        cols_in_matrix = $3,
        stage_count = $4,
        image_may_not_be_present = $5
    WHERE id = $6
  `, request.Name, request.Rows, request.Columns, request.Stages, request.ImageMaybeNotPresent, request.ConfigID)
	if err != nil {
		return err
	}

	stageIdMap := make(map[string]int)
	for stage, _ := range request.Matrix {
		stageInt, _ := strconv.Atoi(stage)
		var stageId int
		db.QueryRow(`
    INSERT INTO test_config_stages (test_config_id, stage_number, creation_date)
    VALUES ($1, $2, $3)
    RETURNING id;
  `, request.ConfigID, stageInt, time.Now()).Scan(&stageId)
		stageIdMap[stage] = stageId
	}

	transaction, err := db.Begin()
	if err != nil {
		return err
	}

	for stage, _ := range request.Matrix {
		for row, _ := range request.Matrix[stage] {
			for col, _ := range request.Matrix[stage][row] {
				_, err = transaction.Exec(`
        INSERT INTO test_config_stage_images (image, image_type, stage_id, alias, row_number, column_number, creation_date)
        SELECT
          image,
          image_type,
          $1,
          CASE WHEN alias = 'user-img'
            THEN 'user-img'
            ELSE replace(md5(random() :: TEXT || clock_timestamp() :: TEXT), '-' :: TEXT, '' :: TEXT) :: VARCHAR(60)
          END AS alias,
          $5,
          $6,
          $2
        FROM random_stage_images
        WHERE test_config_id = $3 AND alias = $4 LIMIT 1
      `, stageIdMap[stage], time.Now(), request.ConfigID, request.Matrix[stage][row][col], row, col)
				if err != nil {
					return err
				}
			}
		}
	}

	err = transaction.Commit()

	if err != nil {
		return err
	}

	return nil
}

func (request Config) SaveAsTest(db *sql.DB) error {
	_, err := db.Exec(`
    INSERT INTO image_trials (subject_id, test_config_id, creation_date)
	VALUES ($1, $2, $3);
  `, request.SubjectID, request.ConfigID, time.Now())
	if err != nil {
		return err
	}

	stageIDMap := make(map[string]int)
	for stage, _ := range request.Matrix {
		stageInt, _ := strconv.Atoi(stage)
		var stageId int
		db.QueryRow(`
    INSERT INTO test_config_stages (test_config_id, stage_number, creation_date)
    VALUES ($1, $2, $3)
    RETURNING id;
  `, request.ConfigID, stageInt, time.Now()).Scan(&stageId)
		stageIDMap[stage] = stageId
	}

	transaction, err := db.Begin()
	if err != nil {
		return err
	}

	for stage, _ := range request.Matrix {
		for row, _ := range request.Matrix[stage] {
			for col, _ := range request.Matrix[stage][row] {
				_, err = transaction.Exec(`
        INSERT INTO test_config_stage_images (image, image_type, stage_id, alias, row_number, column_number, creation_date)
        SELECT
          image,
          image_type,
          $1,
          CASE WHEN alias = 'user-img'
            THEN 'user-img'
            ELSE replace(md5(random() :: TEXT || clock_timestamp() :: TEXT), '-' :: TEXT, '' :: TEXT) :: VARCHAR(60)
          END AS alias,
          $5,
          $6,
          $2
        FROM random_stage_images
        WHERE test_config_id = $3 AND alias = $4 LIMIT 1
      `, stageIDMap[stage], time.Now(), request.ConfigID, request.Matrix[stage][row][col], row, col)
				if err != nil {
					return err
				}
			}
		}
	}

	err = transaction.Commit()

	if err != nil {
		return err
	}

	return nil
}

func GetConfigList(db *sql.DB) []ConfigLabelAndId {
	rows, err := db.Query("SELECT id, name FROM test_configs WHERE id != 0 AND name IS NOT NULL ORDER BY UPPER(name) ASC")
	if err != nil {
		fmt.Println(err)
		// handle error
	}

	var configList []ConfigLabelAndId
	defer rows.Close()
	for rows.Next() {
		configLabelAndId := new(ConfigLabelAndId)
		if err := rows.Scan(&configLabelAndId.Id, &configLabelAndId.Label); err != nil {
			fmt.Println(err)
			fmt.Println(err)
			// handle error
		}
		configList = append(configList, *configLabelAndId)
	}

	return configList
}

func GetConfigById(db *sql.DB, configId int) *Config {
	config := new(Config)

	config.ConfigID = configId

	db.QueryRow(`SELECT name, stage_count, rows_in_matrix, cols_in_matrix, image_may_not_be_present
              FROM test_configs
              WHERE id = $1`, configId).Scan(&config.Name, &config.Stages, &config.Rows, &config.Columns, &config.ImageMaybeNotPresent)

	config.Matrix = *GetMatrixMap(db, configId, false)

	return config
}
