package block

// ChestMinecart (item only) is a minecart with a chest inside.
type ChestMinecart struct{}

func (cm ChestMinecart) EncodeItem() (name string, meta int16) {
	return "minecraft:chest_minecart", 0
}
