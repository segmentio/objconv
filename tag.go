package objconv

import "strings"

// Tag represents the result of parsing the tag of a struct field.
type Tag struct {
	// Name is the field name that should be used when serializing.
	Name string

	// Omitempty is true if the tag had `omitempty` set.
	Omitempty bool

	// Omitzero is true if the tag had `omitzero` set.
	Omitzero bool
}

// ParseTag parses a raw tag obtained from a struct field, returning the results
// as a Tag value.
func ParseTag(tag string) Tag {
	var tokens [2]string
	var omitzero bool
	var omitempty bool

	name, tag := parseNextTagToken(tag)
	tokens[0], tag = parseNextTagToken(tag)
	tokens[1], tag = parseNextTagToken(tag)

	for _, t := range tokens {
		switch t {
		case "omitempty":
			omitempty = true
		case "omitzero":
			omitzero = true
		}
	}

	return Tag{
		Name:      name,
		Omitempty: omitempty,
		Omitzero:  omitzero,
	}
}

func parseNextTagToken(tag string) (token string, next string) {
	if split := strings.IndexByte(tag, ','); split < 0 {
		token = tag
	} else {
		token, next = tag[:split], tag[split+1:]
	}
	return
}
