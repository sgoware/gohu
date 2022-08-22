package mapping

import (
	"github.com/zeromicro/go-zero/core/jsonx"
)

const jsonTagKey = "json"

var jsonUnmarshaler = NewUnmarshaler(jsonTagKey)

// UnmarshalJsonBytes unmarshals content into v.
func UnmarshalJsonBytes(content []byte, v interface{}) error {
	return unmarshalJsonBytes(content, v, jsonUnmarshaler)
}

func unmarshalJsonBytes(content []byte, v interface{}, unmarshaler *Unmarshaler) error {
	var m map[string]interface{}
	if err := jsonx.Unmarshal(content, &m); err != nil {
		return err
	}

	return unmarshaler.Unmarshal(m, v)
}
