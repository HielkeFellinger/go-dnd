package ecs

type RawGameFile struct {
	Version         string         `yaml:"version"`
	TypeTranslation map[string]int `yaml:"type_translation"`
	Items           []RawEntity    `yaml:"base_items"`
	Maps            []RawEntity    `yaml:"base_maps"`
	Chars           []RawEntity    `yaml:"base_characters"`
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
