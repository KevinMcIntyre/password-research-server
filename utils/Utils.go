package utils

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/disintegration/imaging"
)

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
