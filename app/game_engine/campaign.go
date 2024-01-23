package game_engine

type CampaignPool interface {
	GetId() uint
	GetLeadId() string
	GetEngine() Engine
	TransmitEventMessage(message EventMessage)
	GetAllClientIds(filterOut ...string) []string
}

type CampaignClient interface {
	GetId() string
	IsLead() bool
}
