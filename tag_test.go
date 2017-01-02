package objconv

import "testing"

func TestParseTag(t *testing.T) {
	tests := []struct {
		tag string
		res tag
	}{
		{
			tag: "",
			res: tag{},
		},
		{
			tag: "hello",
			res: tag{name: "hello"},
		},
		{
			tag: ",omitempty",
			res: tag{omitempty: true},
		},
		{
			tag: "-",
			res: tag{name: "-"},
		},
		{
			tag: "hello,omitempty",
			res: tag{name: "hello", omitempty: true},
		},
		{
			tag: "-,omitempty",
			res: tag{name: "-", omitempty: true},
		},
	}

	for _, test := range tests {
		t.Run(test.tag, func(t *testing.T) {
			if res := parseTag(test.tag); res != test.res {
				t.Errorf("%s: %#v != %#v", test.tag, test.res, res)
			}
		})
	}
}
