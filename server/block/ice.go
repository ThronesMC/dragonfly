package block

import (
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
)

// Ice is a solid block similar to packed ice.
type Ice struct {
	transparent
}

var iceHash = NextHash()

func (i Ice) Hash() uint64 {
	return iceHash
}

func (i Ice) Model() world.BlockModel {
	return model.Solid{}
}

// BreakInfo ...
func (i Ice) BreakInfo() BreakInfo {
	return newBreakInfo(5, alwaysHarvestable, pickaxeEffective, silkTouchOnlyDrop(i))
}

// LightDiffusionLevel always returns 2.
func (Ice) LightDiffusionLevel() uint8 {
	return 2
}

// LightEmissionLevel returns 2.
func (Ice) LightEmissionLevel() uint8 {
	return 2
}

// Friction ...
func (i Ice) Friction() float64 {
	return 0.98
}

// EncodeItem ...
func (Ice) EncodeItem() (name string, meta int16) {
	return "minecraft:ice", 0
}

// EncodeBlock ...
func (Ice) EncodeBlock() (string, map[string]any) {
	return "minecraft:ice", nil
}
