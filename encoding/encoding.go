package encoding

import (
	"github.com/yaoapp/gou/encoding/base64"
	"github.com/yaoapp/gou/encoding/hex"
	"github.com/yaoapp/gou/encoding/json"
	"github.com/yaoapp/gou/process"
)

func init() {
	process.Register("encoding.base64.Encode", base64.ProcessEncode)
	process.Register("encoding.base64.Decode", base64.ProcessDecode)
	process.Register("encoding.hex.Encode", hex.ProcessEncode)
	process.Register("encoding.hex.Decode", hex.ProcessDecode)
	process.Register("encoding.json.Encode", json.ProcessEncode)
	process.Register("encoding.json.Decode", json.ProcessDecode)
}
