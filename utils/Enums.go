package utils

import (
	"fmt"
	"strings"

	"github.com/johannes-kuhfuss/services_utils/api_error"
)

// Usage: var AccountTypes = Enum{[]EnumItem{{0, "Basic"}, {1, "Advanced"}}}

type EnumItem struct {
	idx int
	val string
}
type Enum struct {
	items []EnumItem
}

func (e *Enum) Value(i int) (string, api_error.ApiErr) {
	val, err := e.ItemByIndex(i)
	if err != nil {
		return "", err
	}
	return val.val, nil
}

func (e *Enum) Index(v string) (int, api_error.ApiErr) {
	val, err := e.ItemByValue(v)
	if err != nil {
		return 0, err
	}
	return val.idx, nil
}

func (e *Enum) Values() []string {
	var names []string
	for _, item := range e.items {
		names = append(names, item.val)
	}
	return names
}

func (e *Enum) AsMap() map[int]string {
	m := make(map[int]string)
	for _, item := range e.items {
		m[item.idx] = item.val
	}
	return m
}

func (e *Enum) FromMap(m map[int]string) {
	var eItem EnumItem
	for index, item := range m {
		eItem.idx = index
		eItem.val = item
		e.items = append(e.items, eItem)
	}
}

func (e *Enum) ItemByValue(v string) (*EnumItem, api_error.ApiErr) {
	for _, item := range e.items {
		if strings.EqualFold(v, item.val) {
			return &item, nil
		}
	}
	return nil, api_error.NewNotFoundError(fmt.Sprintf("No item with value %v found", v))
}

func (e *Enum) ItemByIndex(i int) (*EnumItem, api_error.ApiErr) {
	for _, item := range e.items {
		if item.idx == i {
			return &item, nil
		}
	}
	return nil, api_error.NewNotFoundError(fmt.Sprintf("No item with index %v found", i))
}
