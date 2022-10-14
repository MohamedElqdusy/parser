package formats

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"parser/logger"
)

// Writer represents the parser writer interface.
type Writer interface {
	// Write writes the parsed records, summary or errors to their final destination.
	Write(*Chans)
}

// Json is a writer implementation for writing using json format.
type Json struct {
	summaryWriter io.Writer // used to write to the summary file.
	resultWriter  io.Writer // used to write to the result file.
}

// NewJson returns a new json writer.
func NewJson(fileName string) (Writer, error) {
	var (
		err    error
		sf, rf *os.File
	)
	// create the summary file.
	if sf, err = os.Create(fmt.Sprintf("%s_summary", fileName)); err != nil {
		return nil, err
	}
	// create the parsed record file.
	if rf, err = os.Create(fmt.Sprintf("%s_result", fileName)); err != nil {
		return nil, err
	}
	return &Json{
		summaryWriter: sf,
		resultWriter:  rf,
	}, nil
}

// Write writes the parsed records in json format
func (j *Json) Write(chans *Chans) {
	summaryEncoder := json.NewEncoder(j.summaryWriter)
	resultEncoder := json.NewEncoder(j.resultWriter)

	for {
		select {
		// receiving a valid parsed record.
		case j := <-chans.resource:
			if j != nil {
				// write the resource's valid parsing to the result file
				if err := resultEncoder.Encode(j); err != nil {
					logger.Error(err)
				}
			}
		// receiving a parsing error
		case m := <-chans.err:
			if m != nil {
				// write the resource's parsing errors to the summary file
				if err := summaryEncoder.Encode(m); err != nil {
					logger.Error(err)
				}
			}
		// receiving the last summary stats
		case s := <-chans.summary:
			if s != nil {
				// write the whole parsing summary stats to the summary file
				// it should be at the last line in the file.
				if err := summaryEncoder.Encode(s); err != nil {
					logger.Error(err)
				}
			}
		// exit
		case <-chans.stop:
			return
		}
	}
}
