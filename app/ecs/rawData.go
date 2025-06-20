package ecs

type RawGameFile struct {
	Version         string         `yaml:"version"`
	TypeTranslation map[string]int `yaml:"type_translation"`
	Engine          RawEngine      `yaml:"engine"`
	Items           []RawEntity    `yaml:"base_items"`
	Maps            []RawEntity    `yaml:"base_maps"`
	Factions        []RawEntity    `yaml:"base_factions"`
	Chars           []RawEntity    `yaml:"base_characters"`
	Inventories     []RawEntity    `yaml:"base_inventories"`
	MapContent      []RawEntity    `yaml:"base_map_content"`
	Others          []RawEntity    `yaml:"others"`
}

type RawEngine struct {
	Systems []RawSystem `yaml:"systems"`
}

type RawSystem struct {
}

type RawComponent struct {
	ComponentType string            `yaml:"type"`
	Params        map[string]string `yaml:"params"`
}

type RawEntity struct {
	Id          string         `yaml:"id"`
	Name        string         `yaml:"name"`
	Description string         `yaml:"description"`
	Components  []RawComponent `yaml:"components"`
}
