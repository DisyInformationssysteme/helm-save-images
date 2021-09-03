package cmd

import (
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strings"
)

var (
	imageString string
)

func getImageList(chart string, version string, values []string) []string {
	template, err := getTemplate(chart, version, values)
	if err != nil {
		fmt.Println(err)
	}

	re := regexp.MustCompile("image: (.*)")
	regexImageMatches := re.FindAllStringSubmatch(template, -1)

	// extract simple slice
	allImages := make([]string, 0)
	for _, row := range regexImageMatches {
		imageString = trimQuotes(row[1])
		allImages = append(allImages, imageString)
	}

	// filter for unique images and sort them alphabetically
	uniqueImages := uniqueNonEmptyElementsOf(allImages)
	sort.Strings(uniqueImages)

	return uniqueImages
}

func run(chart string, version string, values []string, DryRun bool) {
	imageList := getImageList(chart, version, values)
	outputFileName := getAllImagesFilename(chart)

	if DryRun {
		printThoseImages(imageList)
	} else {
		pullImages(imageList)
		saveImages(imageList, outputFileName)
	}
}

func pullImages(imageList []string) error {
	for _, image := range imageList {
		if image != "" {
			args := []string{"pull", "-q", image}
			fmt.Printf("Pulling %s\n", image)
			cmd := exec.Command("docker", args...)

			out, err := outputWithRichError(cmd)
			if err != nil {
				fmt.Printf("%s\n", err)
				return err
			}
			fmt.Printf("%s\n", out)
		}
	}
	return nil
}

func saveImages(imageList []string, outputFileName string) error {
	if (strings.TrimSpace(strings.Join(imageList, "")) != "") && (outputFileName != "") {
		args := []string{"save", "-o", outputFileName}
		args = append(args, imageList...)
		fmt.Printf("Save images to %s\n", outputFileName)
		cmd := exec.Command("docker", args...)
		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("%s\n", err)
			return err
		}
		fmt.Printf("%s\n", stdoutStderr)
	}
	return nil
}

func getAllImagesFilename(chart string) string {
	return getChartName(chart) + "-" + getChartVersion(chart) + "-images.tar"
}

func printThoseImages(imageList []string) {
	fmt.Println("Your Helm Chart contains the following images:")
	fmt.Println()
	for _, image := range imageList {
		fmt.Printf("%+v\n", image)
	}
}

func uniqueNonEmptyElementsOf(s []string) []string {
	unique := make(map[string]bool, len(s))
	us := make([]string, len(unique))
	for _, elem := range s {
		if len(elem) != 0 {
			if !unique[elem] {
				us = append(us, elem)
				unique[elem] = true
			}
		}
	}
	return us
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if c := s[len(s)-1]; s[0] == c && (c == '"' || c == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
