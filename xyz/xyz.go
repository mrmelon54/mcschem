package xyz

type XYZ struct {
	X int32 `nbt:"x"`
	Y int32 `nbt:"y"`
	Z int32 `nbt:"z"`
}

func (xyz XYZ) Add(other XYZ) XYZ {
	return XYZ{
		X: xyz.X + other.X,
		Y: xyz.Y + other.Y,
		Z: xyz.Z + other.Z,
	}
}

func (xyz XYZ) Sub(other XYZ) XYZ {
	return XYZ{
		X: xyz.X - other.X,
		Y: xyz.Y - other.Y,
		Z: xyz.Z - other.Z,
	}
}

func (xyz XYZ) Abs() XYZ {
	a := xyz
	if a.X < 0 {
		a.X *= -1
	}
	if a.Y < 0 {
		a.Y *= -1
	}
	if a.Z < 0 {
		a.Z *= -1
	}
	return a
}

func (xyz XYZ) Count() int32 {
	a := xyz.X * xyz.Y * xyz.Z
	if a < 0 {
		a *= -1
	}
	return a
}

func (xyz XYZ) LoopAllBlocks(cb func(XYZ)) {
	var x, y, z int32
	for x = 0; x < xyz.X; x++ {
		for y = 0; y < xyz.Y; y++ {
			for z = 0; z < xyz.Z; z++ {
				cb(XYZ{X: x, Y: y, Z: z})
			}
		}
	}
}
