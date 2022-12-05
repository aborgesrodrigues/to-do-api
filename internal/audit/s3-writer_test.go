package audit

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestDirectoryCreation(t *testing.T) {
	now := time.Now()
	actual := getDateDirectory(now)

	if !strings.Contains(actual, strconv.Itoa(now.Year())) {
		t.Errorf("Expected Year: %d in directory: %s", now.Year(), actual)
	}

	if !strings.Contains(actual, strconv.Itoa(int(now.Month()))) {
		t.Errorf("Expected Month: %d in directory: %s", int(now.Month()), actual)
	}

	if !strings.Contains(actual, strconv.Itoa(now.Day())) {
		t.Errorf("Expected Day: %d in directory: %s", now.Day(), actual)
	}

	if !strings.Contains(actual, strconv.Itoa(now.Hour())) {
		t.Errorf("Expected Hour: %d in directory: %s", now.Hour(), actual)
	}
}
