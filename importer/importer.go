package importer

import (
	"bufio"
	"bytes"
	"html"
	"html/template"
	"os"
	"regexp"
	"strings"

	_ "embed"
)

//go:embed "brunorequest.gotmpl"
var TemplateVo []byte

func convertRestToBruno(fileName string, inputFile []byte, seq int) []byte {
	fileNameSplits := strings.Split(fileName, ".")
	fileName = fileNameSplits[0]
	// Create a BrunoConfig struct and populate it with the variables
	config := &BrunoConfig{
		Meta: &Meta{
			Name: fileName,
		},
		Call: &Call{},
		Body: &Body{},
	}
	// Read the input file line by line
	scanner := bufio.NewScanner(bytes.NewReader(inputFile))
	bodyStarted := false
	bodyLines := []string{}
	for scanner.Scan() {
		line := scanner.Bytes()
		lineString := html.UnescapeString(string(line))
		lineString = strings.ReplaceAll(lineString, `\"`, `"`)

		// If the line starts with a square or curly brace, or we're already in the body area,
		// append a line to the body array
		if strings.HasPrefix(lineString, "{") || strings.HasPrefix(lineString, "[") || bodyStarted {
			bodyStarted = true
			bodyLines = append(bodyLines, lineString)
			continue
		}
		// If the line is empty, skip it
		if len(line) == 0 {
			continue
		}

		// If the line starts with "POST" or "GET", set two variables to the method and the url that follows
		if strings.HasPrefix(lineString, "POST") {
			config.Meta.Verb = "post"
			// Split the line by spaces
			lineSplits := strings.Split(lineString, " ")
			// Set the url to the second element in the array
			if len(lineSplits) > 1 {
				config.Call.Url = lineSplits[1]
			}
		}
		if strings.HasPrefix(lineString, "GET") {
			config.Meta.Verb = "get"
			config.Body.Mode = "none"
			lineSplits := strings.Split(lineString, " ")
			if len(lineSplits) > 1 {
				config.Call.Url = lineSplits[1]
			}
		}
		// If the line starts with a word followed by a colon, set a variable to the headers array
		re := regexp.MustCompile(`(?i)^\w+-?\w*:`)
		match := re.FindStringSubmatch(lineString)
		if match != nil {
			if config.Headers == nil {
				config.Headers = []string{}
			}
			config.Headers = append(config.Headers, string(line))
		}

	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if len(bodyLines) > 0 {
		config.Body.Mode = "json"
		config.Body.Raw = bodyLines
	}
	// Use the BrunoConfig struct to populate the template
	config.Meta.Seq = seq
	// Convert the template to a byte array
	tmpl, err := template.New("generate-" + fileName).Parse(string(TemplateVo))
	if err != nil {
		panic(err)
	}
	var parsed bytes.Buffer
	err = tmpl.Execute(&parsed, config)
	if err != nil {
		panic(err)
	}
	// remove empty lines at the end of the parsed buffer
	for {
		if len(parsed.Bytes()) == 0 {
			break
		}
		if len(parsed.Bytes()) > 1 &&
			parsed.Bytes()[len(parsed.Bytes())-1] == 10 &&
			parsed.Bytes()[len(parsed.Bytes())-2] == 10 {
			parsed.Truncate(len(parsed.Bytes()) - 1)
		} else {
			break
		}
	}
	parsedString := html.UnescapeString(parsed.String())
	// Return the template
	return []byte(parsedString)
}

func WalkDir(inputDir, outputDir string) error {
	// Recursively walk the input directory
	// For each file, read the file and convert it to Bruno using convertRestToBruno
	// Write the output to the output directory creating the same directory structure
	// as the input directory
	// If the output directory does not exist, create it
	inputDirHandle, err := os.Open(inputDir)
	if err != nil {
		return err
	}
	defer inputDirHandle.Close()
	seq := 1
	for {
		// Read the next file in the input directory
		file, err := inputDirHandle.Readdir(1)
		if err != nil {
			return err
		}
		if len(file) == 0 {
			break
		}
		fileName := file[0].Name()
		// if the filename starts with a dot, skip it
		if fileName[0] == '.' {
			continue
		}
		if file[0].IsDir() {
			// If the file is a directory,  create the destination directory if it doesn't exist
			// and recursively call walkDir
			// check if the output directory exists

			_, err := os.Stat(outputDir + "/" + fileName)
			if os.IsNotExist(err) {
				err := os.MkdirAll(outputDir+"/"+fileName, 0o755)
				if err != nil {
					return err
				}
			}
			err = WalkDir(inputDir+"/"+fileName, outputDir+"/"+fileName)
			if err != nil {
				return err
			}
		} else {
			// if the file does not end with .http, skip it
			if fileName[len(fileName)-5:] != ".http" {
				continue
			}
			// If the file is a regular file, read the file and convert it to Bruno
			inputFile, err := os.ReadFile(inputDir + "/" + fileName)
			if err != nil {
				return err
			}
			outputFile := convertRestToBruno(fileName, inputFile, seq)
			seq++
			// Write the output to the output directory creating the same directory structure
			// as the input directory
			// If the output directory does not exist, create it
			fileNameSplits := strings.Split(fileName, ".")
			fileName = fileNameSplits[0] + ".bru"
			err = os.WriteFile(outputDir+"/"+fileName, outputFile, 0o644)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
