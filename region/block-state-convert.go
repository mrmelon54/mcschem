package region

import (
	"bytes"
	"fmt"
	"github.com/Tnze/go-mc/level/block"
	"github.com/Tnze/go-mc/nbt"
	"github.com/Tnze/go-mc/save"
)

func blockStateToId(v save.BlockState) (s block.StateID, err error) {
	b, ok := block.FromID[v.Name]
	if !ok {
		err = fmt.Errorf("unknown block id: %v", v.Name)
		return
	}
	if v.Properties.Data != nil {
		if err2 := v.Properties.Unmarshal(&b); err != nil {
			err = fmt.Errorf("unmarshal block properties fail: %v", err2)
			return
		}
	}
	s, ok = block.ToStateID[b]
	if !ok {
		err = fmt.Errorf("unknown block: %v", b)
		return
	}
	return s, nil
}

func idToBlockState(s block.StateID) (a save.BlockState) {
	b := block.StateList[s]
	a.Name = b.ID()

	buf := new(bytes.Buffer)
	buf.Reset()
	err := nbt.NewEncoder(buf).Encode(b, "")
	if err != nil {
		return
	}
	_, err = nbt.NewDecoder(buf).Decode(&a.Properties)
	return
}
