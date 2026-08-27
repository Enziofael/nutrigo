package query

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	TEMPLATE_NONE = iota
	// Anchors:
	//
	//	- 1: "distinct" {"" | "DISTINCT"}
	//	- n: "columns" {like "id," | "name"}
	//	- 1: "table" {table name or subquery in brackets: "users", "(SELECT ...)"}
	//	- n: "joins" {like "LEFT JOIN ..."}
	//	- ?: "where"
	//	- 1: "where" > "where_condition"
	//	- ?: "group"
	//	- n: "group" > "group_columns" {like "region," | "city"}
	//	- ?: "having"
	//	- 1: "having" > "having_condition" {like "type = client AND age > 18"}
	//	- ?: "order"
	//	- n: "order" > "order_columns" {like "id DESC NULLS FIRST," | "name ASC"}
	//	- ?: "fetch"
	//	- 1: "fetch" > "fetch"  {"ROW" | "n ROWS"}
	//	- 1: "fetch" > "ties" {"WITH TIES" | "ONLY"}
	//	- ?: "offset"
	//	- 1: "offset" > "offset" {n}
	//	- ?: "lock"
	//	- 1: "lock" > "lock" {"UPDATE" | "SHARE"}
	TEMPLATE_SELECT

	// Anchors:
	//
	//	- 1: "table" {table name or subquery in brackets: "users", "(SELECT ...)"}
	//	- n: "columns" {like "id," | "name"}
	//	- n: "values" {like "(3,'Jane')," "(5,'Mike')"}
	//	- ?: "returning"
	//	- n: "returning" -> "returning_columns" {like "id," | "name"}
	TEMPLATE_INSERT

	// Anchors:
	//
	//	- 1: "table" {table name or subquery in brackets: "users", "(SELECT ...)"}
	//	- n: "sets" {like "name = 'Mike'," | "age = 18"}
	//	- ?: "from"
	//	- n: "from" > "from_tables" {like "users," | "cities"}
	//	- ?: "where"
	//	- 1: "where" > "where_condition" {like "age < 18 OR name = "Jane"}
	//	- ?: "returning"
	//	- n: "returning" -> "returning_columns" {like "id," | "name"}
	TEMPLATE_UPDATE

	//	- 1: "table" {table name or subquery in brackets: "users", "(SELECT ...)"}
	//	- ?: "using"
	//	- n: "using" > "using_tables" {like "users," | "cities"}
	//	- ?: "where"
	//	- 1: "where" > "where_condition" {like "age < 18 OR name = "Jane"}
	//	- ?: "returning"
	//	- n: "returning" -> "returning_columns" {like "id," | "name"}
	TEMPLATE_DELETE

	// Anchors:
	//
	//	- 1: "a" {Baked select query string. Use [SelectTemplate]}
	//	- 1: "operator" {"UNION" | "UNION ALL" | "INTERSECT" | "EXCEPT"}
	//	- 1: "b" {Baked select query string. Use [SelectTemplate]}
	TEMPLATE_SELECT_SET

	// Anchors:
	//
	//	- 1: "set" {Baked select set query string. Use [SelectSetTemplate]}
	//	- ?: "order"
	//	- n: "order" > "order_columns"
	//	- ?: "fetch"
	//	- 1: "fetch" > "fetch" {"ROW" | "n ROWS"}
	//	- 1: "fetch" > "ties" {"WITH TIES" | "ONLY"}
	//	- ?: "offset"
	//	- 1: "offset" > "offset" {n}
	TEMPLATE_SELECT_SET_SORTING
)

var selectTemplatePrototype = MustCompileTemplate(SelectTemplate)
var selectSetTemplatePrototype = MustCompileTemplate(SelectSetTemplate)
var selectSetSortingTemplatePrototype = MustCompileTemplate(SelectSetSortingTemplate)

var insertTemplatePrototype = MustCompileTemplate(InsertTemplate)
var updateTemplatePrototype = MustCompileTemplate(UpdateTemplate)
var deleteTemplatePrototype = MustCompileTemplate(DeleteTemplate)

