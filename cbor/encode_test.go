package cbor

import "testing"

type counter struct {
	n int
}

func (c *counter) Write(b []byte) (int, error) {
	c.n += len(b)
	return len(b), nil
}

func BenchmarkMarshal(b *testing.B) {
	for _, test := range cborTests {
		b.Run(testName(test), func(b *testing.B) {
			n := 0

			for i := 0; i != b.N; i++ {
				s, err := Marshal(test)
				if err != nil {
					b.Fatal(err)
				}
				n = len(s)
			}

			b.SetBytes(int64(n))
		})
	}
}

func BenchmarkEncoder(b *testing.B) {
	c := &counter{}
	e := NewEncoder(c)

	for _, test := range cborTests {
		b.Run(testName(test), func(b *testing.B) {
			for i := 0; i != b.N; i++ {
				c.n = 0
				if err := e.Encode(test); err != nil {
					b.Fatal(err)
				}
			}
			b.SetBytes(int64(c.n))
		})
	}
}
