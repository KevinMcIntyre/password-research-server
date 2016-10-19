package utils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
)

func WritePid() {
	err := ioutil.WriteFile("SERVER_PID", []byte(strconv.Itoa(os.Getpid())), 0644)
	if err != nil {
		log.Fatalf("Error writing SERVER_PID\n")
	}
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func RoundUp(input float64, places int) float64 {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Ceil(digit)
	return round / pow
}

func ParsePGArray(array string) ([]string, error) {
	var out []string
	var arrayOpened, quoteOpened, escapeOpened bool
	item := &bytes.Buffer{}
	for _, r := range array {
		switch {
		case !arrayOpened:
			if r != '{' {
				return nil, errors.New("Doesn't appear to be a postgres array.  Doesn't start with an opening curly brace.")
			}
			arrayOpened = true
		case escapeOpened:
			item.WriteRune(r)
			escapeOpened = false
		case quoteOpened:
			switch r {
			case '\\':
				escapeOpened = true
			case '"':
				quoteOpened = false
				if item.String() == "NULL" {
					item.Reset()
				}
			default:
				item.WriteRune(r)
			}
		case r == '}':
			// done
			out = append(out, item.String())
			return out, nil
		case r == '"':
			quoteOpened = true
		case r == ',':
			// end of item
			out = append(out, item.String())
			item.Reset()
		default:
			item.WriteRune(r)
		}
	}
	return nil, errors.New("Doesn't appear to be a postgres array.  Premature end of string.")
}

func ImageFileToByteArray(fileName string) []byte {
	file, err := os.Open(fileName)

	if err != nil {
		fmt.Println(err)
		// handle error
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	bytes := make([]byte, size)

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)

	return bytes
}

func ResizeImage(fileName string) {
	img, err := imaging.Open(fileName)
	croppedImg := imaging.CropCenter(imaging.Fit(img, 200, 200, imaging.Linear), 120, 120)
	err = imaging.Save(croppedImg, fileName)
	if err != nil {
		// handle error
	}
}

func IsEmptyString(str string) bool {
	if len(strings.TrimSpace(str)) < 1 {
		return true
	}
	return false
}

func IsValidEmail(email string) bool {
	Re := regexp.MustCompile(`.+\@.+\..+`)
	return Re.MatchString(email)
}

func IntInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func MsToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(0, msInt*int64(time.Millisecond)), nil
}

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
