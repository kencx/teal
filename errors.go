package teal

import "errors"

var ErrDoesNotExist = errors.New("the item does not exist")
var ErrNoRows = errors.New("no items found")
