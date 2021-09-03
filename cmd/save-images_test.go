package cmd

import (
	"os"
	"testing"
)

func TestGetImageList(t *testing.T) {
	// TODO
}

func TestSaveImages(t *testing.T) {
	testTable := []struct {
		imageList      []string
		outputFilename string
	}{
		{[]string{""}, ""},
		{[]string{"library/alpine:3.14.2"}, ""},
		{[]string{"library/alpine:3.14.2"}, "alpine-single-image.tar"},
		{[]string{"library/alpine:3.14.2", "library/alpine:3.13.6"}, "alpine-two-images.tar"},
	}
	for _, table := range testTable {

		pullImagesErr := pullImages(table.imageList)
		if pullImagesErr != nil {
			t.Errorf("docker pull was not successfull")
		}
		saveImageErr := saveImages(table.imageList, table.outputFilename)
		if saveImageErr != nil {
			t.Errorf("docker save was not successfull")
		}
		os.Remove(table.outputFilename)
	}
}
