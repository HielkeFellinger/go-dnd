package ecs_components

import "github.com/hielkefellinger/go-dnd/app/ecs"

func MapIntToTypeV0(rawId int) uint64 {
	switch rawId {
	case 0:
		return ecs.PositionComponentType
	case 1:
		return ecs.AreaComponentType
	case 2:
		return ecs.RangeComponentType
	case 3:
		return ecs.DamageComponentType
	case 4:
		return ecs.RestoreComponentType
	case 5:
		return ecs.ItemComponentType
	case 6:
		return ecs.AmountComponentType
	case 7:
		return ecs.WeightComponentType
	case 8:
		return ecs.InventoryComponentType
	case 9:
		return ecs.LevelComponentType
	case 10:
		return ecs.TypeComponentType
	case 11:
		return ecs.ValutaComponentType
	case 12:
		return ecs.ResourceComponentType
	case 13:
		return ecs.TransportComponentType
	case 14:
		return ecs.TurnDistanceComponentType
	case 15:
		return ecs.VisibilityComponentType
	case 16:
		return ecs.HealthComponentType
	case 17:
		return ecs.StatComponentType
	case 18:
		return ecs.FactionComponentType
	case 19:
		return ecs.CharacterComponentType
	case 20:
		return ecs.MapComponentType
	case 21:
		return ecs.ImageComponentType
	case 22:
		return ecs.PlayerComponentType
	case 23:
		return ecs.LootComponentType
	case 24:
		return ecs.BlockerComponentType
	case 25:
		return ecs.MapContentComponentType

	case 40:
		return ecs.ControlsRelationComponentType
	case 41:
		return ecs.HasRelationComponentType
	case 42:
		return ecs.RequiresRelationComponentType
	case 43:
		return ecs.CreatesRelationComponentType
	case 44:
		return ecs.FilterRelationComponentType
	case 45:
		return ecs.MapItemRelationComponentType
	case 46:
		return ecs.MapLinkRelationComponentType
	default:
		return ecs.UnknownComponentType
	}
}

func MapTypeToConstructorFunction(componentType uint64) func() ecs.Component {
	switch componentType {
	case ecs.PositionComponentType:
		return NewPositionComponent
	case ecs.AreaComponentType:
		return NewAreaComponent
	case ecs.RangeComponentType:
		return NewRangeComponent
	case ecs.DamageComponentType:
		return NewDamageComponent
	case ecs.RestoreComponentType:
		return NewRestoreComponent
	case ecs.ItemComponentType:
		return NewItemComponent
	case ecs.AmountComponentType:
		return NewAmountComponent
	case ecs.WeightComponentType:
		return NewWeightComponent
	case ecs.InventoryComponentType:
		return NewInventoryComponent
	case ecs.LevelComponentType:
		return NewLevelComponent
	case ecs.TypeComponentType:
		return NewTypeComponent
	case ecs.ValutaComponentType:
		return NewValutaComponent
	case ecs.ResourceComponentType:
		return NewResourceComponent
	case ecs.TransportComponentType:
		return NewTransportComponent
	case ecs.TurnDistanceComponentType:
		return NewTurnDistanceComponent
	case ecs.VisibilityComponentType:
		return NewVisibilityComponent
	case ecs.HealthComponentType:
		return NewHealthComponent
	case ecs.StatComponentType:
		return NewStatComponent
	case ecs.FactionComponentType:
		return NewFactionComponent
	case ecs.CharacterComponentType:
		return NewCharacterComponent
	case ecs.MapComponentType:
		return NewMapComponent
	case ecs.ImageComponentType:
		return NewImageComponent
	case ecs.PlayerComponentType:
		return NewPlayerComponent
	case ecs.LootComponentType:
		return NewLootComponent
	case ecs.BlockerComponentType:
		return NewBlockerComponent
	case ecs.MapContentComponentType:
		return NewMapContentComponent

	case ecs.ControlsRelationComponentType:
		return NewControlsRelationComponent
	case ecs.HasRelationComponentType:
		return NewHasRelationComponent
	case ecs.RequiresRelationComponentType:
		return NewRequiresRelationComponent
	case ecs.CreatesRelationComponentType:
		return NewCreatesRelationComponent
	case ecs.FilterRelationComponentType:
		return NewFilterRelationComponent
	case ecs.MapItemRelationComponentType:
		return NewMapItemRelationComponent
	case ecs.MapLinkRelationComponentType:
		return NewMapLinkRelationComponent
	default:
		return nil
	}
}

func MapStringToFilterMode(filterMode string) ecs.FilterMode {
	switch filterMode {
	case "block":
		return ecs.BlockFilterMode
	case "allow":
		return ecs.AllowFilterMode
	case "less":
		return ecs.LessThanFilterMode
	case "more":
		return ecs.MoreThanFilterMode
	case "shouldHave":
		return ecs.ShouldHaveFilterMode
	case "shouldNotHave":
		return ecs.ShouldNotHaveFilterMode
	default:
		return ecs.UnknownFilterMode
	}
}
func MapFilterModeToString(filterMode ecs.FilterMode) string {
	switch filterMode {
	case ecs.BlockFilterMode:
		return "block"
	case ecs.AllowFilterMode:
		return "allow"
	case ecs.LessThanFilterMode:
		return "less"
	case ecs.MoreThanFilterMode:
		return "more"
	case ecs.ShouldHaveFilterMode:
		return "shouldHave"
	case ecs.ShouldNotHaveFilterMode:
		return "shouldNotHave"

	default:
		return "unknown"
	}
}
