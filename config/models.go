package config

import (
	"encoding/json"
	"fmt"
	"time"
)

type Duration time.Duration

func (d *Duration) UnmarshalJSON(b []byte) error {
	var durationStr string
	if err := json.Unmarshal(b, &durationStr); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return fmt.Errorf("failed to parse duration: %w", err)
	}

	*d = Duration(duration)
	return nil
}
