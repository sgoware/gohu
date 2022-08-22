package mapping

import "encoding/json"

func Struct2Struct(input interface{}, output interface{}) error {
	bytes, _ := json.Marshal(input)
	// TODO: 待修改mapping内部代码
	err := UnmarshalJsonBytes(bytes, output)
	if err != nil {
		return err
	}
	return nil
}
