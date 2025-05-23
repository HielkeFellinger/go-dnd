package game_engine

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

const SpaceGameTest string = "../../content/default/entities.yml"
const SpaceGame string = "./content/default/entities.yml"

func loadGame(gameFile string) ecs.World {

	log.Println("Loading raw/base Game File")
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)

	data, err := os.ReadFile(gameFile)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Unmarshal raw/base Game Data")
	var game ecs.RawGameFile
	if err := yaml.Unmarshal(data, &game); err != nil {
		log.Fatalln(err)
	}

	idToUuidDict := make(map[string]uuid.UUID)
	uuidToEntityDict := make(map[uuid.UUID]ecs.Entity)
	log.Println("Parsing raw/base Game Items (Entities)")
	err, items := parseRawEntity(game.Items, idToUuidDict, uuidToEntityDict)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("Parsing raw/base Game Characters (Entities)")
	err, chars := parseRawEntity(game.Chars, idToUuidDict, uuidToEntityDict)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("Parsing raw/base Game Factions (Entities)")
	err, factions := parseRawEntity(game.Factions, idToUuidDict, uuidToEntityDict)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("Parsing raw/base Game Maps (Entities)")
	err, maps := parseRawEntity(game.Maps, idToUuidDict, uuidToEntityDict)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("Parsing raw/base Game Inventories (Entities)")
	err, inventories := parseRawEntity(game.Inventories, idToUuidDict, uuidToEntityDict)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("Parsing raw/base Game Map Content (Entities)")
	err, mapContent := parseRawEntity(game.MapContent, idToUuidDict, uuidToEntityDict)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("Parsing raw/base Game Other Content (Entities)")
	err, otherItems := parseRawEntity(game.Others, idToUuidDict, uuidToEntityDict)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("Parsing raw/base Game Items (Entity Components)")
	parseRawComponentsOfEntity(game, game.Items, uuidToEntityDict, idToUuidDict)
	log.Println("Parsing raw/base Game Characters (Entity Components)")
	parseRawComponentsOfEntity(game, game.Chars, uuidToEntityDict, idToUuidDict)
	log.Println("Parsing raw/base Game Factions (Entity Components)")
	parseRawComponentsOfEntity(game, game.Factions, uuidToEntityDict, idToUuidDict)
	log.Println("Parsing raw/base Game Maps (Entity Components)")
	parseRawComponentsOfEntity(game, game.Maps, uuidToEntityDict, idToUuidDict)
	log.Println("Parsing raw/base Game Inventory (Entity Components)")
	parseRawComponentsOfEntity(game, game.Inventories, uuidToEntityDict, idToUuidDict)
	log.Println("Parsing raw/base Game Map Content (Entity Components)")
	parseRawComponentsOfEntity(game, game.MapContent, uuidToEntityDict, idToUuidDict)
	log.Println("Parsing raw/base Game Other Items (Entity Components)")
	parseRawComponentsOfEntity(game, game.Others, uuidToEntityDict, idToUuidDict)

	// Add the fully updated Entities to the world
	log.Println("Done loading raw/base Game. Now filling world")
	world := ecs.NewBaseWorld()
	if errAdd := world.AddEntities(items); errAdd != nil {
		log.Fatalln(errAdd.Error())
	}
	if errAdd := world.AddEntities(chars); errAdd != nil {
		log.Fatalln(errAdd.Error())
	}
	if errAdd := world.AddEntities(factions); errAdd != nil {
		log.Fatalln(errAdd.Error())
	}
	if errAdd := world.AddEntities(maps); errAdd != nil {
		log.Fatalln(errAdd.Error())
	}
	if errAdd := world.AddEntities(inventories); errAdd != nil {
		log.Fatalln(errAdd.Error())
	}
	if errAdd := world.AddEntities(mapContent); errAdd != nil {
		log.Fatalln(errAdd.Error())
	}
	if errAdd := world.AddEntities(otherItems); errAdd != nil {
		log.Fatalln(errAdd.Error())
	}
	log.Println("Done filling world")
	return &world
}

func parseRawComponentsOfEntity(game ecs.RawGameFile, rawEntities []ecs.RawEntity, uuidToEntityDict map[uuid.UUID]ecs.Entity,
	idToUuidDict map[string]uuid.UUID) {
	for _, rawEntity := range rawEntities {
		match := uuidToEntityDict[idToUuidDict[rawEntity.Id]]

		for _, rawComponent := range rawEntity.Components {
			rawType := game.TypeTranslation[rawComponent.ComponentType]
			parsedType := ecs_components.MapIntToTypeV0(rawType)
			newComponent := ecs_components.MapTypeToConstructorFunction(parsedType)()

			// Check if relational component or regular
			if relComponent, ok := newComponent.(ecs.RelationalComponent); ok {
				// Find Target
				targetMatch := uuidToEntityDict[idToUuidDict[rawComponent.Params["entity"]]]
				if err := relComponent.LoadFromRawComponentRelation(rawComponent, targetMatch); err != nil {
					log.Fatalf(err.Error() + " Raw Entity ID: " + rawEntity.Id)
				}
			} else {
				if err := newComponent.LoadFromRawComponent(rawComponent); err != nil {
					log.Fatalf(err.Error() + " Raw Entity ID: " + rawEntity.Id)
				}
			}

			if err := match.AddComponent(newComponent); err != nil {
				log.Fatalf(err.Error() + " Raw Entity ID: " + rawEntity.Id)
			}
		}
	}
}

func parseRawEntity(rawEntities []ecs.RawEntity,
	idToUuidDict map[string]uuid.UUID,
	uuidToEntityDict map[uuid.UUID]ecs.Entity) (error, []ecs.Entity) {

	entities := make([]ecs.Entity, len(rawEntities))

	for index, rawEntity := range rawEntities {
		// Test if ID is unique
		if _, match := idToUuidDict[rawEntity.Id]; match {
			log.Fatalf("Duplicate Raw Entity with ID: '%s'. Could not load game", rawEntity.Id)
		}

		// Create and fill the new Entity
		entity := ecs.NewEntity()

		// Test if ID is a UUID, if so use this!
		if savedUuid, err := uuid.Parse(rawEntity.Id); err == nil {
			entity.Id = savedUuid
		}

		err := entity.LoadFromRawEntity(rawEntity)
		if err != nil {
			log.Fatalf(err.Error() + " Raw Entity ID: " + rawEntity.Id)
		}

		// Update the maps
		idToUuidDict[rawEntity.Id] = entity.Id
		uuidToEntityDict[entity.Id] = &entity
		entities[index] = &entity
	}

	return nil, entities
}
