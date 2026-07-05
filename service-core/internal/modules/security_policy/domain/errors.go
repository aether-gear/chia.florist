package domain

import "errors"

var (
	ErrRuleNotFound    = errors.New("waf rule not found")
	ErrDuplicateRuleID = errors.New("waf rule id already exists")
	ErrInvalidPattern  = errors.New("waf rule pattern is not a valid regular expression")
)

var (
	ErrIPRecordNotFound = errors.New("ip access control record not found")
)

var (
	ErrFilterEntryNotFound = errors.New("filter config entry not found")
)
