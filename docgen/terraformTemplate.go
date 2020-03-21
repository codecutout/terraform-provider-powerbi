package docgen

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"os"
	"path/filepath"
	"strings"
)

// PopulateTerraformDocs update template fields inline in files in the folderpath
func PopulateTerraformDocs(folderpath string, providerName string, provider *schema.Provider) error {
	err := filepath.Walk(folderpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		fileName := info.Name()
		fileNameWithoutExtension := fileName[0:strings.LastIndex(fileName, ".")]

		if fileNameWithoutExtension == providerName {
			populateProviderDoc(path, provider)
		}
		for resourceName, resourceValue := range provider.ResourcesMap {
			if fileNameWithoutExtension == resourceName {
				populateResourceDoc(path, resourceValue)
			}
		}

		return nil

	})
	if err != nil {
		return nil
	}
	return nil
}

func populateProviderDoc(filepath string, provider *schema.Provider) error {
	return populateTemplate(filepath, struct {
		NonComputedParameters string
		ComputedParameters    string
	}{
		propertyDocumentation(provider.Schema, func(propKey string, propValue *schema.Schema) bool {
			return !propValue.Computed
		}),
		propertyDocumentation(provider.Schema, func(propKey string, propValue *schema.Schema) bool {
			return propValue.Computed
		}),
	})
}

func populateResourceDoc(filepath string, resource *schema.Resource) error {
	return populateTemplate(filepath, struct {
		NonComputedParameters string
		ComputedParameters    string
	}{
		propertyDocumentation(resource.Schema, func(propKey string, propValue *schema.Schema) bool {
			return !propValue.Computed
		}),
		propertyDocumentation(resource.Schema, func(propKey string, propValue *schema.Schema) bool {
			return propValue.Computed
		}),
	})
}

func propertyDocumentation(propertySchemas map[string]*schema.Schema, filter func(propKey string, propValue *schema.Schema) bool) string {
	builder := &strings.Builder{}
	writePropertyDocumentation(builder, propertySchemas, filter)
	return strings.Trim(builder.String(), " \r\n")
}
func writePropertyDocumentation(writer *strings.Builder, propertySchemas map[string]*schema.Schema, filter func(propKey string, propValue *schema.Schema) bool) {

	filteredPropertySchames := make(map[string]*schema.Schema)
	for key, value := range propertySchemas {
		if filter(key, value) {
			filteredPropertySchames[key] = value
		}
	}

	nestedObjects := make(map[string]map[string]*schema.Schema)

	for filteredProperty, filteredPropertySchema := range filteredPropertySchames {

		descriptionSuffix := ""
		res, isNestedResource := filteredPropertySchema.Elem.(*schema.Resource)
		if isNestedResource {
			nestedObjects[filteredProperty] = res.Schema
			descriptionSuffix = capitilizeFirstCharacter(indefiniteArticle(filteredProperty)) + " `" + filteredProperty + "` block is defined below."
		}

		tagString := buildTagString(filteredPropertySchema)
		if tagString != "" {
			tagString = tagString + " "
		}

		writeLine(writer, "* `", filteredProperty, "` - ", tagString, joinSentances(filteredPropertySchema.Description, descriptionSuffix))

	}

	for argumentName, nestedObject := range nestedObjects {
		writeLine(writer, "---")
		writeLine(writer, capitilizeFirstCharacter(indefiniteArticle(argumentName)), " `", argumentName, "` block supports the following:")
		writePropertyDocumentation(writer, nestedObject, filter)
	}
}

func buildTagString(attribute *schema.Schema) string {
	var tags []string

	if attribute.Optional {
		tags = append(tags, "Optional")
		if attribute.Default != "" && attribute.Default != nil {
			tags = append(tags, fmt.Sprintf("Default: `%v`", attribute.Default))
		}
	} else if attribute.Required {
		tags = append(tags, "Required")
	}

	if attribute.ForceNew {
		tags = append(tags, "Forces new resource")
	}

	if len(tags) == 0 {
		return ""
	}

	joinedTags := strings.Join(tags, ", ")
	return "(" + joinedTags + ")"
}
