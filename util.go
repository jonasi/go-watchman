package watchman

func stringParams(strs ...string) []interface{} {
	p := make([]interface{}, len(strs))

	for i := range strs {
		p[i] = strs[i]
	}

	return p
}
