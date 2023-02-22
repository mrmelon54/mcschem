package marshall

import (
	"errors"
	"github.com/MrMelon54/mcschem/region"
	"github.com/MrMelon54/mcschem/structure"
	"github.com/MrMelon54/mcschem/xyz"
	"golang.org/x/exp/maps"
	"io"
	"time"
)

var ErrNoLitematicRegionExists = errors.New("no litematic region exists")

type LitematicFormat struct {
	name string
	Data structure.LitematicData
}

func NewLitematic(version, mcVersion int32, name, desc, author string) *LitematicFormat {
	t := time.Now().Unix()
	return &LitematicFormat{
		name: "litematic",
		Data: structure.LitematicData{
			Version:              version,
			MinecraftDataVersion: mcVersion,
			Metadata: structure.LitematicMetadata{
				EnclosingSize: xyz.XYZ{},
				Author:        author,
				Description:   desc,
				Name:          name,
				RegionCount:   0,
				TimeCreated:   t,
				TimeModified:  t,
				TotalBlocks:   0,
				TotalVolume:   0,
			},
			Regions: map[string]structure.LitematicRegion{},
		},
	}
}

func ParseLitematic(reader io.Reader) (*LitematicFormat, error) {
	name, data, err := decodeGeneric[structure.LitematicData](reader)
	if err != nil {
		return nil, err
	}
	return &LitematicFormat{name: name, Data: data}, nil
}

func (l *LitematicFormat) Save(writer io.Writer) error {
	return encodeGeneric[structure.LitematicData](writer, l.name, l.Data)
}

func (l *LitematicFormat) ListRegions() []string {
	return maps.Keys(l.Data.Regions)
}

func (l *LitematicFormat) GetRegion(name string) (*region.Access, error) {
	if r, ok := l.Data.Regions[name]; ok {
		return region.FromLitematicRegion(r)
	}
	return nil, ErrNoLitematicRegionExists
}

func (l *LitematicFormat) CreateRegion(name string) (*region.Access, error) {
	r := region.Empty()
	return r, l.AddRegion(name, r)
}

func (l *LitematicFormat) AddRegion(name string, access *region.Access) error {
	o, err := access.ToLitematicRegion()
	if err != nil {
		return err
	}
	l.Data.Regions[name] = o
	return nil
}

func (l *LitematicFormat) CalculateMetadata() {
	l.Data.Metadata.RegionCount = int32(len(l.Data.Regions))
	l.Data.Metadata.EnclosingSize = l.calculateEnclosingSize()
	l.Data.Metadata.TotalVolume = l.Data.Metadata.EnclosingSize.Count()
}

func (l *LitematicFormat) calculateEnclosingSize() xyz.XYZ {
	sizes := make([]xyz.XYZ, 0)
	for _, r := range l.Data.Regions {
		sizes = append(sizes, r.Position)
		sizes = append(sizes, r.Position.Add(r.Size))
	}
	max := xyz.MaxOf(sizes...)
	min := xyz.MinOf(sizes...)
	return max.Sub(min)
}
