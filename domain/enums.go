package domain

import (
	"fmt"
	"strings"

	"github.com/johannes-kuhfuss/services_utils/api_error"
)

type EnumItem struct {
	index int
	value string
}
type EnumList struct {
	items []EnumItem
}

func (e *EnumList) Value(index int) (string, api_error.ApiErr) {
	for _, item := range e.items {
		if item.index == index {
			return item.value, nil
		}
	}
	return "", api_error.NewNotFoundError(fmt.Sprintf("No item with index %v found", index))
}

func (e *EnumList) Index(value string) (int, api_error.ApiErr) {
	for index, item := range e.items {
		if strings.EqualFold(value, item.value) {
			return index, nil
		}
	}
	return 0, api_error.NewNotFoundError(fmt.Sprintf("No item with name %v found", value))
}

func (e *EnumList) Values() []string {
	var names []string
	for _, item := range e.items {
		names = append(names, item.value)
	}
	return names
}

func (e *EnumList) AsMap() map[int]string {
	m := make(map[int]string)
	for _, item := range e.items {
		m[item.index] = item.value
	}
	return m
}

func (e *EnumList) FromMap(m map[int]string) {
	var eItem EnumItem
	for index, item := range m {
		eItem.index = index
		eItem.value = item
		e.items = append(e.items, eItem)
	}
}
