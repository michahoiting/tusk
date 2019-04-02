package marshal

// ListMap is a map from string to string optionally represented in yaml as a
// list of strings. A list of strings will be unmarshalled as if both the key
// and value are the same.
type ListMap map[string]string

// UnmarshalYAML unmarshals a list or map always into a map.
func (lm *ListMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var list []string
	listCandidate := UnmarshalCandidate{
		Unmarshal: func() error { return unmarshal(&list) },
		Assign: func() {
			listMap := make(map[string]string, len(list))
			for _, name := range list {
				listMap[name] = name
			}
			*lm = listMap
		},
	}

	var listMap map[string]string
	mapCandidate := UnmarshalCandidate{
		Unmarshal: func() error { return unmarshal(&listMap) },
		Assign:    func() { *lm = listMap },
	}

	return UnmarshalOneOf(listCandidate, mapCandidate)
}
