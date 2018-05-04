package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"
)

func render_template(tmpl_file string, pairs map[string]interface{}) string {
	file, err := os.Open(tmpl_file)
	if err != nil {
		LogWarning.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var srcContent bytes.Buffer
	for scanner.Scan() {
		t := fmt.Sprintln(scanner.Text())
		if strings.Index(t, "<%file:") > -1 {
			if strings.Index(t, "%>") > -1 {
				re := regexp.MustCompile("\\<\\%file:(.*?)\\%\\>")
				match := re.FindStringSubmatch(t)

				includeFileName := fmt.Sprintf("%s/%s", path.Dir(tmpl_file), match[1])
				includeContent, err := ioutil.ReadFile(includeFileName)
				if err != nil {
					LogWarning.Println(err)
				}
				LogInfo.Println("including file :", includeFileName)
				srcContent.WriteString(string(includeContent))
			} else {
				LogWarning.Println("Found incomplete tag in include from file ", tmpl_file)
			}
		} else {
			srcContent.WriteString(t)
		}
	}

	if err := scanner.Err(); err != nil {
		LogWarning.Println(err)
	}

	pairs["hostname"], err = os.Hostname()
	check(err)

	var tmpl = template.Must(template.ParseFiles(tmpl_file))

	var bytes bytes.Buffer
	writer := bufio.NewWriter(&bytes)

	err = tmpl.Execute(writer, pairs)
	check(err)

	err = writer.Flush()
	check(err)

	LogDebug.Println(bytes.String)

	return (bytes.String())

}
