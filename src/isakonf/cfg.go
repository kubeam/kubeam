package main

import (
	"errors"
	"os"
	"regexp"
	"strings"
)

//func init() {
//}

// Load adds or updates entries in an existing map with string keys
// and string values using a configuration file.
//
// The filename paramter indicates the configuration file to load ...
// the dest parameter is the map that will be updated.
//
// The configuration file entries should be constructed in key=value
// syntax.  A # symbol at the beginning of a line indicates a comment.
// Blank lines are ignored.
func Loadconfig(filename string, dest map[string]interface{}) error {
	var re *regexp.Regexp
	var pat = "[#].*\\n|\\s+\\n|\\S+[=]|.*\n"
	re, _ = regexp.Compile(pat)

	fi, err := os.Stat(filename)
	if err != nil {
		return err
	}
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	buff := make([]byte, fi.Size())
	f.Read(buff)
	f.Close()
	str := string(buff)
	if !strings.HasSuffix(str, "\n") {
		str += "\n"
	}
	s2 := re.FindAllString(str, -1)

	for i := 0; i < len(s2); {
		if strings.HasPrefix(s2[i], "#") {
			i++
		} else if strings.HasSuffix(s2[i], "=") {
			key := strings.ToLower(s2[i])[0 : len(s2[i])-1]
			i++
			if strings.HasSuffix(s2[i], "\n") {
				val := s2[i][0 : len(s2[i])-1]
				if strings.HasSuffix(val, "\r") {
					val = val[0 : len(val)-1]
				}
				i++
				dest[key] = val
			}
		} else if strings.Index(" \t\r\n", s2[i][0:1]) > -1 {
			i++
		} else {
			return errors.New("Unable to process line in cfg file containing " + s2[i])
		}
	}
	return nil
}
