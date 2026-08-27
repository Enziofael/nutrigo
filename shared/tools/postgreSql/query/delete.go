package query

//	- 1: "table" {table name or subquery in brackets: "users", "(SELECT ...)"}
//	- ?: "using"
//	- n: "using" > "using_tables" {like "users," | "cities"}
//	- ?: "where"
//	- 1: "where" > "where_condition" {like "age < 18 OR name = "Jane"}
//	- ?: "returning"
//	- n: "returning" -> "returning_columns" {like "id," | "name"}
const DeleteTemplate = `
DELETE FROM >table<
[using:
USING
	>>using_tables<<
]
[where:
WHERE
	>where_condition<
]
[returning:
RETURNING
	>>returning_columns<<
]
`
