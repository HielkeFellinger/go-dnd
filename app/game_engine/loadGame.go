package game_engine

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func loadGame() ecs.BaseWorld {

	log.Println("Loading raw/base Game File")
	data, err := os.ReadFile("../../content/space/entities.yml")
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
	log.Println("Parsing raw/base Game Items")
	err, items := parseRawEntity(game.Items, idToUuidDict, uuidToEntityDict)
	world.Entities = append(world.Entities, items...)

	log.Println("Parsing raw/base Game Characters")
	err, chars := parseRawEntity(game.Chars, idToUuidDict, uuidToEntityDict)
	world.Entities = append(world.Entities, chars...)

	log.Println("Parsing raw/base Game Maps")
	err, maps := parseRawEntity(game.Maps, idToUuidDict, uuidToEntityDict)
	world.Entities = append(world.Entities, maps...)

	// Add Components (TEST) <<<---- To function
	for _, rawEntity := range game.Items {
		match := uuidToEntityDict[idToUuidDict[rawEntity.Id]]

		for _, rawComponent := range rawEntity.Components {
			rawType := game.TypeTranslation[rawComponent.ComponentType]
			parsedType := ecs_components.MapIntToTypeV0(rawType)

			newComponent := ecs_components.MapTypeToConstructorFunction(parsedType)()
			if err := newComponent.LoadFromRawComponent(rawComponent); err != nil {
				log.Fatalf(err.Error())
			}
			match.AddComponent(newComponent)
		}
	}

	return world
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
