package upload

import (
	"errors"
	"os"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadSingleFile(
	c *gin.Context,
	fieldName string,
	filePath string,
) error {
	file, err := c.FormFile(fieldName)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		return err
	}

	return c.SaveUploadedFile(file, path.Join(filePath, filepath.Base(file.Filename)))
}

func UploadMultipleFile(
	c *gin.Context,
	fieldName string,
	filePath string,
) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	files, ok := form.File[fieldName]
	if !ok {
		return errors.New("couldn't any files to uploaded")
	}

	err = os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		return err
	}

	errs := make([]*error, 0, len(files))
	for _, file := range files {
		err := c.SaveUploadedFile(file, path.Join(filePath, filepath.Base(file.Filename)))

		if err != nil {
			errs = append(errs, &err)
		}
	}

	if len(errs) > 0 {
		errStr := ""
		for _, err := range errs {
			errStr += (*err).Error() + "\n"
		}
		return errors.New(errStr)
	}

	return nil
}
