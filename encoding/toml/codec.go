package toml

import (
	"bytes"

	"github.com/BurntSushi/toml"
	"github.com/jacexh/gopkg/encoding"
)

type tomlCodec struct{}

func (c tomlCodec) Name() string {
	return "toml"
}

func (c tomlCodec) Marshal(v interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := toml.NewEncoder(buf).Encode(v)
	if err == nil {
		return buf.Bytes(), nil
	}
	return nil, err
}

func (c tomlCodec) Unmarshal(data []byte, v interface{}) error {
	return toml.Unmarshal(data, v)
}

func init() {
	encoding.RegisterCodec(tomlCodec{})
}
