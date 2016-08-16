package objconv

import (
	"reflect"
	"strings"
	"unicode"
)

// StructTag represents the result of parsing a struct field tag that was
// laid out with the standard key:"value" format.
type StructTag map[string]Tag

// Tag represents the result of parsing the json tag of a struct field.
type Tag struct {
	// Name is the field name that should be used when serializing.
	Name string

	// Omitempty is true if the struct field json tag had `omitempty` set.
	Omitempty bool

	// Skip is true if the struct field json tag started with `-`.
	Skip bool
}

// ParseStructTag parses the tag of a struct field that may or may not
// returing the result as a StructTag.
func ParseStructTag(t reflect.StructTag) StructTag {
	tags := make(StructTag, 1)

	for tag := strings.TrimLeftFunc(string(t), unicode.IsSpace); len(tag) != 0; tag = strings.TrimLeftFunc(tag, unicode.IsSpace) {
		var key string
		var val string

		key, val, tag = parseNextTagPair(tag)
		tags[key] = ParseTag(val)
	}

	return tags
}

// ParseTag parses a raw json tag obtained from a struct field,
// returining the results as a Tag value.
func ParseTag(tag string) Tag {
	name, tag := parseNextTagToken(tag)
	token, _ := parseNextTagToken(tag)
	return Tag{
		Name:      name,
		Skip:      name == "-",
		Omitempty: token == "omitempty",
	}
}

func parseNextTagPair(tag string) (key string, val string, next string) {
	key, next = parseNextTagKey(tag)
	val, next = parseNextTagVal(next)
	return
}

func parseNextTagKey(tag string) (key string, next string) {
	if split := strings.IndexByte(tag, ':'); split < 0 {
		key = tag
	} else {
		key, next = strings.TrimSpace(tag[:split]), tag[split+1:]
	}
	return
}

func parseNextTagVal(tag string) (val string, next string) {
	if len(tag) == 0 || tag[0] != '"' {
		next = tag
	} else if split := strings.IndexByte(tag[1:], '"'); split < 0 {
		val = tag[1:]
	} else {
		val, next = tag[1:split+1], tag[split+2:]
	}
	return
}

func parseNextTagToken(tag string) (token string, next string) {
	if split := strings.IndexByte(tag, ','); split < 0 {
		token = tag
	} else {
		token, next = tag[:split], tag[split+1:]
	}
	return
}
