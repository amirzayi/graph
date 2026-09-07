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

func (iow iowriter) Log(event string, userID int) error {
	fmt.Fprintf(iow.w, "%s occured by user id: %d\n", event, userID)
	return nil
}
