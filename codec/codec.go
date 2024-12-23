package codec

import "io"

type Header struct {
	ServiceMethod string
	Seq           uint64
	Error         string
}

// Codec 消息体编码和解码的接口
type Codec interface {
	io.Closer
	ReadHeader(header *Header) error              // 读取头部信息
	ReadBody(body interface{}) error              // 读取body信息
	Write(header *Header, body interface{}) error // 写信息
}

// NewCodecFunc 构造函数
type NewCodecFunc func(closer io.ReadWriteCloser) Codec

type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json"
)

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	// value是构造函数而非实例
	NewCodecFuncMap[GobType] = NewGobCodec
}
