package text

import "fmt"

// A Title is a string which cannot exceed 30 characters in length.
type Title string

func (t *Title) UnmarshalText(p []byte) (err error) {
	*t, err = ParseTitle(string(p))
	return
}

func (t *Title) Scan(v any) (err error) {
	switch v := v.(type) {
	case string:
		*t, err = ParseTitle(v)
	default:
		return fmt.Errorf("converting type %T to a text.Title", v)
	}
	return err
}

func ParseTitle(s string) (Title, error) {
	if n := len(s); n < 1 || n > 30 {
		return "", ErrInvalidLength
	}
	return Title(s), nil
}
