package template

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
)

var (
	ErrTemplateProcessing = errors.New("err happened during template processing")
)

// Process performs provided template processing with the provided content details.
func Process(dir, templateName, outputName string, data interface{}) error {
	templatePath := filepath.Join(dir, templateName)
	outputPath := filepath.Join(dir, outputName)

	tmpl, err := template.New(templateName).ParseFiles(templatePath)
	if err != nil {
		return errors.Wrap(err, ErrTemplateProcessing.Error())
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return errors.Wrap(err, ErrTemplateProcessing.Error())
	}

	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return errors.Wrap(err, ErrTemplateProcessing.Error())
	}

	return nil
}
