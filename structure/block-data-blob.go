package structure

type BlockDataBlob struct {
	palette map[string]int32
	data    []byte
	xLen    int32
	yLen    int32
	zLen    int32
}

func NewBlockDataBlob(palette map[string]int32, data []byte, xLen, yLen, zLen int32) *BlockDataBlob {
	return &BlockDataBlob{palette: palette, data: data, xLen: xLen, yLen: yLen, zLen: zLen}
}

func (d *BlockDataBlob) ReadPalette(i byte) string {
	for j := range d.palette {
		if d.palette[j] == int32(i) {
			return j
		}
	}
	return ""
}

func (d *BlockDataBlob) WritePalette(v string, i byte) {
	d.palette[v] = int32(i)
}

func (d *BlockDataBlob) ReadBlock(x, y, z int32) byte {
	return d.data[y*d.xLen*d.zLen+z*d.xLen+x]
}

func (d *BlockDataBlob) WriteBlock(x, y, z int32, v byte) {
	d.data[y*d.xLen*d.zLen+z*d.xLen+x] = v
}

func (d *BlockDataBlob) SingleBlockData(id string) *BlockDataBlob {
	c := byte(d.palette[id])
	b := make([]byte, len(d.data))
	for i := range d.data {
		if d.data[i] == c {
			b[i] = 1
		} else {
			b[i] = 0
		}
	}
	return NewBlockDataBlob(map[string]int32{"minecraft:air": 0, id: 1}, b, d.xLen, d.yLen, d.zLen)
}

func (d *BlockDataBlob) Output() (map[string]int32, []byte, int32, int32, int32) {
	return d.palette, d.data, d.xLen, d.yLen, d.zLen
}
