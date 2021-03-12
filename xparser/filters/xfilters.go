package filters

var filterList map[FilterType]interface{}

type IFilter interface {
	Execute(cvalue interface{}, value interface{}) bool
}

type IFilterTotals interface {
	CalcTotals(field string, cvalue interface{}, data *[]map[string]interface{}) interface{}
}

type FilterType string

func GetFilter(filter FilterType) (IFilter, bool) {
	c, found := filterList[filter]
	if v, ok := c.(IFilter); ok {
		return v, found
	}
	return nil, found
}

func FilterSupportsTotals(filter FilterType) (IFilterTotals, bool) {
	c, found := filterList[filter]
	if v, ok := c.(IFilterTotals); ok {
		return v, found
	}
	return nil, found
}

func RegisterFilter(filter FilterType, f IFilter) {
	if filterList == nil {
		filterList = make(map[FilterType]interface{}, 0)
	}
	filterList[filter] = f
}
