package models
import (
  "fmt"
  "database/sql"
  "time"
)

type ImageUploadResponse struct {
  SubjectId    int
  CollectionId int
  Aliases      []string
}

type ImageDiscardRequest struct {
  SubjectId    int
  CollectionId int
  ImageAlias   string
}

func DiscardImage(db *sql.DB, subjectId int, collectionId int, imageAlias string) {
  _, err := db.Exec("DELETE FROM uploaded_images WHERE subject_id = $1 AND collection_id = $2 AND alias = $3", subjectId, collectionId, imageAlias)
  if err != nil {
    fmt.Println(err)
    // handle error
  }
}

func SaveUpload(db *sql.DB, uploaded bool, subjectId int, collectionId int, imageBytes []byte, imageType string) string {
  var alias string
  err := db.QueryRow("INSERT INTO uploaded_images (subject_id, collection_id, image, image_type, alias, creation_date) VALUES($1, $2, $3, $4, replace(md5(random()::text || clock_timestamp()::text), '-'::text, ''::text)::varchar(60), $5) returning alias;",
    subjectId,
    collectionId,
    imageBytes,
    imageType,
    time.Now(),
    ).Scan(&alias)
  if err != nil {
    fmt.Println(err)
    // handle error
  }
  return alias
}
