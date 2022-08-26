package uuid

import (
	"github.com/speps/go-hashids/v2"
	"github.com/spf13/cast"
)

func NewRandomString(src string, salt string, length int) string {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = length
	h, _ := hashids.NewWithData(hd)
	randomString, _ := h.Encode(cast.ToIntSlice([]uint8(src)))
	return randomString
}
