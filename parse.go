package objconv

// The Parser interface must be implemented by types that provide decoding of a
// specific format (like json, resp, ...).
type Parser interface {
	/*
			ParseBegin(*Reader) error

			ParseEnd(*Reader) error

			ParseType(*Reader) (Type, error)

			ParseBool(*Reader) (bool, error)

			ParseInt(*Reader) (int64, error)

			ParseUint(*Reader) (uint64, error)

			ParseFloat(*Reader) (float64, error)

			ParseString(*Reader, *bytes.Buffer) error

			ParseBytes(*Reader, *bytes.Buffer) error

			ParseTime(*Reader) (time.Time, error)

			ParseDuration(*Reader) (time.Duration, error)

			ParseArrayBegin(*Reader) (int, error)

			ParseArrayEnd(*Reader) error

			ParseArrayNext(*Reader) error

			ParseMapBegin(*Reader) (int, error)

			ParseMapEnd(*Reader) error

			ParseMapValue(*Reader) error

			ParseMapNext(*Reader) error

		Parse(*Reader, interface{}) (interface{}, error)
	*/
}
