package process

import (
	"fmt"
	"parser/formats"
	"parser/logger"
	"parser/resources"
	"reflect"
)

// Processor is the parsing engine
type Processor struct {
	Input   formats.Reader  // represents the raw data reader
	Output  formats.Writer  // represents the parsed data writer
	Summary formats.Summary // the final parsing stats
}

// NewDownloadProcessor creates a processor for an url
func NewDownloadProcessor(url string) (Processor, error) {
	var (
		p        Processor
		err      error
		filePath string
	)
	// download the file
	if filePath, err = downloadFile(url); err != nil {
		return Processor{}, err
	}
	// create a local file processor
	if p, err = NewFileProcessor(filePath); err != nil {
		return Processor{}, err
	}
	return p, nil
}

// NewFileProcessor creates a processor for a local file
func NewFileProcessor(filePath string) (Processor, error) {
	var (
		csv  formats.Reader
		json formats.Writer
		err  error
	)
	// create a csv reader
	if csv, err = formats.NewCsv(filePath); err != nil {
		return Processor{}, fmt.Errorf("couldn't start CSV reader: %w", err)
	}
	logger.Info("created csv reader's processor for file ", filePath)

	// creating Json result writer
	if json, err = formats.NewJson(filePath); err != nil {
		return Processor{}, fmt.Errorf("couldn't start Json writer: %w", err)
	}
	logger.Info("created json writer's processor for file ", filePath)
	return Processor{
		Input:  csv,
		Output: json,
	}, nil
}

// Process starts the parsing process.
func (p *Processor) Process(resource resources.Resource) error {
	var (
		s   formats.Summary
		err error
	)
	logger.Info("Processor started processing with resource type ", reflect.TypeOf(resource))
	// start writing
	channels := formats.NewChannels()
	go p.Output.Write(&channels)

	// start reading
	if s, err = p.Input.Read(resource, &channels); err != nil {
		logger.Error("Error Reading the resources from CSV: %w", err)
	}
	// save the summary stats to the processor
	p.Summary = s
	// write the final processing stats tot the summary report
	// it should be at the last line.
	channels.SendSummary(&p.Summary)
	// close the channels
	channels.Close()
	logger.Info(fmt.Sprintf("%+v", p.Summary))
	return nil
}
