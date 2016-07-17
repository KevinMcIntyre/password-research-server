package models
import (
  "database/sql"
  "fmt"
)

type CollectionLabelAndId struct {
  Id    int
  Label string
}

func GetCollectionList(db *sql.DB) ([]CollectionLabelAndId) {
  rows, err := db.Query("SELECT id, label FROM collections WHERE id != 0 ORDER BY UPPER(label) ASC")
  if (err != nil) {
    fmt.Println(err);
    // handle error
  }

  var collectionList []CollectionLabelAndId
  defer rows.Close()
  for rows.Next() {
    CollectionLabelAndId := new(CollectionLabelAndId)
    if err := rows.Scan(&CollectionLabelAndId.Id, &CollectionLabelAndId.Label); err != nil {    fmt.Println(err);
      fmt.Println(err);
      // handle error
    }
    collectionList = append(collectionList, *CollectionLabelAndId)
  }

  return collectionList
}
