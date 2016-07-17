package models

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type RandomImage struct {
	Alias  string
	Stage  int
	Row    int
	Column int
}

type ImageMatrixResponse struct {
	Id     int
	Matrix map[string]map[string]map[string]string
}

func ReplaceRandomImage(db *sql.DB, configId int, collectionId int, selectedAlias string, replacementAlias string, replacementType string) *RandomImage {
	replacedImage := new(RandomImage)

	if replacementType == "random-img" {
		db.QueryRow(`UPDATE random_stage_images
    SET stage_number = previous.stagenum, row_number = previous.rownum, column_number = previous.colnum, replacement_alias = NULL
    FROM (
      SELECT stage_number AS stagenum, row_number AS rownum, column_number AS colnum FROM random_stage_images WHERE alias = $2 AND test_config_id = $1
    ) previous
    WHERE alias = $3 AND test_config_id = $1 RETURNING stage_number, row_number, column_number, alias;`,
			configId,
			selectedAlias,
			replacementAlias).Scan(&replacedImage.Stage, &replacedImage.Row, &replacedImage.Column, &replacedImage.Alias)
	} else if replacementType == "user-img" {
    db.QueryRow(`
      UPDATE random_stage_images SET alias ='user-img', image = NULL
      WHERE alias = $2 AND test_config_id = $1
      RETURNING stage_number, row_number, column_number, alias;
    `,
      configId,
      selectedAlias).Scan(&replacedImage.Stage, &replacedImage.Row, &replacedImage.Column, &replacedImage.Alias)
  } else {
		db.QueryRow(`WITH replacing AS (
      SELECT 0 as id, test_config_id, stage_number, row_number, column_number
      FROM random_stage_images WHERE alias = $2 AND test_config_id = $1
    ), collection_img AS(
      SELECT 0 as id, image, image_type
      FROM saved_images WHERE alias = $3 AND collection_id =$4
    )
    INSERT INTO random_stage_images (
      image, image_type, alias, test_config_id, stage_number, row_number, column_number, creation_date
    )
    SELECT
    collection_img.image,
    collection_img.image_type,
    replace(md5(random()::text || clock_timestamp()::text), '-'::text, ''::text)::varchar(60),
    replacing.test_config_id,
    replacing.stage_number,
    replacing.row_number,
    replacing.column_number,
    $5
    FROM collection_img JOIN replacing ON replacing.id = collection_img.id
    RETURNING stage_number, row_number, column_number, alias;`,
			configId,
			selectedAlias,
			replacementAlias,
			collectionId,
			time.Now()).Scan(&replacedImage.Stage, &replacedImage.Row, &replacedImage.Column, &replacedImage.Alias)
	}

	go db.Exec("DELETE FROM random_stage_images WHERE test_config_id = $1 AND (alias = $2 OR replacement_alias IS NOT NULL);", configId, selectedAlias)

	return replacedImage
}

func getRandomImagesByConfigId(db *sql.DB, configId int) *[]*RandomImage {
	var randomImages []*RandomImage

	rows, err := db.Query("SELECT alias, stage_number, row_number, column_number FROM random_stage_images WHERE test_config_id = $1;", configId)
	if err != nil {
		fmt.Println(err)
		// handle error
	}
	for rows.Next() {
		randomImage := new(RandomImage)
		if err := rows.Scan(&randomImage.Alias, &randomImage.Stage, &randomImage.Row, &randomImage.Column); err != nil {
			fmt.Println(err)
			// handle error
		}
		randomImages = append(randomImages, randomImage)
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		// handle error
	}
	rows.Close()

	return &randomImages
}

func getStageImagesByConfigId(db *sql.DB, configId int) *[]*RandomImage {
  var randomImages []*RandomImage

  rows, err := db.Query(`
    SELECT image.alias, stage.stage_number, image.row_number, image.column_number
    FROM test_config_stage_images image
    JOIN test_config_stages stage ON image.stage_id = stage.id
    WHERE stage.test_config_id = $1
    ORDER BY stage.stage_number ASC, image.row_number, image.column_number ASC
  `, configId)
  if err != nil {
    fmt.Println(err)
    // handle error
  }
  for rows.Next() {
    randomImage := new(RandomImage)
    if err := rows.Scan(&randomImage.Alias, &randomImage.Stage, &randomImage.Row, &randomImage.Column); err != nil {
      fmt.Println(err)
      // handle error
    }
    randomImages = append(randomImages, randomImage)
  }

  if err := rows.Err(); err != nil {
    fmt.Println(err)
    // handle error
  }
  rows.Close()

  return &randomImages
}

func GetMatrixMap(db *sql.DB, configId int, randomImages bool) *map[string]map[string]map[string]string {
	var matrixMap = make(map[string]map[string]map[string]string)
  var imagePointers *[]*RandomImage
  if (randomImages) {
    imagePointers = getRandomImagesByConfigId(db, configId)
  } else {
    imagePointers = getStageImagesByConfigId(db, configId)
  }
	for _, imagePointer := range *imagePointers {
		image := *imagePointer
		addToMatrixMap(matrixMap, strconv.Itoa(image.Stage), strconv.Itoa(image.Row), strconv.Itoa(image.Column), image.Alias)
	}

	return &matrixMap
}

func addToMatrixMap(matrixMap map[string]map[string]map[string]string, stage, row, column, alias string) {
	stageMap, ok := matrixMap[stage]
	if !ok {
		stageMap = make(map[string]map[string]string)
		matrixMap[stage] = stageMap
	}

	rowMap, ok := matrixMap[stage][row]
	if !ok {
		rowMap = make(map[string]string)
		matrixMap[stage][row] = rowMap
	}

	matrixMap[stage][row][column] = alias
}
