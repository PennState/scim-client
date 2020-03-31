package scim

import "fmt"

type CodecOperation string

const (
	Marshal   CodecOperation = "Marshal"
	Unmarshal CodecOperation = "Unmarshal"
)

type CodecError struct {
	Err  string
	Op   CodecOperation
	Body []byte
}

func (ce CodecError) Error() string {
	return fmt.Sprintf("Err: %s, Operation: %s, Body: %s", ce.Err, ce.Op, string(ce.Body))
}
