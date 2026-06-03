package linkutil

import "testing"

func TestParseRange(t *testing.T) {
	start, end, err := ParseRange("[5, 10]")
	if err != nil {
		t.Fatal(err)
	}
	if start != 5 || end != 10 {
		t.Fatalf("got [%d,%d)", start, end)
	}
}

func TestContentRange(t *testing.T) {
	if got := ContentRange(0, 10, 42); got != "links 0-10/42" {
		t.Fatalf("got %q", got)
	}
}
