package docgen

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"strings"
)

func writeLine(writer io.StringWriter, values ...string) {
	for _, value := range values {
		writer.WriteString(value)
	}
	writer.WriteString("\r\n")
}

func indefiniteArticle(noun string) string {
	//not always accurate, but good enough for our use case
	isMatch, _ := regexp.MatchString("^([aeio]|un|ul|hour)", noun)
	if isMatch {
		return "an"
	}
	return "a"
}

func joinSentances(sentances ...string) string {
	var cleanedSentances []string
	for _, sentance := range sentances {
		var cleanedSentance = capitilizeFirstCharacter(strings.Trim(sentance, ". "))
		if cleanedSentance != "" {
			cleanedSentances = append(cleanedSentances, cleanedSentance)
		}
	}

	joinedSentances := strings.Join(cleanedSentances, ". ")
	if joinedSentances != "" {
		joinedSentances += "."
	}
	return joinedSentances
}

func capitilizeFirstCharacter(str string) string {
	runes := []rune(str)
	if len(runes) == 0 {
		return ""
	} else if len(runes) == 1 {
		return strings.ToUpper(string(runes[0]))
	}
	return strings.ToUpper(string(runes[0])) + string(runes[1:])
}

func toMap(in interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMap only accepts structs; got %T", v)
	}

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanInterface() {
			out[v.Type().Field(i).Name] = v.Field(i).Interface()
		}

	}
	return out, nil
}

func writeToFile(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}
