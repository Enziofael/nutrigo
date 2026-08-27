package query

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
const UpdateTemplate = `
UPDATE >table<
SET
	>>sets<<
[from:
FROM
	>>from_tables<<
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
