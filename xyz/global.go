package xyz

import sf "github.com/sa-/slicefunk"

func minIntOf(vars ...int32) int32 {
	min := vars[0]
	for _, i := range vars {
		if min > i {
			min = i
		}
	}
	return min
}

func maxIntOf(vars ...int32) int32 {
	max := vars[0]
	for _, i := range vars {
		if max < i {
			max = i
		}
	}
	return max
}

func MinOf(vars ...XYZ) XYZ {
	return XYZ{
		X: minIntOf(sf.Map[XYZ, int32](vars, func(a XYZ) int32 {
			return a.X
		})...),
		Y: minIntOf(sf.Map[XYZ, int32](vars, func(a XYZ) int32 {
			return a.Y
		})...),
		Z: minIntOf(sf.Map[XYZ, int32](vars, func(a XYZ) int32 {
			return a.Z
		})...),
	}
}

func MaxOf(vars ...XYZ) XYZ {
	return XYZ{
		X: maxIntOf(sf.Map[XYZ, int32](vars, func(a XYZ) int32 {
			return a.X
		})...),
		Y: maxIntOf(sf.Map[XYZ, int32](vars, func(a XYZ) int32 {
			return a.Y
		})...),
		Z: maxIntOf(sf.Map[XYZ, int32](vars, func(a XYZ) int32 {
			return a.Z
		})...),
	}
}
