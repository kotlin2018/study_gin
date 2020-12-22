// +build jsoniter

package json

import  "github.com/json-iterator/go" //是一款快且灵活的 JSON 解析器,它最多能比普通的解析器快 10 倍之多

var (
	json 			= jsoniter.ConfigCompatibleWithStandardLibrary
	Marshal 		= json.Marshal
	Unmarshal 		= json.Unmarshal
	MarshalIndent 	= json.MarshalIndent
	NewDecoder 		= json.NewDecoder
	NewEncoder 		= json.NewEncoder
)
