package ecs

const (
	UnknownComponentType      uint64 = 0
	PositionComponentType     uint64 = 1 << 0
	AreaComponentType         uint64 = 1 << 1
	RangeComponentType        uint64 = 1 << 2
	DamageComponentType       uint64 = 1 << 3
	RestoreComponentType      uint64 = 1 << 4
	ItemComponentType         uint64 = 1 << 5
	AmountComponentType       uint64 = 1 << 6
	WeightComponentType       uint64 = 1 << 7
	InventoryComponentType    uint64 = 1 << 8
	LevelComponentType        uint64 = 1 << 9
	TypeComponentType         uint64 = 1 << 10
	ValutaComponentType       uint64 = 1 << 11
	ResourceComponentType     uint64 = 1 << 12
	TransportComponentType    uint64 = 1 << 13
	TurnDistanceComponentType uint64 = 1 << 14
	VisibilityComponentType   uint64 = 1 << 15
	HealthComponentType       uint64 = 1 << 16
	StatComponentType         uint64 = 1 << 17
	FactionComponentType      uint64 = 1 << 18
	CharacterComponentType    uint64 = 1 << 19
	MapComponentType          uint64 = 1 << 20
	ImageComponentType        uint64 = 1 << 21
	PlayerComponentType       uint64 = 1 << 22
	LootComponentType         uint64 = 1 << 23
	BlockerComponentType      uint64 = 1 << 24
	MapContentComponentType   uint64 = 1 << 25

	/* Relational ComponentTypes */

	ControlsRelationComponentType uint64 = 1 << 40
	HasRelationComponentType      uint64 = 1 << 41
	RequiresRelationComponentType uint64 = 1 << 42
	CreatesRelationComponentType  uint64 = 1 << 43
	FilterRelationComponentType   uint64 = 1 << 44
	MapItemRelationComponentType  uint64 = 1 << 45
	MapLinkRelationComponentType  uint64 = 1 << 46
)

// @todo merge with const? maybe...
var TypeNameToNthBit = map[uint64]TypeTranslationStruct{
	PositionComponentType:     newTypeTranslation("Position", 0),
	AreaComponentType:         newTypeTranslation("Area", 1),
	RangeComponentType:        newTypeTranslation("Range", 2),
	DamageComponentType:       newTypeTranslation("Damage", 3),
	RestoreComponentType:      newTypeTranslation("Restore", 4),
	ItemComponentType:         newTypeTranslation("Item", 5),
	AmountComponentType:       newTypeTranslation("Amount", 6),
	WeightComponentType:       newTypeTranslation("Weight", 7),
	InventoryComponentType:    newTypeTranslation("Inventory", 8),
	LevelComponentType:        newTypeTranslation("Level", 9),
	TypeComponentType:         newTypeTranslation("Type", 10),
	ValutaComponentType:       newTypeTranslation("Valuta", 11),
	ResourceComponentType:     newTypeTranslation("Resource", 12),
	TransportComponentType:    newTypeTranslation("Transport", 13),
	TurnDistanceComponentType: newTypeTranslation("TurnDistance", 14),
	VisibilityComponentType:   newTypeTranslation("Visibility", 15),
	HealthComponentType:       newTypeTranslation("Health", 16),
	StatComponentType:         newTypeTranslation("Stat", 17),
	FactionComponentType:      newTypeTranslation("Faction", 18),
	CharacterComponentType:    newTypeTranslation("Character", 19),
	MapComponentType:          newTypeTranslation("Map", 20),
	ImageComponentType:        newTypeTranslation("Image", 21),
	PlayerComponentType:       newTypeTranslation("Player", 22),
	LootComponentType:         newTypeTranslation("Loot", 23),
	BlockerComponentType:      newTypeTranslation("Blocker", 24),
	MapContentComponentType:   newTypeTranslation("MapContent", 25),

	/* Relational ComponentTypes */

	ControlsRelationComponentType: newTypeTranslation("ControlsRelation", 40),
	HasRelationComponentType:      newTypeTranslation("HasRelation", 41),
	RequiresRelationComponentType: newTypeTranslation("RequiresRelation", 42),
	CreatesRelationComponentType:  newTypeTranslation("CreatesRelation", 43),
	FilterRelationComponentType:   newTypeTranslation("FilterRelation", 44),
	MapItemRelationComponentType:  newTypeTranslation("MapItemRelation", 45),
	MapLinkRelationComponentType:  newTypeTranslation("MapLinkRelation", 46),
}

type TypeTranslationStruct struct {
	Name  string
	BitNr uint64
}

func newTypeTranslation(name string, bitNr uint64) TypeTranslationStruct {
	return TypeTranslationStruct{
		Name:  name,
		BitNr: bitNr,
	}
}
