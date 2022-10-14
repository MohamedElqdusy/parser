package formats

import (
	"parser/resources"
)

// Chans used to communicat between reader and writer.
type Chans struct {
	resource chan resources.Resource
	err      chan error
	summary  chan *Summary
	stop     chan bool
}

// Summary represents the final parsing summary
type Summary struct {
	Total   int64 // total processed records
	Valid   int64 // valid record count
	Skipped int64 // skipped record count that have not been written to the result
}

// NewChannels returns a new channels
func NewChannels() Chans {
	return Chans{resource: make(chan resources.Resource),
		summary: make(chan *Summary),
		err:     make(chan error),
		stop:    make(chan bool)}
}

// SendSummary used to send the last write for summary stats.
func (c *Chans) SendSummary(s *Summary) {
	c.summary <- s
}

// Close the channels
func (c *Chans) Close() {
	close(c.resource)
	close(c.err)
	close(c.summary)
	c.stop <- true
}
