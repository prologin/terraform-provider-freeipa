package api

import (
	"strings"
	"time"

	b64 "encoding/base64"
)

type IPATime struct {
	time.Time
}

const ipaTimeLayout = "20060102150405Z"

func (ipat *IPATime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		return nil
	}

	ipat.Time, err = time.Parse(ipaTimeLayout, s)
	if err != nil {
		return err
	}

	return nil
}

func (ipat *IPATime) MarshalJSON() ([]byte, error) {
	if ipat == nil {
		return []byte("null"), nil
	}

	return []byte(ipat.Time.Format(ipaTimeLayout)), nil
}

func (ipat *IPATime) String() string {
	return ipat.Time.Format(time.RFC3339)
}

type Base64EncodedSecret struct {
	string
}

func (s *Base64EncodedSecret) UnmarshalJSON(b []byte) (err error) {
	return nil
}

func (s *Base64EncodedSecret) MarshalJSON() ([]byte, error) {
	return []byte(s.string), nil
}

func (s *Base64EncodedSecret) Decode() string {
	secret, err := b64.StdEncoding.DecodeString(s.string)

	if err != nil {
		return ""
	}

	return string(secret)
}
