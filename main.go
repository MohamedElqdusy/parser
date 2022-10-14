package main

import (
	"os"
	"parser/logger"
	"parser/process"
	"parser/resources"
)

func main() {
	log := logger.NewLogger()
	logger.InitLogger(log)
	defer logger.Sync()

	var (
		p   process.Processor
		err error
	)

	// pass the url from the cli
	if len(os.Args) != 2 {
		panic("You should pass two CLI arguments!")
	}

	// creates a new processor for the file at the url
	if p, err = process.NewDownloadProcessor(os.Args[1]); err != nil {
		logger.Error(err)
		return
	}
	// starts parsing process
	if err := p.Process(&resources.Customer{}); err != nil {
		logger.Error(err)
	}
	logger.Info("Done!")
}
