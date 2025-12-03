package parser

import "errors"

var (
	ErrTargetNotFound = errors.New("target not found")
	ErrParsingFailed  = errors.New("parsing failed")
)
