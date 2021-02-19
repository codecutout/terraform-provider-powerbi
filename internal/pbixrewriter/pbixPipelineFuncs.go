package pbixrewriter

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
)

// SetDatasetIDPipelineFunc modifies the dataset used within a PBIX
func SetDatasetIDPipelineFunc(datasetID string) PipelineFunc {
	initialCatalogRegex := regexp.MustCompile(`Initial Catalog=[^;]*;`)
	pbiModelDatabaseNameRegex := regexp.MustCompile(`"PbiModelDatabaseName":"[^"]*"`)

	return func(file *zip.File, reader io.Reader, next PipelineFuncNext) error {
		if file.Name == "SecurityBindings" {
			// security bindings needs to be deleted, otherwise file appears corrupt
			// if opening file on machine it was previously opened on
			return nil
		}

		if file.Name == "Connections" {
			bytes, err := ioutil.ReadAll(reader)
			if err != nil {
				return err
			}
			connections := string(bytes)
			connections = initialCatalogRegex.ReplaceAllString(connections, fmt.Sprintf(`Initial Catalog=%s;`, datasetID))
			connections = pbiModelDatabaseNameRegex.ReplaceAllString(connections, fmt.Sprintf(`"PbiModelDatabaseName":"%s"`, datasetID))
			return next(file, strings.NewReader(connections))
		}

		return next(file, reader)
	}
}
