// Foo

// Read and write DRS files.
//
// # Format
//
//	+--------+-------+--------------+
//	| Offset | Size  | Name         |
//	|--------|------ |--------------|
//	|      0 |   40  | Header       |
//	|     40 |       | TableInfo[0] |
//	| ...                           |
//	|        |       | TableInfo[N] |
//	|        |       | FileInfo[0]  |
//	| ...                           |
//	|        |       | FileInfo[M]  |
//	|        |       | FileData[0]  |
//	| ...                           |
//	|        |       | FileData[M]  |
//	+--------+-------+--------------+
package drs
