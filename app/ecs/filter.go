package ecs

type Filter interface {
	GetFilterMode() FilterMode
	GetFilterComponents() []Component
	DoesFilterApplyValid(entity Entity) bool
}

type BaseFilter struct {
	filterMode       FilterMode
	filterValue      int
	filterComponents []Component
}

func NewBaseFilter(filterMode FilterMode, filterValue int) BaseFilter {
	return BaseFilter{
		filterMode:       filterMode,
		filterValue:      filterValue,
		filterComponents: make([]Component, 0),
	}
}

func (f *BaseFilter) GetFilterMode() FilterMode {
	return f.filterMode
}

func (f *BaseFilter) GetFilterComponents() []Component {
	return f.filterComponents
}

func (f *BaseFilter) DoesFilterApplyValid(entity Entity) bool {
	if f.filterMode == AllowFilterMode || f.filterMode == ShouldHaveFilterMode {
		for _, component := range f.filterComponents {
			if !entity.HasComponentType(component.ComponentType()) {
				return false
			}
		}
	} else if f.filterMode == BlockFilterMode || f.filterMode == ShouldNotHaveFilterMode {
		for _, component := range f.filterComponents {
			if entity.HasComponentType(component.ComponentType()) {
				return false
			}
		}
	} else if f.filterMode == LessThanFilterMode {
		hasMatch := false
		for _, component := range f.filterComponents {
			if entity.HasComponentType(component.ComponentType()) {
				matches := entity.GetAllComponentsOfType(component.ComponentType())
				for _, match := range matches {
					hasMatch = hasMatch || match.IsLessThanValue(f.filterValue)
				}
			}
		}
		if !hasMatch {
			return hasMatch
		}
	} else if f.filterMode == MoreThanFilterMode {
		hasMatch := false
		for _, component := range f.filterComponents {
			if entity.HasComponentType(component.ComponentType()) {
				matches := entity.GetAllComponentsOfType(component.ComponentType())
				for _, match := range matches {
					hasMatch = hasMatch || match.IsMoreThanValue(f.filterValue)
				}
			}
		}
		if !hasMatch {
			return hasMatch
		}
	}
	return true
}
