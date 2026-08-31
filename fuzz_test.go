package sanswitch

import (
	"bytes"
	"errors"
	"testing"
)

func FuzzParseVersion(f *testing.F) {
	for _, seed := range []string{"v9.2.0a", "9.1", "v8.2.3", "", "v9.x"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, value string) {
		version, err := ParseVersion(value)
		if err != nil {
			if !errors.Is(err, ErrInvalidVersion) {
				t.Fatalf("ParseVersion(%q) error = %v", value, err)
			}
			return
		}
		if !version.Valid() {
			t.Fatalf("ParseVersion(%q) returned an invalid version", value)
		}
		if reparsed, err := ParseVersion(version.String()); err != nil || reparsed != version {
			t.Fatalf("version round trip: got %v, %v", reparsed, err)
		}
	})
}

func FuzzReadResponseBodyLimit(f *testing.F) {
	f.Add([]byte("<Response></Response>"), uint16(1024))
	f.Add([]byte("oversized"), uint16(3))
	f.Fuzz(func(t *testing.T, data []byte, rawLimit uint16) {
		limit := int64(rawLimit) + 1
		body, err := readResponseBody(bytes.NewReader(data), limit)
		if int64(len(data)) > limit {
			if !errors.Is(err, ErrResponseBodyTooLarge) {
				t.Fatalf("oversized body error = %v", err)
			}
			return
		}
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if !bytes.Equal(body, data) {
			t.Fatalf("body changed during bounded read")
		}
	})
}
