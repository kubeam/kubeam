package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	//"gopkg.in/gcfg.v1"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Description string
	Templates   map[string]map[string]string
}

func main() {
	// Init Loggers:
	// File descriptors in order: Trace, Debug, Info, Warning, Error
	// set to ** ioutil.Discard ** to stop recording those logs
	InitLogger(ioutil.Discard, ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	verbose := flag.Bool("v", false, "Turn on verbose")
	verboseplus := flag.Bool("vv", false, "Turn on more verbose")
	flag.Parse()

	for _, templateConfigFile := range flag.Args() {
		if *verbose {
			fmt.Printf("Processing %v\n", templateConfigFile)
		}

		filename, _ := filepath.Abs(templateConfigFile)

		yamlFile, err := ioutil.ReadFile(filename)
		check(err)

		var config Config
		mymap := make(map[string]interface{})

		err = yaml.Unmarshal(yamlFile, &config)
		check(err)

		fmt.Printf("Description: %#v\n", config.Description)
		for key, value := range config.Templates {
			if *verbose {
				fmt.Printf("###\nFile :[%v]\nTemplate :[%v]\nDatasource :[%v]\nDestination :[%v]\nMode :[%v]\n", key, value["template"], value["datasource"], value["destination"], value["mode"])
			}
			if (value["action"] == "delete") || (value["action"] == "erase") {
				fmt.Printf("Removing file %v", value["destination"])
				err := os.Remove(value["destination"])
				if err != nil {
					log.Fatal(err)
					check(err)
				}
			} else if value["action"] == "copy" {
				if _, ok := value["source"]; ok {
					copyFile(value["source"], value["destination"])
					if *verbose {
						fmt.Printf("Copied [%v] --> [%v]\n", value["source"], value["destination"])
					}
				} else {
					panic("action [copy] requires parameter [source]")
				}
			} else {
				if _, ok := value["sourcetype"]; ok {
					if strings.HasPrefix(value["sourcetype"], "env") {
						LoadEnv(mymap)
					}
					// default to config type is a .ini
				} else {
					err := Loadconfig(value["datasource"], mymap)
					if err != nil {
						log.Fatal(err)
						check(err)
					}
				}

				rendered := []byte(render_template(value["template"], mymap))

				if *verboseplus {
					fmt.Println(string(rendered))
				}
				var fileMode int64
				fileMode, err = strconv.ParseInt(value["mode"], 8, 0)
				check(err)
				err = ioutil.WriteFile(value["destination"], rendered, os.FileMode(fileMode))
				check(err)
				fmt.Printf(" Wrote : %v\n", value["destination"])
			}
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
