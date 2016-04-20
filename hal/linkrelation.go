// go2hal v0.1.0
// Copyright (c) 2016 Patrick Moule
// License: MIT

package hal

import "errors"

type LinkRelation struct {
	name      string
	curieLink *LinkObject
}

func NewLinkRelation(name string) (*LinkRelation, error) {
	if name == "" {
		return nil, errors.New("LinkRelation requires a name value.")
	}

	return &LinkRelation{name: name}, nil
}

func (lr *LinkRelation) Name() string {
	return lr.name
}

func (lr *LinkRelation) FullName() string {
	if lr.curieLink == nil {
		return lr.Name()
	}

	return lr.curieLink.Name + ":" + lr.Name()
}

func (lr *LinkRelation) SetCurieLink(curieLink *LinkObject) {
	lr.curieLink = curieLink
}

func (lr *LinkRelation) CurieLink() LinkObject {
	return *lr.curieLink
}