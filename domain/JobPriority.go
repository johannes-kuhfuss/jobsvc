package domain

import (
	"github.com/johannes-kuhfuss/services_utils/enums"
)

const (
	DefaultJobPriority int32 = 30
)

var (
	JobPriority = enums.Enum{
		Items: []enums.EnumItem{{Idx: 50, Val: "realtime"}, {Idx: 40, Val: "high"}, {Idx: 30, Val: "medium"}, {Idx: 20, Val: "low"}, {Idx: 10, Val: "idle"}},
	}
)

func IsValidPriority(prio string) bool {
	_, err := JobPriority.ItemByValue(prio)
	return err == nil
}
