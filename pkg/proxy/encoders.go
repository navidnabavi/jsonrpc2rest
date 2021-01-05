package proxy

import "encoding/json"

//EncodePayload Encodes Palyload to json
func EncodePayload(params map[string]interface{}, contentType uint8) ([]byte, error) {
	return json.Marshal(params)
}
