package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type GobCodec struct {
	conn io.ReadWriteCloser // 通过tcp或者Unix建立socket时得到的连接
	buf  *bufio.Writer      // 缓冲区
	dec  *gob.Decoder       // 解码
	enc  *gob.Encoder       // 编码
}

// 强制检查是否GobCodec实现了Codec接口
var _ Codec = (*GobCodec)(nil)

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}

// ReadHeader 就是解析header,一般是由指定编码过来的数据
func (g GobCodec) ReadHeader(header *Header) error {
	return g.dec.Decode(header)
}

func (g GobCodec) ReadBody(body interface{}) error {
	return g.dec.Decode(body)
}

func (g GobCodec) Write(header *Header, body interface{}) (err error) {
	// 写完之后要清空buf
	defer func() {
		_ = g.buf.Flush()
		if err != nil {
			_ = g.Close()
		}
	}()
	if err := g.enc.Encode(header); err != nil {
		log.Println("rpc error: gob error encoding header: ", err)
		return err
	}
	if err := g.enc.Encode(body); err != nil {
		log.Println("rpc error: gob error encoding body: ", err)
		return err
	}

	return nil
}
func (g GobCodec) Close() error {
	// 关闭连接
	return g.conn.Close()
}
