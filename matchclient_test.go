package matchclient

import (
	"regexp"
	"testing"
)

// TestHelloName calls matchclient.Hello with a name, checking
// for a valid return value.
func TestHello(t *testing.T) {
	name := "Viki"
	want := regexp.MustCompile(`\b` + name + `\b`)
	msg, err := Hello(name)
	if !want.MatchString(msg) || err != nil {
		t.Fatalf(`Hello("Viki") = %q, %v, want match for %#q, nil`, msg, err, want)
	}
}

// TestHelloEmpty calls matchclient.Hello with an empty string,
// checking for an error.
func TestHelloEmpty(t *testing.T) {
	msg, err := Hello("")
	if msg != "" || err == nil {
		t.Fatalf(`Hello("") = %q, %v, want "", error`, msg, err)
	}
}
