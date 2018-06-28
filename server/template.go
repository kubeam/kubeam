package server

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

	"github.com/Masterminds/sprig"
	"gopkg.in/yaml.v2"
)

//LookupList describes cluster lookup list
type LookupList struct {
	Description string
	Lookup      map[string]map[string]string
}

const (
	//LookupFile ...
	LookupFile = iota + 1
	//LookupKey ...
	LookupKey = iota + 1
	//LookupSubkey ...
	LookupSubkey = iota + 1
	//LookupDefaultValue ...
	LookupDefaultValue = iota + 1
)

/*RenderTemplate reads and parses the yaml file with the application
resource values*/
func RenderTemplate(tmpfile string, pairs map[string]interface{}) string {

	file, err := os.Open(tmpfile)
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

				includeFileName := fmt.Sprintf("%s/%s", path.Dir(tmpfile), match[1])
				includeContent, err := ioutil.ReadFile(includeFileName)
				if err != nil {
					LogWarning.Println(err)
				}
				LogInfo.Println("including file :", includeFileName)
				LogDebug.Println("includeContent", string(includeContent))
				srcContent.WriteString(string(includeContent))
			} else {
				LogWarning.Println("Found incomplete tag in include from file ", tmpfile)
			}
		} else if strings.Index(t, "<%LookupFile:") > -1 {
			LogDebug.Println("Rendering LookupFile")
			var lookup LookupList
			re := regexp.MustCompile("\\<\\%LookupFile:(.*?),(.*?),(.*?),(.*?)\\%\\>")

			/*
				//
				// Fist we need to find if there is a template within the lookup definition
				t := fasttemplate.New(t, "{{", "}}")
				s := t.ExecuteString(pairs)
			*/
			//var tmpl = template.Must(template.ParseFiles(t))
			// Create a new template and parse the letter into it.
			// Get the Sprig function map.
			fmap := sprig.TxtFuncMap()
			var tmpl = template.Must(template.New("LookupFile").Funcs(fmap).Parse(t))

			var bytes bytes.Buffer
			writer := bufio.NewWriter(&bytes)

			err = tmpl.Execute(writer, pairs)
			check(err)

			err = writer.Flush()
			check(err)

			LogDebug.Println(bytes.String())

			match := re.FindStringSubmatch(bytes.String())

			if len(match) == 0 {
				LogError.Println("invalid LookupFile: syntax ", t)
				//BUG/FIX: Should push up a error to rest calling function
				continue
			}

			LogDebug.Println("LookupFile: ", match[LookupFile])
			LogDebug.Println("LookupKey: ", match[LookupKey])
			LogDebug.Println("LookupSubkey: ", match[LookupSubkey])
			LogDebug.Println("LookupDefaultValue: ", match[LookupDefaultValue])

			yamlFile, err := ioutil.ReadFile(fmt.Sprintf(match[LookupFile]))
			if err != nil {
				LogError.Println("reading LookupFile ", match[LookupFile])
				//return "", errors.New(fmt.Sprintf( "Could not lockup file: %v", match) )
			}

			err = yaml.Unmarshal(yamlFile, &lookup)
			check(err)

			var lookupvalue string
			var ok bool
			LogDebug.Println(lookup.Lookup)
			if lookupvalue, ok = lookup.Lookup[match[LookupKey]][match[LookupSubkey]]; ok {
				LogDebug.Println("Found lookup value in file :", lookupvalue)
			} else {
				lookupvalue = match[LookupDefaultValue]
				LogDebug.Println("Using default lookup Value :", lookupvalue)
			}

			srcContent.WriteString(re.ReplaceAllString(bytes.String(), lookupvalue))

		} else {
			srcContent.WriteString(t)
		}
	}

	if err := scanner.Err(); err != nil {
		LogWarning.Println(err)
	}

	//var tmpl = template.Must(template.ParseFiles(tmpl_file))
	// Get the Sprig function map.
	fmap := sprig.TxtFuncMap()
	var tmpl = template.Must(template.New("rendered_template").Funcs(fmap).Parse(srcContent.String()))

	var bytes bytes.Buffer
	writer := bufio.NewWriter(&bytes)

	err = tmpl.Execute(writer, pairs)
	check(err)

	err = writer.Flush()
	check(err)

	LogDebug.Println(bytes.String())

	return (bytes.String())

}
