package main

import (
  "image"
  "github.com/disintegration/imaging"
  "strconv"
)

func main() {

  files := []string{"test-images/01.jpg", "test-images/02.jpg", "test-images/03.jpg", "test-images/04.jpg", "test-images/05.jpg", "test-images/06.jpg"}

  var thumbnails []image.Image
  for _, file := range files {
    img, err := imaging.Open(file)
    if err != nil {
      panic(err)
    }
    croppedThumb := imaging.CropCenter(imaging.Fit(img, 200, 200, imaging.Linear), 120, 120)
    thumbnails = append(thumbnails, croppedThumb)
  }

  for i, thumb := range thumbnails {
    err := imaging.Save(thumb, "result-images/0" + strconv.Itoa(i + 1) + ".jpg")
    if err != nil {
      panic(err)
    }
  }
}
