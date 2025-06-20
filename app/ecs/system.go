package ecs

var defaultPriority uint64 = 32

type System interface {
	GetPriority() uint64
	GetSelectionFilters() []Filter
	GetResolutionFiltersFilters() []Filter
	ApplyGlobal(world World) (bool, error)
	ApplyLocally(entity Entity, world World) (bool, error)
}

type BaseSystem struct {
	selectionFilters  []Filter
	resolutionFilters []Filter
	priority          uint64
}

func NewBaseSystem() BaseSystem {
	return BaseSystem{
		priority:          defaultPriority,
		selectionFilters:  make([]Filter, 0),
		resolutionFilters: make([]Filter, 0),
	}
}

func (s *BaseSystem) GetPriority() uint64 {
	return s.priority
}

func (s *BaseSystem) GetSelectionFilters() []Filter {
	return s.selectionFilters
}

func (s *BaseSystem) GetResolutionFiltersFilters() []Filter {
	return s.selectionFilters
}

func (s *BaseSystem) ApplyGlobal(world World) (bool, error) {
	// @TODO
	return false, nil
}

func (s *BaseSystem) ApplyLocally(entity Entity, world World) (bool, error) {
	// @TODO
	return false, nil
}

func (s *BaseSystem) doesEntityMeetSelectionFiltering(entity Entity) bool {
	for _, filter := range s.selectionFilters {
		if filter.DoesFilterApplyValid(entity) {
			return false
		}
	}
	return true
}
