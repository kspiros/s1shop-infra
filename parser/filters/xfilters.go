package filters

var filterList map[FilterType]IFilter

type IFilter interface {
	Execute(cvalue interface{}, value interface{}) bool
}

type FilterType string

func GetFilter(filter FilterType) (IFilter, bool) {
	c, found := filterList[filter]
	return c, found
}

func RegisterFilter(filter FilterType, f IFilter) {
	if filterList == nil {
		filterList = make(map[FilterType]IFilter, 0)
	}
	filterList[filter] = f
}
