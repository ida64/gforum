package utils

import (
	"fmt"
	"time"
)

func TimeAgo(inputTime time.Time) string {
	duration := time.Since(inputTime)
	if duration.Minutes() < 1 {
		return fmt.Sprintf("%.0f seconds ago", duration.Seconds())
	} else if duration.Hours() < 1 {
		return fmt.Sprintf("%.0f minutes ago", duration.Minutes())
	} else if duration.Hours() < 24 {
		return fmt.Sprintf("%.0f hours ago", duration.Hours())
	}

	return inputTime.Format("January 2, 2006")
}
