package game_engine

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

const SpaceGame string = "../../content/space/entities.yml"

func loadGame(gameFile string) ecs.BaseWorld {

	log.Println("Loading raw/base Game File")
	data, err := os.ReadFile(gameFile)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Unmarshal raw/base Game Data")
	var game ecs.RawGameFile
	if err := yaml.Unmarshal(data, &game); err != nil {
		log.Fatalln(err)
	}
	log.Printf("The data '%v'", game)

	idToUuidDict := make(map[string]uuid.UUID)
	uuidToEntityDict := make(map[uuid.UUID]ecs.Entity)
	world := ecs.BaseWorld{}
	log.Println("Parsing raw/base Game Items (Entities)")
	err, items := parseRawEntity(game.Items, idToUuidDict, uuidToEntityDict)
	if err != nil {
		log.Fatalln(err.Error())
	}
	world.Entities = append(world.Entities, items...)

	log.Println("Parsing raw/base Game Characters (Entities)")
	err, chars := parseRawEntity(game.Chars, idToUuidDict, uuidToEntityDict)
	if err != nil {
		log.Fatalln(err.Error())
	}
	world.Entities = append(world.Entities, chars...)

	log.Println("Parsing raw/base Game Maps (Entities)")
	err, maps := parseRawEntity(game.Maps, idToUuidDict, uuidToEntityDict)
	if err != nil {
		log.Fatalln(err.Error())
	}
	world.Entities = append(world.Entities, maps...)

	log.Println("Parsing raw/base Game Items (Entity Components)")
	parseRawComponentsOfEntity(game, game.Items, uuidToEntityDict, idToUuidDict)
	log.Println("Parsing raw/base Game Characters (Entity Components)")
	parseRawComponentsOfEntity(game, game.Chars, uuidToEntityDict, idToUuidDict)
	log.Println("Parsing raw/base Game Maps (Entity Components)")
	parseRawComponentsOfEntity(game, game.Maps, uuidToEntityDict, idToUuidDict)

	return world
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

			match.AddComponent(newComponent)
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
		entity.WithName(rawEntity.Name).WithDescription(rawEntity.Description)

		// Update the maps
		idToUuidDict[rawEntity.Id] = entity.Id
		uuidToEntityDict[entity.Id] = &entity
		entities[index] = &entity
	}

	return nil, entities
}
