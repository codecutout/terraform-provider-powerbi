package docgen

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type property struct {
	Name   string
	Schema *schema.Schema
}

type nestedObjectProperty struct {
	Name         string
	NestedSchema map[string]*schema.Schema
}

// PopulateTerraformDocs update template fields inline in files in the folderpath
func PopulateTerraformDocs(folderpath string, providerName string, provider *schema.Provider) error {
	err := filepath.Walk(folderpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		fileName := info.Name()
		fileNameWithoutExtension := fileName[0:strings.LastIndex(fileName, ".")]

		if fileNameWithoutExtension == "index" {
			populateProviderDoc(path, provider)
		}
		for resourceName, resourceValue := range provider.ResourcesMap {
			resourceNameWithoutProvider := resourceName[strings.Index(resourceName, "_")+1:]
			if resourceNameWithoutProvider == fileNameWithoutExtension || resourceName == fileNameWithoutExtension {
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

	sortedProperties := sortProperties(filteredPropertySchames)

	nestedObjects := make([]nestedObjectProperty, 0)

	for _, prop := range sortedProperties {

		descriptionSuffix := ""
		res, isNestedResource := prop.Schema.Elem.(*schema.Resource)
		if isNestedResource {
			nestedObjects = append(nestedObjects, nestedObjectProperty{Name: prop.Name, NestedSchema: res.Schema})
			descriptionSuffix = capitilizeFirstCharacter(indefiniteArticle(prop.Name)) + " [`" + prop.Name + "`](#prop-" + prop.Name + ") block is defined below."
		}

		tagString := buildTagString(prop.Schema)
		if tagString != "" {
			tagString = tagString + " "
		}

		writeLine(writer, "* `", prop.Name, "` - ", tagString, joinSentances(prop.Schema.Description, descriptionSuffix))

	}

	for _, nestedProp := range nestedObjects {
		writeLine(writer, "---")
		writeLine(writer,
			"<a id=\"prop-", nestedProp.Name, "\"></a>",
			capitilizeFirstCharacter(indefiniteArticle(nestedProp.Name)),
			" `", nestedProp.Name, "` block supports the following:")
		writePropertyDocumentation(writer, nestedProp.NestedSchema, filter)
	}
}

func sortProperties(propertySchemas map[string]*schema.Schema) []property {
	propertyList := make([]property, 0, len(propertySchemas))
	for name, schema := range propertySchemas {
		propertyList = append(propertyList, property{Name: name, Schema: schema})
	}

	// Property names are sorted last
	sort.SliceStable(propertyList, func(i, j int) bool {
		return propertyList[i].Name < propertyList[j].Name
	})

	// ForceNew fields are sorted second
	sort.SliceStable(propertyList, func(i, j int) bool {
		return propertyList[i].Schema.ForceNew && !propertyList[j].Schema.ForceNew
	})

	// Required fields are sorted first
	sort.SliceStable(propertyList, func(i, j int) bool {
		return propertyList[i].Schema.Required && !propertyList[j].Schema.Required
	})

	return propertyList
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
