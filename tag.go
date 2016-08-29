package objconv

import "strings"

// Tag represents the result of parsing the json tag of a struct field.
type Tag struct {
	// Name is the field name that should be used when serializing.
	Name string

	// Omitempty is true if the struct field json tag had `omitempty` set.
	Omitempty bool
}

// ParseTag parses a raw json tag obtained from a struct field,
// returining the results as a Tag value.
func ParseTag(tag string) Tag {
	name, tag := parseNextTagToken(tag)
	token, _ := parseNextTagToken(tag)
	return Tag{Name: name, Omitempty: token == "omitempty"}
}

func parseNextTagToken(tag string) (token string, next string) {
	if split := strings.IndexByte(tag, ','); split < 0 {
		token = tag
	} else {
		token, next = tag[:split], tag[split+1:]
	}
	return
}