// Use constants:
//   - [TEMPLATE_SELECT]
//   - [TEMPLATE_INSERT]
//   - [TEMPLATE_UPDATE]
//   - [TEMPLATE_DELETE]
//   - [TEMPLATE_SELECT_SET]
//   - [TEMPLATE_SELECT_SET_SORTING]
//
// nil will be returned otherwise
func ClonePrototype(queryType int) *QueryTemplate {
	switch queryType {
	case TEMPLATE_SELECT:
		return selectTemplatePrototype.Clone()
	case TEMPLATE_INSERT:
		return insertTemplatePrototype.Clone()
	case TEMPLATE_UPDATE:
		return updateTemplatePrototype.Clone()
	case TEMPLATE_DELETE:
		return deleteTemplatePrototype.Clone()
	case TEMPLATE_SELECT_SET:
		return selectSetTemplatePrototype.Clone()
	case TEMPLATE_SELECT_SET_SORTING:
		return selectSetSortingTemplatePrototype.Clone()
	default:
		return nil
	}
}

type QueryTemplate struct {
	template        string
	singleAnchors   map[string]*TemplateAnchor
	multipleAnchors map[string]*TemplateAnchor
	optionalAnchors map[string]*TemplateAnchor
}

type TemplateAnchor struct {
	anchorStart  int
	anchorEnd    int
	anchorType   byte
	contentStart int
	contentEnd   int
}

const (
	anchor_none = iota
	anchor_single
	anchor_multiple
	anchor_optional
)

var anchorSingleRegexp = regexp.MustCompile(`>([a-zA-Z_]+)<`)
var anchorMultipleRegexp = regexp.MustCompile(`(?m)^(?:.*)(>>([a-zA-Z_]+)<<)(?:.*)$`)
var anchorOptionalRegexp = regexp.MustCompile(`(?s)\[([a-zA-Z_]+):(.*?)\]`)
var anchorClearRegexp = regexp.MustCompile(`(?s)(?:>)(?:.*?)(?:<)`)

func (src *QueryTemplate) Clone() *QueryTemplate {
	return &QueryTemplate{
		template:        src.template,
		singleAnchors:   cloneAnchors(src.singleAnchors),
		multipleAnchors: cloneAnchors(src.multipleAnchors),
		optionalAnchors: cloneAnchors(src.optionalAnchors),
	}
}

func cloneAnchors(anchors map[string]*TemplateAnchor) map[string]*TemplateAnchor {
	cloned := make(map[string]*TemplateAnchor, len(anchors))
	for key, val := range anchors {
		a := *val
		cloned[key] = &a
	}
	return cloned
}

func MustCompileTemplate(template string) *QueryTemplate {
	qt := QueryTemplate{
		template:        template,
		singleAnchors:   make(map[string]*TemplateAnchor),
		multipleAnchors: make(map[string]*TemplateAnchor),
		optionalAnchors: make(map[string]*TemplateAnchor),
	}

	matchesSingle := anchorSingleRegexp.FindAllStringSubmatchIndex(template, -1)
	for _, match := range matchesSingle {
		a := TemplateAnchor{
			anchorType:  anchor_single,
			anchorStart: match[0],
			anchorEnd:   match[1],
		}
		name := template[match[2]:match[3]]
		if qt.singleAnchors[name] != nil {
			panic("single anchors with the same name is not allowed: " + name + "occurs at " + strconv.Itoa(qt.singleAnchors[name].anchorStart) + " and " + strconv.Itoa(a.anchorStart))
		}
		qt.singleAnchors[name] = &a
	}

	matchesMultiple := anchorMultipleRegexp.FindAllStringSubmatchIndex(template, -1)
	for _, match := range matchesMultiple {
		a := TemplateAnchor{
			anchorType:   anchor_multiple,
			anchorStart:  match[0],
			anchorEnd:    match[1],
			contentStart: match[2],
			contentEnd:   match[3],
		}
		name := template[match[4]:match[5]]
		if qt.multipleAnchors[name] != nil {
			panic("multiple anchors with the same name is not allowed: " + name + "occurs at " + strconv.Itoa(qt.multipleAnchors[name].anchorStart) + " and " + strconv.Itoa(a.anchorStart))
		}
		qt.multipleAnchors[name] = &a
	}

	matchesOptional := anchorOptionalRegexp.FindAllStringSubmatchIndex(template, -1)
	for _, match := range matchesOptional {
		a := TemplateAnchor{
			anchorType:   anchor_optional,
			anchorStart:  match[0],
			anchorEnd:    match[1],
			contentStart: match[4],
			contentEnd:   match[5],
		}
		name := template[match[2]:match[3]]
		if qt.optionalAnchors[name] != nil {
			panic("optional anchors with the same name is not allowed: " + name + "occurs at " + strconv.Itoa(qt.optionalAnchors[name].anchorStart) + " and " + strconv.Itoa(a.anchorStart))
		}
		qt.optionalAnchors[name] = &a
	}

	return &qt
}

