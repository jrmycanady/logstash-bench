package bench

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
)

// BuildConfig builds a new testing config based on the filter provided.
func BuildConfig(codec string, inputFilePath string, outputFilePath string, filterFilePath string, sinceDBPath string) (string, error) {
	input := BuildInput(codec, inputFilePath, sinceDBPath)
	output := BuildOutput(outputFilePath)
	filter, err := FilterFromFile(filterFilePath)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s %s", input, filter, output), nil
}

// BuildInput builds the input configuration.
func BuildInput(codec string, path string, sinceDBPath string) string {
	return fmt.Sprintf("input {file {codec => \"%s\" mode => \"tail\" path => \"%s\" sincedb_path => \"%s\" start_position => \"beginning\"} }", codec, path, sinceDBPath)
}

// BuildOutput builds the output configuration.
func BuildOutput(path string) string {
	return fmt.Sprintf("output {file {path => \"%s\"} }", path)
}

// FilterFromFile creates a filter based on the file provided. This includes
// addting the timestamp for performance testing.
func FilterFromFile(path string) (string, error) {
	var output bytes.Buffer

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	sData := string(data)
	s := bufio.NewScanner(strings.NewReader(sData))
	s.Split(bufio.ScanRunes)

	found := false

	for s.Scan() {
		t := s.Text()
		if t == "\n" || t == "\r" {
			continue
		}
		output.WriteString(t)
		if s.Text() == "{" && !found {
			found = true
			output.WriteString(`ruby { code => "event.set('processed_at', Time.now());"} `)
		}
	}

	return output.String(), nil

}
