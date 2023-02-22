package marshall

import (
	"compress/gzip"
	"github.com/Tnze/go-mc/nbt"
	"io"
)

func decodeGeneric[T any](r io.Reader) (name string, data T, err error) {
	var r2 io.Reader
	r2, err = gzip.NewReader(r)
	if err != nil {
		return
	}
	decoder := nbt.NewDecoder(r2)
	name, err = decoder.Decode(&data)
	if err != nil {
		return
	}
	return
}

func encodeGeneric[T any](w io.Writer, name string, data T) error {
	w2 := gzip.NewWriter(w)
	encoder := nbt.NewEncoder(w2)
	err := encoder.Encode(data, name)
	if err != nil {
		return err
	}
	return w2.Close()
}
