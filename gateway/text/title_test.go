package text_test

import (
	"encoding/json"
	"strings"
	"testing"

	"go.adoublef.dev/is"
	. "go.petal-hub.io/gateway/text"
)

func Test_ParseTitle(t *testing.T) {
	type testcase struct {
		s   string
		err error
	}
	tt := map[string]testcase{
		"Simple": {
			s: "My Title",
		},
		"Empty": {
			s:   "",
			err: ErrInvalidLength,
		},
		"TooLong": {
			s:   strings.Repeat("Abcdefghijklmnopqrstuvwxyz", 5),
			err: ErrInvalidLength,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			is := is.NewRelaxed(t)

			_, err := ParseTitle(tc.s)
			is.Err(err, tc.err)
		})
	}
}

func Test_Title_UnmarshalJSON(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var in = `"My title is valid"`
		var title Title
		err := json.NewDecoder(strings.NewReader(in)).Decode(&title)
		is.NoErr(err)
	})

	t.Run("Err", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var in = `""`
		var title Title
		err := json.NewDecoder(strings.NewReader(in)).Decode(&title)
		is.Err(err, ErrInvalidLength)
	})
}

// TODO test for db scanning
