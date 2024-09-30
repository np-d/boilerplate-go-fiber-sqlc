package converter

import "encoding/json"

func StructToMap(src any) (*map[string]any, error) {
	out := make(map[string]any)
	data, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &out)
	if err != nil {
		return nil, err
	}
	return &out, err
}
