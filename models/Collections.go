package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/KevinMcIntyre/password-research-server/utils"
)

type NewCollectionRequest struct {
	CollectionLabel string
}

func (request NewCollectionRequest) Validate() []string {
	errorFields := make([]string, 0)

	if utils.IsEmptyString(request.CollectionLabel) {
		errorFields = append(errorFields, "collectionLabel")
	}

	return errorFields
}

func (collection NewCollectionRequest) Save(db *sql.DB) int {
	var newCollectionId int
	err := db.QueryRow("INSERT INTO collections(label, creation_date) VALUES($1, $2) returning id", collection.CollectionLabel, time.Now()).Scan(&newCollectionId)
	if err != nil {
		fmt.Println(err)
		// handle error
	}
	return newCollectionId
}
