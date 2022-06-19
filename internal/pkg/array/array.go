package array

func StringsToMap(v []string) map[string]struct{} {
	result := map[string]struct{}{}
	for _, s := range v {
		result[s] = struct{}{}
	}
	return result
}
