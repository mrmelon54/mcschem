package structure

type SchematicData struct {
	Version     int32
	DataVersion int32
	Metadata    SchematicMetadata
	Offset      [3]int32
	Width       int32
	Height      int32
	Length      int32
	PaletteMax  int32
	Palette     map[string]int32
	BlockData   []byte
}

type SchematicMetadata struct {
	Name      string
	Author    string
	Generator string
	Date      int64
	WEOriginX int32
	WEOriginY int32
	WEOriginZ int32
	WEOffsetX int32
	WEOffsetY int32
	WEOffsetZ int32
}
