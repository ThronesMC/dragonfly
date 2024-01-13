package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

type Magma struct {
	solid
}

var magmaHash = NextHash()

func (m Magma) Hash() uint64 {
	return magmaHash
}

// LightDiffusionLevel always returns 2.
func (Magma) LightDiffusionLevel() uint8 {
	return 2
}

// LightEmissionLevel returns 3.
func (Magma) LightEmissionLevel() uint8 {
	return 3
}

// EntityInside ...
func (m Magma) EntityInside(_ cube.Pos, _ *world.World, e world.Entity) {
	if _, ok := e.(flammableEntity); ok {
		if l, ok := e.(livingEntity); ok && !l.AttackImmune() {
			l.Hurt(1, MagmaDamageSource{})
		}
	}
}

func (Magma) EncodeItem() (name string, meta int16) {
	return "minecraft:magma", 0
}

func (Magma) EncodeBlock() (string, map[string]any) {
	return "minecraft:magma", nil
}

// MagmaDamageSource is used for damage caused by being on Magma.
type MagmaDamageSource struct{}

func (MagmaDamageSource) ReducedByResistance() bool { return true }
func (MagmaDamageSource) ReducedByArmour() bool     { return true }
func (MagmaDamageSource) Fire() bool                { return true }
