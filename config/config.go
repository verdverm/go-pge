package config

import (
	"bytes"
	"io"
	"log"
	"strings"
)

func ParseConfig(contents []byte, mapper func(field, value string, config interface{}) (err error), config interface{}) error {
	data := bytes.NewBuffer(contents)

	var field, value string
	var err error

	lineNum := 0
	for {
		l, buferr := data.ReadString('\n') // parse line-by-line
		lineNum++
		l = strings.TrimSpace(l)

		if buferr != nil {
			if buferr != io.EOF {
				return buferr
			}

			if len(l) == 0 {
				break
			}
		}

		if len(l) == 0 {
			continue
		} // empty line

		// strip comments
		if p := strings.Index(l, "#"); p >= 0 { // comment
			l = l[:p]
		}
		// for comments that take a whole line
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}

		// fill field and value strings
		if p := strings.Index(l, "="); p >= 0 { // comment
			field = strings.TrimSpace(l[:p-1])
			value = strings.TrimSpace(l[p+1:])
		} else {
			log.Fatalf("No '=' on line %d of config file", lineNum)
		}

		err = mapper(field, value, config)

		if err != nil {
			return err
		}
		// Reached end of data
		if buferr == io.EOF {
			break
		}
	}
	return nil
}
