package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/johannes-kuhfuss/services_utils/date"
)

type HistoryItem struct {
	Date    time.Time
	Message string
}

type HistoryList struct {
	Entries []HistoryItem
}

func (h *HistoryList) Add(date time.Time, msg string) {
	var newEntry HistoryItem
	newEntry.Date = date
	newEntry.Message = msg
	h.Entries = append(h.Entries, newEntry)
}

func (h *HistoryList) AddNow(msg string) {
	now, _ := date.GetNowLocal("")
	h.Add(*now, msg)
}

func (h *HistoryList) ToString() string {
	var history string
	for _, entry := range h.Entries {
		history = history + entry.Date.Format(date.ApiDateLayout) + ": " + entry.Message + "\n"
	}
	return history
}

func (h HistoryList) Value() (driver.Value, error) {
	return json.Marshal(h)
}

func (h *HistoryList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &h)
}
