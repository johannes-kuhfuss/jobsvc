package utils

func SliceContainsString(slice []string, searchStr string) bool {
	for _, val := range slice {
		if val == searchStr {
			return true
		}
	}
	return false
}
