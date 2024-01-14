package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// Hopper is a redstone component. It can be used to collect items that drop on it.
type Hopper struct {
	solid

	// Facing is the direction the pumpkin is facing.
	Facing cube.Direction
}

var hopperHash = NextHash()

func (h Hopper) Hash() uint64 {
	return hopperHash | uint64(h.Facing)<<8
}

// UseOnBlock ...
func (h Hopper) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(w, pos, face, h)
	if !used {
		return
	}
	h.Facing = user.Rotation().Direction().Opposite()

	place(w, pos, h, user, ctx)
	return placed(ctx)
}

// BreakInfo ...
func (h Hopper) BreakInfo() BreakInfo {
	return newBreakInfo(1, alwaysHarvestable, pickaxeEffective, oneOf(h))
}

// EncodeItem ...
func (h Hopper) EncodeItem() (name string, meta int16) {
	return "minecraft:hopper", 0
}

// EncodeBlock ...
func (h Hopper) EncodeBlock() (string, map[string]any) {
	return "minecraft:hopper", nil
}

func allHoppers() (hoppers []world.Block) {
	for i := cube.Direction(0); i <= 3; i++ {
		hoppers = append(hoppers, Hopper{Facing: i})
	}
	return
}
