package audit

import (
	"fmt"
	"io"
)

type iowriter struct {
	w io.Writer
}

func NewIOWriter(w io.Writer) iowriter {
	return iowriter{w: w}
}

func (iow iowriter) Log(event string, userID int, traceID string) error {
	_, err := fmt.Fprintf(iow.w, "%s occured by user id: %d on trace id: %s\n", event, userID, traceID)
	return err
}
