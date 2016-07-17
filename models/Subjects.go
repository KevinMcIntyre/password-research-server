package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/KevinMcIntyre/password-research-server/utils"
)

type NewSubjectRequest struct {
	FirstName string
	LastName  string
	Email     string
	Birthday  string
}

type SubjectNameAndId struct {
	Id   int
	Name string
}

type SubjectProfile struct {
	FirstName        string
	LastName         string
	Email            string
	Birthday         time.Time
	Password         sql.NullString
	PasswordStrength sql.NullInt64
	PasswordEntropy  sql.NullFloat64
	PinNumber        sql.NullString
}

func (request NewSubjectRequest) Validate() ([]string, time.Time) {
	errorFields := make([]string, 0)
	var birthday time.Time
	var err error
	if utils.IsEmptyString(request.FirstName) {
		errorFields = append(errorFields, "firstName")
	}

	if utils.IsEmptyString(request.LastName) {
		errorFields = append(errorFields, "lastName")
	}

	if utils.IsEmptyString(request.Email) || !utils.IsValidEmail(request.Email) {
		errorFields = append(errorFields, "email")
	}

	if utils.IsEmptyString(request.Birthday) {
		errorFields = append(errorFields, "birthday")
	} else {
		layout := "1/2/2006"
		birthday, err = time.Parse(layout, request.Birthday)

		if err != nil {
			errorFields = append(errorFields, "birthday")
		}
	}

	return errorFields, birthday
}

func SaveNewSubject(db *sql.DB, profile NewSubjectRequest, birthday time.Time) int {
	var newSubjectId int
	err := db.QueryRow("INSERT INTO subjects(email, first_name, last_name, birth_date, creation_date) VALUES($1, $2, $3, $4, $5) returning id", profile.Email, profile.FirstName, profile.LastName, birthday, time.Now()).Scan(&newSubjectId)
	if err != nil {
		fmt.Println(err)
		// handle error
	}
	return newSubjectId
}

func SaveSubjectPassword(db *sql.DB, subjectId int, password string, score int, entropy float64) {
	err := db.QueryRow("UPDATE subjects SET password = $2, password_strength = $3, password_entropy = $4 WHERE id = $1", subjectId, password, score, entropy)
	if err != nil {
		fmt.Println(err)
		// handle error
	}
}

func SaveSubjectPin(db *sql.DB, subjectId int, pinNumber string) {
	err := db.QueryRow("UPDATE subjects SET pin_number = $2 WHERE id = $1", subjectId, pinNumber)
	if err != nil {
		fmt.Println(err)
		// handle error
	}
}

func GetSubjectProfileById(db *sql.DB, subjectId int) SubjectProfile {
	var profile SubjectProfile
	err := db.QueryRow("SELECT first_name, last_name, email, birth_date, password, password_strength, password_entropy, pin_number FROM subjects WHERE id=$1 AND id != 0", subjectId).Scan(&profile.FirstName, &profile.LastName, &profile.Email, &profile.Birthday, &profile.Password, &profile.PasswordStrength, &profile.PasswordEntropy, &profile.PinNumber)
	if err != nil {
		fmt.Println(err)
		// handle error
	}
	return profile
}

func GetSubjectList(db *sql.DB) []SubjectNameAndId {
	rows, err := db.Query("SELECT id, first_name, last_name FROM subjects WHERE id != 0 ORDER BY UPPER(last_name) ASC, UPPER(first_name) ASC")
	if err != nil {
		fmt.Println(err)
		// handle error
	}

	var subjectList []SubjectNameAndId
	for rows.Next() {
		SubjectNameAndId := new(SubjectNameAndId)
		var firstName string
		var lastName string
		if err := rows.Scan(&SubjectNameAndId.Id, &firstName, &lastName); err != nil {
			fmt.Println(err)
			fmt.Println(err)
			// handle error
		}
		SubjectNameAndId.Name = lastName + ", " + firstName
		subjectList = append(subjectList, *SubjectNameAndId)
	}

	return subjectList
}
