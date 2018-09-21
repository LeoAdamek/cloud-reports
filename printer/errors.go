package printer

import "fmt"

type OutOfBoundsError struct {
	Destination Cursor
}

func oob(dst Cursor) *OutOfBoundsError {
	return  &OutOfBoundsError{Destination: dst}
}

func (e OutOfBoundsError) Error() string {
	return fmt.Sprintf("Requested operation moved print out-of-bounds. %s", e.Destination.Pos())
}