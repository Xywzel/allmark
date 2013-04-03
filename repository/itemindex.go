package repository

import (
	"fmt"
)

type ItemIndex struct {
	path  string
	items []*Item
}

func NewItemIndex(directory string, items []*Item) *ItemIndex {
	return &ItemIndex{
		path:  directory,
		items: items,
	}
}

func (itemIndex *ItemIndex) String() string {
	return fmt.Sprintf("Item Index %s", itemIndex.path)
}

func (itemIndex *ItemIndex) Walk(walkFunc func(item *Item)) {
	for _, item := range itemIndex.items {
		item.Walk(walkFunc)
	}
}

func (itemIndex *ItemIndex) Path() string {
	return itemIndex.path
}
