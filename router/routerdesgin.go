package router

import (
	"reflect"
	"regexp"
)

type controllerInfo struct {
	regex          *regexp.Regexp
	params         map[int]string
	controllerType reflect.Type
}
