package uuid

import (
	"strconv"

	"github.com/google/uuid"
)

type UUID struct {
	uuid.UUID
}

func Parse(s string) (UUID, error) {

	result, err := uuid.Parse(s)
	if err != nil {
		return UUID{}, err
	}

	return UUID{result}, nil
}

func MustParse(s string) UUID {

	result, err := Parse(s)
	if err != nil {
		panic(`uuid: Parse(` + quote(s) + `): ` + err.Error())
	}

	return result
}

func (u *UUID) UnmarshalJSON(b []byte) error {
	return u.UUID.UnmarshalText(b)
}

func (u UUID) MarshalJSON() ([]byte, error) {
	return u.UUID.MarshalText()
}

func quote(s string) string {
	if strconv.CanBackquote(s) {
		return "`" + s + "`"
	}
	return strconv.Quote(s)
}
