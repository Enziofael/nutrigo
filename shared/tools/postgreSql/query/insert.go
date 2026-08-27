package query

// Anchors:
//
//	- 1: "table" {table name or subquery in brackets: "users", "(SELECT ...)"}
//	- n: "columns" {like "id," | "name"}
//	- n: "values" {like "(3,'Jane')," "(5,'Mike')"}
//	- ?: "returning"
//	- n: "returning" -> "returning_columns" {like "id," | "name"}
const InsertTemplate = `
INSERT INTO >table<(
	>>columns<<
) VALUES 
	>>values<<
[returning:
RETURNING
	>>returning_columns<<
]
`
