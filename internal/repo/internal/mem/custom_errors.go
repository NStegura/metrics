package mem

import "fmt"

type BackupError struct {
	prevErr error
}

func (b BackupError) Error() string {
	return fmt.Sprintf("Make backup err: %s", b.prevErr)
}
