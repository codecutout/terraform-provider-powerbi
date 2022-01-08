package pbixrewriter

import (
	"archive/zip"
	"io"
	"os"
)

// PipelineFunc defines a unit of processing for rewriting PBIX files
type PipelineFunc func(file *zip.File, reader io.Reader, next PipelineFuncNext) error

// PipelineFuncNext defines the next function inside a PipelineFunc
type PipelineFuncNext func(file *zip.File, reader io.Reader) error

// RewritePbix rewrites a PBIX file applying the given PipelineFuncs
func RewritePbix(input *zip.Reader, output *zip.Writer, pipelineFuncs []PipelineFunc) error {
	for _, inputItem := range input.File {
		inputItemReader, err := inputItem.Open()
		if err != nil {
			return err
		}
		defer inputItemReader.Close()

		pipeline := nestPipelineFunc(0, append(pipelineFuncs, buildWriterPipelineFunc(output)))

		pipeline(inputItem, inputItemReader)

	}
	return nil
}

// RewritePbixFiles rewrites a PBIX file applying the given PipelineFuncs
func RewritePbixFiles(inputPbixFile string, outputPbixFile string, pipelineFuncs []PipelineFunc) error {
	zipReader, err := zip.OpenReader(inputPbixFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	targetFile, err := os.Create(outputPbixFile)
	if err != nil {
		return err
	}
	targetZipWriter := zip.NewWriter(targetFile)
	defer targetZipWriter.Close()

	return RewritePbix(&zipReader.Reader, targetZipWriter, pipelineFuncs)
}

func buildWriterPipelineFunc(writer *zip.Writer) PipelineFunc {

	return func(file *zip.File, reader io.Reader, next PipelineFuncNext) error {
		header, err := zip.FileInfoHeader(file.FileInfo())
		if err != nil {
			return err
		}
		header.Name = file.Name

		outputItemWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(outputItemWriter, reader)
		if err != nil {
			return err
		}
		return nil
	}
}

func nestPipelineFunc(index int, pipelineFuncs []PipelineFunc) PipelineFuncNext {
	return func(file *zip.File, reader io.Reader) error {
		nextIndex := index + 1
		if nextIndex >= len(pipelineFuncs) {
			return pipelineFuncs[index](file, reader, func(file *zip.File, reader io.Reader) error { return nil })
		}
		return pipelineFuncs[index](file, reader, nestPipelineFunc(index+1, pipelineFuncs))
	}
}
