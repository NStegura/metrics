package custom_errors

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

type ParseUrlError struct {
	Url string
}

func (e ParseUrlError) Error() string {
	return fmt.Sprintf("bad request url: %s", e.Url)
}
