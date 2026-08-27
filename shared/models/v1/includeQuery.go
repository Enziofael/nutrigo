package models

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
)

// AllowedParams() should return string with allowed params (a-z, _),
// separeted by comma
//
//	func (Request) AllowedParams() string {
//		return "param,param,nested_includable<>"
//	}
//
// GetIncludeQuery() should pass main IncludableRequest
// and its nested ones to the NewIncludeQuery()
//
//	func (req Request) GetIncludeQuery() IncludeQuery {
//		return NewIncludeQuery(req, NestedIncludable{})
//	}
type IncludableRequest interface {
	AllowedParams() string
	GetIncludeQuery() IncludeQuery
}

type IncludeQuery struct {
	Name         string
	Params       *[]string
	Nested       *[]IncludeQuery
	nestedRegexp *regexp.Regexp
}

var commonIncludeQueryRegexp *regexp.Regexp = regexp.MustCompile(`^(?:(?:[a-z_<>]+)[,]?)*$`)
var paramRegexp *regexp.Regexp = regexp.MustCompile(`^[a-z_]+$`)
var nestedRegexp *regexp.Regexp = regexp.MustCompile(`^([a-z_]+)<([a-z_<>,]*)>$`)

func preSplit(s string) []string {
	var parts []string
	var current strings.Builder
	depth := 0
	for _, ch := range s {
		switch ch {
		case ',':
			if depth == 0 {
				parts = append(parts, current.String())
				current.Reset()
			} else {
				current.WriteRune(ch)
			}
		case '<':
			depth++
			current.WriteRune(ch)
		case '>':
			depth--
			current.WriteRune(ch)
		default:
			current.WriteRune(ch)
		}
	}
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	return parts
}

func ValidateIncludeQuery(s string) bool {
	return commonIncludeQueryRegexp.MatchString(s)
}

func NewIncludeQuery(req IncludableRequest, nested ...IncludableRequest) IncludeQuery {
	if !ValidateIncludeQuery(req.AllowedParams()) {
		log.Panicf("%s in NewIncludeQuery(): %s", ErrInvalidIncludeQuery.Error(), req.AllowedParams())
	}

	parts := preSplit(req.AllowedParams())
	var params []string
	var nesteds []IncludeQuery

	nestedIdx := 0

	for _, part := range parts {
		switch {
		case paramRegexp.MatchString(part):
			params = append(params, part)

		case nestedRegexp.MatchString(part):
			if nestedIdx >= len(nested) {
				log.Panicf("Not enough nested requests for param: %s", part)
			}
			submatch := nestedRegexp.FindStringSubmatch(part)
			name := submatch[1]
			nestedRegex := regexp.MustCompile(name + `<([a-z_<>,]*)>`)
			nestedQuery := nested[nestedIdx].GetIncludeQuery()
			nestedQuery.nestedRegexp = nestedRegex
			nesteds = append(nesteds, nestedQuery)
			nestedIdx++
		default:
			log.Panicf("Unhandled param string: %s", part)
		}
	}

	if len(params) == 0 {
		params = nil
	}
	if len(nesteds) == 0 {
		nesteds = nil
	}

	t := reflect.TypeOf(req)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	name := t.Name()

	return IncludeQuery{
		Name:   name,
		Params: &params,
		Nested: &nesteds,
	}
}

func (iq IncludeQuery) Parse(s string) (IncludeQuery, error) {
	if !ValidateIncludeQuery(s) {
		return IncludeQuery{}, fmt.Errorf("%w:\"%s\"", ErrInvalidIncludeQuery, s)
	}

	var parsedParams []string
	var parsedNesteds []IncludeQuery

	parts := preSplit(s)

	for _, part := range parts {

		matchedParam := false
		if iq.Params != nil {
			for _, param := range *iq.Params {
				if param == part {
					parsedParams = append(parsedParams, param)
					matchedParam = true
					break
				}
			}
		}
		if matchedParam {
			continue
		}

		if iq.Nested != nil {
			for _, nested := range *iq.Nested {
				if nested.nestedRegexp != nil {
					submatch := nested.nestedRegexp.FindStringSubmatch(part)
					if submatch != nil {
						nestedContent := submatch[1]
						parsedNested, err := nested.Parse(nestedContent)
						if err != nil {
							return parsedNested, err //forwarding the results of an unsuccessful nested request parsing
						}
						parsedNesteds = append(parsedNesteds, parsedNested)
						break
					}
				}
			}
		}
	}

	if len(parsedParams) == 0 {
		parsedParams = nil
	}
	if len(parsedNesteds) == 0 {
		parsedNesteds = nil
	}

	return IncludeQuery{
		Name:         iq.Name,
		Params:       &parsedParams,
		Nested:       &parsedNesteds,
		nestedRegexp: iq.nestedRegexp,
	}, nil
}

func (iq IncludeQuery) Blank() IncludeQuery {
	return IncludeQuery{
		Name:         iq.Name,
		Params:       nil,
		Nested:       nil,
		nestedRegexp: iq.nestedRegexp,
	}
}

var ErrNestedIncludeNotFound error = errors.New("include query not found in")

func (iq IncludeQuery) FindNested(name string) (IncludeQuery, error) {
	// No Nested
	if iq.Nested == nil {
		return IncludeQuery{}, fmt.Errorf("\"%s\" %w \"%s\"", name, ErrNestedIncludeNotFound, iq.Name)
	}

	// Search in Nested
	for _, nested := range *iq.Nested {
		if nested.Name == name {
			return nested, nil
		}
	}

	// Not found in Nested
	return IncludeQuery{}, fmt.Errorf("\"%s\" %w \"%s\"", name, ErrNestedIncludeNotFound, iq.Name)
}

const BFSThreshold = 10

func (iq IncludeQuery) FindNestedRecursive(name string) (IncludeQuery, error) {
	// No Nested
	if iq.Nested == nil {
		return IncludeQuery{}, fmt.Errorf("\"%s\" %w \"%s\"", name, ErrNestedIncludeNotFound, iq.Name)
	}

	// Check Nested
	for _, nested := range *iq.Nested {
		if nested.Name == name {
			return nested, nil
		}
	}
	// Search in Nested
	for _, nested := range *iq.Nested {
		if found, err := nested.FindNestedRecursive(name); err == nil {
			return found, nil
		}
	}

	// Not found in Nested
	return IncludeQuery{}, fmt.Errorf("\"%s\" %w \"%s\"", name, ErrNestedIncludeNotFound, iq.Name)
}

func (iq IncludeQuery) ContainsParam(name string) bool {
	if iq.Params == nil {
		return false
	}
	for _, param := range *iq.Params {
		if param == name {
			return true
		}
	}
	return false
}

func (iq IncludeQuery) ContainsParamRecursive(name string) bool {
	if iq.ContainsParam(name) {
		return true
	}
	if iq.Nested == nil {
		return false
	}
	for _, nested := range *iq.Nested {
		if nested.ContainsParamRecursive(name) {
			return true
		}
	}
	return false
}

func (iq IncludeQuery) Print(level int) {
	tab := strings.Repeat("\t", level)
	fmt.Printf("%s%s\n", tab, iq.Name)

	if iq.Params != nil && len(*iq.Params) > 0 {
		fmt.Printf("\t%sparams:\n", tab)
		for _, param := range *iq.Params {
			fmt.Printf("\t\t%s\"%s\"\n", tab, param)
		}
	}

	if iq.Nested != nil && len(*iq.Nested) > 0 {
		fmt.Printf("\t%snested:\n", tab)
		for _, nested := range *iq.Nested {
			nested.Print(level + 2)
		}
	}
}
