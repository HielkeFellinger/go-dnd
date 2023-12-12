package game_engine

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
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
	var game RawGameFile
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

	return world
}

func parseRawEntity(rawEntities []RawEntity,
	idToUuidDict map[string]uuid.UUID,
	uuidToEntityDict map[uuid.UUID]ecs.Entity) (error, []ecs.Entity) {

	entities := make([]ecs.Entity, len(rawEntities))

	for rawEntity := range rawEntities {

		// Test if ID is unique
		// Test if

	}

	return nil, entities
}
