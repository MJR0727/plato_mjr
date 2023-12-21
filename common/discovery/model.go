package discovery

import "encoding/json"

type EndpointInfo struct {
	IP       string                 `json:"ip"`
	Port     string                 `json:"port"`
	MetaData map[string]interface{} `json:"meta"`
}

func (epi *EndpointInfo) Marshal() string {
	bytes, err := json.Marshal(epi)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func UnMarshal(data []byte) (*EndpointInfo, error) {
	epi := &EndpointInfo{}
	err := json.Unmarshal(data, epi)
	if err != nil {
		return nil, err
	}
	return epi, nil
}
