package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

func showChart(chart string) string {
	args := []string{"show", "chart", chart}
	cmd := exec.Command("helm", args...)
	out, err := outputWithRichError(cmd)
	if err != nil {
		fmt.Println("unexpected error: err:", err)
		os.Exit(1)
	}
	return string(out)
}

func getChartVersion(chart string) string {
	re := regexp.MustCompile(`(?m)^(version:).([0-9]+\.[0-9]+\.[0-9]+(?:-[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)?(?:\+[0-9A-Za-z-]+)?)`)
	showChart := showChart(chart)
	match := re.FindStringSubmatch(showChart)
	return match[2]
}

func getChartName(chart string) string {
	re := regexp.MustCompile("(?m)^(name:).(.*)")
	showChart := showChart(chart)
	match := re.FindStringSubmatch(showChart)
	return match[2]
}

func getTemplate(chart string, version string, values []string) (string, error) {
	args := []string{"template", chart, "--version", version}
	for i, s := range values {
		args = append(args, "-f")
		args = append(args, s)
		_ = i
	}
	cmd := exec.Command("helm", args...)
	out, err := outputWithRichError(cmd)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
