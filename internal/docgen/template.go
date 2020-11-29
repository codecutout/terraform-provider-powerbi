package docgen

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

func populateTemplate(filename string, model interface{}) error {
	template, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	templateString := string(template)
	re, err := regexp.Compile(`(?s)<!--\s*docgen:([a-zA-Z][a-zA-Z0-9]*)\s*-->(.*?)<!--\s*/docgen\s*-->`)
	if err != nil {
		return err
	}

	modelMap, err := toMap(model)
	if err != nil {
		return err
	}
	newString := re.ReplaceAllStringFunc(templateString, func(match string) string {
		modelProp := re.FindStringSubmatch(match)[1]
		val, ok := modelMap[modelProp]
		if ok {
			return fmt.Sprintf(`<!-- docgen:%s -->
%v
<!-- /docgen -->`, modelProp, val)
		}
		return fmt.Sprintf(`<!-- docgen:%s -->
<!-- /docgen -->`, modelProp)
	})
	return writeToFile(filename, newString)
}
