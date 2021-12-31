package domain

import (
	"github.com/johannes-kuhfuss/services_utils/enums"
)

var (
	JobPriority = enums.Enum{
		Items: []enums.EnumItem{{Idx: 0, Val: "realtime"}, {Idx: 1, Val: "high"}, {Idx: 2, Val: "medium"}, {Idx: 3, Val: "low"}, {Idx: 4, Val: "idle"}},
	}
	DefaultJobPriority int32 = 2
)

func IsValidPriority(prio string) bool {
	_, err := JobPriority.ItemByValue(prio)
	return err == nil
}
