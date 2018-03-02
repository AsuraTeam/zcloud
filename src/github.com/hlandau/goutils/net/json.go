package net

import "encoding/base64"
import "encoding/json"
import "strings"

type Base64 []byte

func (b *Base64) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	*b, err = base64.StdEncoding.DecodeString(s)
	return err
}

func (b Base64) MarshalJSON() ([]byte, error) {
	return json.Marshal(base64.StdEncoding.EncodeToString(b))
}

type Base64up []byte

func (b *Base64up) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	*b, err = base64.URLEncoding.DecodeString(RepadBase64(s))
	return err
}

func (b Base64up) MarshalJSON() ([]byte, error) {
	return json.Marshal(strings.TrimRight(base64.URLEncoding.EncodeToString(b), "="))
}

func RepadBase64(s string) string {
	if m := len(s) % 4; m >= 2 {
		return s + "=="[m-2:]
	}

	return s
}

// untested
func RepadBase32(s string) string {
	if m := len(s) % 8; m >= 6 {
		return s + "======"[m-6:]
	}

	return s
}
