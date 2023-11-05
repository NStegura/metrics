package customerrors

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

type ParseURLError struct {
	URL string
}

func (e ParseURLError) Error() string {
	return fmt.Sprintf("bad request url: %s", e.URL)
}
