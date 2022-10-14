package resources

import (
	"bytes"
	"fmt"
	"strings"
)

type Resource interface {
	NewZeroValue() Resource
	Validate() error
}

type ValidateError struct {
	Field  string
	Tag    string
	Value  interface{}
	Record interface{}
}

func (v ValidateError) Error() string {
	return fmt.Sprintf("Error in Record [%v] :Field validation for '%s' failed on the '%s' tag  with value: %+v", v.Record, v.Field, v.Tag, v.Value)
}

type ValidateErrors []ValidateError

func (vs ValidateErrors) Error() string {
	buff := bytes.NewBufferString("")

	var e ValidateError
	for i := 0; i < len(vs); i++ {
		e = vs[i]
		buff.WriteString(e.Error())
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}
