package objconv

import "fmt"

// Assert checks if the boolean gived as first argument is true, if not it
// panics with an error message formatted off of the remaining arguments.
func Assert(pass bool, args ...interface{}) {
	if !pass {
		panic(fmt.Sprint(args...))
	}
}

// Assertf checks if the boolean given as first argument is true, it not it
// panics with an error message formatted off of the format string and remaining
// arguments.
func Assertf(pass bool, msg string, args ...interface{}) {
	if !pass {
		panic(fmt.Sprintf(msg, args...))
	}
}

// AssertErr checks if the error given as argument is nil, if not it panics with
// the error itself as value.
func AssertErr(err error) {
	if err != nil {
		panic(err)
	}
}
