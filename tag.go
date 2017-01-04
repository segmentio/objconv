package objconv

import "strings"

// tag represents the result of parsing the tag of a struct field.
type tag struct {
	// Name is the field name that should be used when serializing.
	name string

	// Omitempty is true if the tag had `omitempty` set.
	omitempty bool

	// Omitzero is true if the tag had `omitzero` set.
	omitzero bool
}

// ParseTag parses a struct field tag in s, returning the name of the field and
// its properties.
func ParseTag(s string) (name string, omitempty bool, omitzero bool) {
	t := parseTag(s)
	return t.name, t.omitempty, t.omitzero
}

// parseTag parses a raw tag obtained from a struct field, returning the results
// as a tag value.
func parseTag(s string) tag {
	var tokens [2]string
	var omitzero bool
	var omitempty bool

	name, s := parseNextTagToken(s)
	tokens[0], s = parseNextTagToken(s)
	tokens[1], _ = parseNextTagToken(s)

	for _, t := range tokens {
		switch t {
		case "omitempty":
			omitempty = true
		case "omitzero":
			omitzero = true
		}
	}

	return tag{
		name:      name,
		omitempty: omitempty,
		omitzero:  omitzero,
	}
}

func parseNextTagToken(s string) (token string, next string) {
	if split := strings.IndexByte(s, ','); split < 0 {
		token = s
	} else {
		token, next = s[:split], s[split+1:]
	}
	return
}
