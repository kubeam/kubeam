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

	"gopkg.in/yaml.v2"
)

type LookupList struct {
	Description string
	Lookup      map[string]map[string]string
}

const (
	LOOKUP_FILE          = iota + 1
	LOOKUP_KEY           = iota + 1
	LOOKUP_SUBKEY        = iota + 1
	LOOKUP_DEFAULT_VALUE = iota + 1
)

func render_template(tmpl_file string, pairs map[string]interface{}) string {

	file, err := os.Open(tmpl_file)
	if err != nil {
		LogWarning.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	//var srcContent string
	var srcContent bytes.Buffer
	for scanner.Scan() {
		t := fmt.Sprintln(scanner.Text())
		if strings.Index(t, "<%file:") > -1 {
			LogDebug.Println("Including file external file")
			if strings.Index(t, "%>") > -1 {
				re := regexp.MustCompile("\\<\\%file:(.*?)\\%\\>")
				match := re.FindStringSubmatch(t)
				if len(match) == 0 {
					LogError.Println("invalid file: syntax ", t)
					continue
				}

				includeFileName := fmt.Sprintf("%s/%s", path.Dir(tmpl_file), match[1])
				includeContent, err := ioutil.ReadFile(includeFileName)
				if err != nil {
					LogWarning.Println(err)
				}
				LogInfo.Println("including file :", includeFileName)
				LogDebug.Println("includeContent", string(includeContent))
				srcContent.WriteString(string(includeContent))
			} else {
				LogWarning.Println("Found incomplete tag in include from file ", tmpl_file)
			}
		} else if strings.Index(t, "<%lookup_file:") > -1 {
			LogDebug.Println("Rendering lookup_file")
			var lookup LookupList
			re := regexp.MustCompile("\\<\\%lookup_file:(.*?),(.*?),(.*?),(.*?)\\%\\>")

			/*
				//
				// Fist we need to find if there is a template within the lookup definition
				t := fasttemplate.New(t, "{{", "}}")
				s := t.ExecuteString(pairs)
			*/
			//var tmpl = template.Must(template.ParseFiles(t))
			// Create a new template and parse the letter into it.
			var tmpl = template.Must(template.New("lookup_file").Parse(t))

			var bytes bytes.Buffer
			writer := bufio.NewWriter(&bytes)

			err = tmpl.Execute(writer, pairs)
			check(err)

			err = writer.Flush()
			check(err)

			LogDebug.Println(bytes.String)

			match := re.FindStringSubmatch(bytes.String())

			if len(match) == 0 {
				LogError.Println("invalid lookup_file: syntax ", t)
				//BUG/FIX: Should push up a error to rest calling function
				continue
			}

			LogDebug.Println("lookup_file: ", match[LOOKUP_FILE])
			LogDebug.Println("lookup_key: ", match[LOOKUP_KEY])
			LogDebug.Println("lookup_subkey: ", match[LOOKUP_SUBKEY])
			LogDebug.Println("lookup_default_value: ", match[LOOKUP_DEFAULT_VALUE])

			yamlFile, err := ioutil.ReadFile(fmt.Sprintf(match[LOOKUP_FILE]))
			if err != nil {
				LogError.Println("reading lookup_file ", match[LOOKUP_FILE])
				//return "", errors.New(fmt.Sprintf( "Could not lockup file: %v", match) )
			}

			err = yaml.Unmarshal(yamlFile, &lookup)
			check(err)

			var lookup_value string
			var ok bool
			LogDebug.Println(lookup.Lookup)
			if lookup_value, ok = lookup.Lookup[match[LOOKUP_KEY]][match[LOOKUP_SUBKEY]]; ok {
				LogDebug.Println("Found lookup value in file :", lookup_value)
			} else {
				lookup_value = match[LOOKUP_DEFAULT_VALUE]
				LogDebug.Println("Using default lookup Value :", lookup_value)
			}

			srcContent.WriteString(re.ReplaceAllString(bytes.String(), lookup_value))

		} else {
			srcContent.WriteString(t)
		}
	}

	if err := scanner.Err(); err != nil {
		LogWarning.Println(err)
	}

	//var tmpl = template.Must(template.ParseFiles(tmpl_file))
	var tmpl = template.Must(template.New("rendered_template").Parse(srcContent.String()))

	var bytes bytes.Buffer
	writer := bufio.NewWriter(&bytes)

	err = tmpl.Execute(writer, pairs)
	check(err)

	err = writer.Flush()
	check(err)

	LogDebug.Println(bytes.String())

	return (bytes.String())

}