func (qt *QueryTemplate) MustRecompileTemplate(template string) {
	re := MustCompileTemplate(template)
	qt.template = template
	qt.singleAnchors = re.singleAnchors
	qt.multipleAnchors = re.multipleAnchors
	qt.optionalAnchors = re.optionalAnchors
}

func (qt *QueryTemplate) MustReplaceSingle(anchorName string, value string) {
	a := qt.singleAnchors[anchorName]
	if a == nil {
		panic("single anchor \"" + anchorName + "\" not found\ntemplate:\n" + qt.template)
	}

	var sb strings.Builder
	sb.WriteString(qt.template[:a.anchorStart])
	sb.WriteString(value)
	sb.WriteString(qt.template[a.anchorEnd:])

	qt.MustRecompileTemplate(sb.String())
}

func (qt *QueryTemplate) MustPlaceMultiple(anchorName string, value string) {
	a := qt.multipleAnchors[anchorName]
	if a == nil {
		panic("multiple anchor \"" + anchorName + "\" not found\ntemplate:\n" + qt.template)
	}

	var sb strings.Builder
	sb.WriteString(qt.template[:a.contentStart])
	sb.WriteString(value)
	sb.WriteString("\n")
	sb.WriteString(qt.template[a.anchorStart:])

	qt.MustRecompileTemplate(sb.String())
}

func (qt *QueryTemplate) MustEnableOptional(anchorName string) {
	a := qt.optionalAnchors[anchorName]
	if a == nil {
		panic("optional anchor \"" + anchorName + "\" not found\ntemplate:\n" + qt.template)
	}

	var sb strings.Builder
	sb.WriteString(qt.template[:a.anchorStart])
	sb.WriteString(qt.template[a.contentStart:a.contentEnd])
	sb.WriteString(qt.template[a.anchorEnd:])

	qt.MustRecompileTemplate(sb.String())
}

func (qt *QueryTemplate) Bake() string {
	var sb strings.Builder
	var cursor int

	matches := anchorClearRegexp.FindAllStringIndex(qt.template, -1)
	for _, match := range matches {
		sb.WriteString(qt.template[cursor:match[0]])
		cursor = match[1]
	}
	sb.WriteString(qt.template[cursor:])
	qt.template = sb.String()

	sb.Reset()
	cursor = 0

	matches = anchorOptionalRegexp.FindAllStringIndex(qt.template, -1)
	for _, match := range matches {
		sb.WriteString(qt.template[cursor:match[0]])
		cursor = match[1]
		fmt.Printf("=================optional tidy iteration\n%s\n", sb.String())
	}
	sb.WriteString(qt.template[cursor:])

	lines := strings.Split(sb.String(), "\n")
	tidyLines := []string{}
	for _, line := range lines {
		line = strings.TrimRight(line, " \t")
		if line != "" {
			tidyLines = append(tidyLines, line)
		}
	}
	qt.template = strings.Join(tidyLines, "\n")

	return qt.template
}

// Shortcut for [*QueryTemplate.MustRecompileTemplate]
func (qt *QueryTemplate) MRecompile(template string) {
	qt.MustRecompileTemplate(template)
}

// Shortcut for [*QueryTemplate.MustReplaceSingle]
func (qt *QueryTemplate) MSingle(anchorName string, value string) {
	qt.MustReplaceSingle(anchorName, value)
}

// Shortcut for [*QueryTemplate.MustPlaceMultiple]
func (qt *QueryTemplate) MMultiple(anchorName string, value string) {
	qt.MustPlaceMultiple(anchorName, value)
}

// Shortcut for [*QueryTemplate.MustEnableOptional]
func (qt *QueryTemplate) MOptiobal(anchorName string) {
	qt.MustEnableOptional(anchorName)
}

// Shortcut for [*QueryTemplate.MustReplaceSingle]
func (qt *QueryTemplate) MSingleIf(anchorName string, value string, shouldReplace bool) {
	if !shouldReplace {
		return
	}
	qt.MustReplaceSingle(anchorName, value)
}

// Shortcut for [*QueryTemplate.MustPlaceMultiple]
func (qt *QueryTemplate) MMultipleIf(anchorName string, value string, shouldReplace bool) {
	if !shouldReplace {
		return
	}
	qt.MustPlaceMultiple(anchorName, value)
}

// Shortcut for [*QueryTemplate.MustEnableOptional]
func (qt *QueryTemplate) MOptiobalIf(anchorName string, shouldReplace bool) {
	if !shouldReplace {
		return
	}
	qt.MustEnableOptional(anchorName)
}
