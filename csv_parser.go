package csv_parser

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
)

// MultipartWriter converts the contents of the passed/received file into
// a multipart.Writer form file so the data can be passed to next service.
// It returns the correct HTTP Content-Type, the form content as bytes, and a
// possible error.
func MultipartWriter(formFile *multipart.FileHeader, fileName string, createdBy string) (string, bytes.Buffer, error) {
	// Create a new Form with file and the correct Content-Type
	// (multipart.Writer.CreateFormFile takes care of the details)
	var buffer bytes.Buffer
	multipartWriter := multipart.NewWriter(&buffer)
	{
		f, err := formFile.Open()
		if err != nil {
			return "", buffer, fmt.Errorf("error opening '%s' file: %v", fileName, err)
		}
		defer f.Close()

		// Create a Go representation of a form with a file.
		fileWriter, err := multipartWriter.CreateFormFile(fileName, fileName)
		if err != nil {
			return "", buffer, fmt.Errorf("error creating multipart form file '%s': %v",
				fileName, err)
		}

		// Copy the contents of form file received into the new form file.
		_, err = io.Copy(fileWriter, f)
		if err != nil {
			return "", buffer, fmt.Errorf("error copying '%s' file to multipart: %v",
				fileName, err)
		}

		// NOTE: FIELD, not FILE, fileWriter was used above. This is a FIELD
		//       that contains the logged-in user's znumber.
		fieldWriter, err := multipartWriter.CreateFormField("created-by")
		if err != nil {
			return "", buffer, fmt.Errorf("error creating multipart form field 'created-by': %v",
				err)
		}

		// Copy the ZNumber `createdBy` into the newly created form.
		_, err = io.WriteString(fieldWriter, createdBy)
		if err != nil {
			return "", buffer, fmt.Errorf("error copying '%s' file to multipart: %v",
				fileName, err)
		}

	}

	// Closing this thing sets the boundary of the data (required in HTTP file
	// sending protocol) which makes all this work without hard-to-track bugs.
	multipartWriter.Close()

	return multipartWriter.FormDataContentType(), buffer, nil
}

// MultipartReader reads the records from the passed/received CSV-encoded file
// and returns the parsed data (in a 2D String slice format) to the next service.
// It also gives an option to include or exclude file header.
func MultipartReader(formFile *multipart.FileHeader, includeHeader bool) ([][]string, error) {
	fileName := formFile.Filename

	// Open the file
	src, err := formFile.Open()
	if err != nil {
		return nil, fmt.Errorf("error opening '%s' file: %v", fileName, err)
	}
	defer src.Close()

	// Create a Reader that reads records from the CSV-encoded file
	r := csv.NewReader(src)

	// r.FieldsPerRecord = 0
	r.TrimLeadingSpace = true

	if !includeHeader {
		// Skipping first row (line), the headers
		if _, err := r.Read(); err != nil {
			return nil, fmt.Errorf("error parsing the first line in '%s' file: %v", fileName, err)
		}
	}

	// Keep reading everything else
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error parsing the '%s' file: %v", fileName, err)
	}

	return records, nil
}
