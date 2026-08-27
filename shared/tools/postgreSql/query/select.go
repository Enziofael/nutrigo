package query

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
const SelectTemplate = `
SELECT >distinct<
	>>columns<<
FROM 
	>table<
>>joins<<
[where:
WHERE 
	>where_condition<
]
[group:
GROUP BY
	>>group_columns<<
]
[having:
HAVING
	>having_condition<
]
[order:
ORDER BY
	>>order_columns<<
]
[fetch:
FETCH FIRST >fetch< >ties<
]
[offset:
OFFSET >offset<
]
[lock:
FOR >lock<
]
`

//[join:[inner_join:INNER][left_join:LEFT][right_join:RIGHT][outer_join:FULL OUTER][cross_join:CROSS] JOIN >join_table< [on:ON >join_condition<][using:USING (>join_condition<)]]

// Anchors:
//
//	- 1: "a" {Baked select query string. Use [SelectTemplate]}
//	- 1: "operator" {"UNION" | "UNION ALL" | "INTERSECT" | "EXCEPT"}
//	- 1: "b" {Baked select query string. Use [SelectTemplate]}
const SelectSetTemplate = `
(>a<)
>operator<
(>b<)
`

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
const SelectSetSortingTemplate = `
>set<
[order:
ORDER BY
	>>order_columns<<
]
[fetch:
FETCH FIRST >fetch< >ties<
]
[offset:
OFFSET >offset<
]
`
