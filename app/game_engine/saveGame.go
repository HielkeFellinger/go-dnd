package game_engine

import (
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func saveGame(world ecs.World, gameFile string) error {

	rawGame := ecs.RawGameFile{}
	rawGame.Version = "0.1"

	// Save TypeTranslation;
	rawGame.TypeTranslation = make(map[string]int, len(ecs.TypeNameToNthBit))
	for _, value := range ecs.TypeNameToNthBit {
		rawGame.TypeTranslation[value.Name] = int(value.BitNr)
	}

	// Split per type;
	log.Println("   - Saving Item (Entities)")
	rawGame.Items = parseEntityIntoRawEntity(world.GetItemEntities())
	log.Println("   - Saving Character (Entities)")
	rawGame.Chars = parseEntityIntoRawEntity(world.GetCharacterEntities())
	log.Println("   - Saving Map (Entities)")
	rawGame.Maps = parseEntityIntoRawEntity(world.GetMapEntities())
	log.Println("   - Saving Faction (Entities)")
	rawGame.Factions = parseEntityIntoRawEntity(world.GetFactionEntities())
	log.Println("   - Saving Inventories (Entities)")
	rawGame.Inventories = parseEntityIntoRawEntity(world.GetInventoryEntities())
	log.Println("   - Saving Map Content (Entities)")
	rawGame.MapContent = parseEntityIntoRawEntity(world.GetMapContentEntities())
	log.Println("   - Saving Other (Entities)")
	rawGame.Others = parseEntityIntoRawEntity(world.GetOtherEntities())

	gameFileContent, err := yaml.Marshal(rawGame)
	if err != nil {
		return err
	}

	// Save the file
	log.Println("Attempting to save campaign to file")
	if writeErr := os.WriteFile(gameFile, gameFileContent, 0644); writeErr != nil {
		return writeErr
	}
	log.Println("Saved the game")

	return nil
}

func parseEntityIntoRawEntity(entities []ecs.Entity) []ecs.RawEntity {
	rawEntities := make([]ecs.RawEntity, 0)

	for _, entity := range entities {
		// Skip nil entity; may be a leftover of slices.delete not reducing the total size of the underlying array.
		if entity == nil {
			continue
		}

		rawEntity := ecs.RawEntity{
			Id:          entity.GetId().String(),
			Name:        entity.GetName(),
			Description: entity.GetDescription(),
			Components:  parseComponentsToRawComponents(entity.GetAllComponents()),
		}
		rawEntities = append(rawEntities, rawEntity)
	}
	return rawEntities
}

func parseComponentsToRawComponents(components []ecs.Component) []ecs.RawComponent {
	rawComponents := make([]ecs.RawComponent, 0)

	for _, component := range components {
		// Skip nil components; may be a leftover of slices.delete not reducing the total size of the underlying array.
		if component == nil {
			continue
		}

		if rawComponent, err := component.ParseToRawComponent(); err == nil {
			rawComponents = append(rawComponents, rawComponent)
		} else {
			log.Printf("Error parsing component of type: '%v' with error message: '%s'", component, err.Error())
		}
	}
	return rawComponents
}
