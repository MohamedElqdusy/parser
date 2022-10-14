package formats

import (
	"bufio"
	"encoding/csv"
	"github.com/jszwec/csvutil"
	"io"
	"os"
	"parser/logger"
	"parser/resources"
)

// Reader represents the parsing reader interface.
type Reader interface {

	// Read map the raw data tot a specific resource.
	Read(resources.Resource, *Chans) (Summary, error)
}

// Csv read the csv data.
type Csv struct {
	reader      io.Reader
	readingFile *os.File
	summery     Summary
}

// NewCsv returns a new csv reader
func NewCsv(filePath string) (Reader, error) {
	var (
		f   *os.File
		err error
	)
	if f, err = os.Open(filePath); err != nil {
		return nil, err
	}
	return &Csv{reader: bufio.NewReader(f), summery: Summary{}}, nil

}

// Read csv read implementation.
func (c *Csv) Read(res resources.Resource, channels *Chans) (Summary, error) {
	csvReader := csv.NewReader(c.reader)
	dec, err := csvutil.NewDecoder(csvReader)
	if err != nil {
		return Summary{}, err
	}

	//header := dec.Header()
	for {
		// update the total count summery
		c.summery.Total += 1
		// use a new zero value every iteration
		res = res.NewZeroValue()
		if err := dec.Decode(res); err == io.EOF {
			break
		} else if err != nil {
			// send decoding errors to the summary report then break
			logger.Error("found an error while decoding: %w", err)
			channels.err <- err
			// update the skipped count summery
			c.summery.Skipped += 1
			break
		}

		// Validates the csv row after mapping to a resource
		if err = res.Validate(); err != nil {
			// send validation's error to the summery report
			channels.err <- err
			// update the skipped count summery
			c.summery.Skipped += 1
		} else {
			// send the valid parsed record to the result file.
			channels.resource <- res
			// update valid summery
			c.summery.Valid += 1
		}
	}

	return c.summery, nil
}
