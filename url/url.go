package url

import (
	"encoding/json"
	"net/url"
	"strconv"
)

func quote(s string) string {
	if strconv.CanBackquote(s) {
		return "`" + s + "`"
	}
	return strconv.Quote(s)
}

type URL struct {
	url.URL
}

func Parse(s string) (URL, error) {

	result, err := url.Parse(s)
	if err != nil {
		return URL{}, err
	}

	return URL{*result}, nil
}

func MustParse(s string) URL {

	result, err := Parse(s)
	if err != nil {
		panic(`url: Parse(` + quote(s) + `): ` + err.Error())
	}

	return result
}

func (u *URL) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	return u.URL.UnmarshalBinary([]byte(s))
}

func (u URL) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}
