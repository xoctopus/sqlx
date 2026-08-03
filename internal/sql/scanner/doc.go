// Package scanner maps *sql.Rows into Go values for sqlx.
//
// Public helpers such as [github.com/xoctopus/sqlx/pkg/helper.Scan] and
// QueryAndScan delegate here. Adaptor catalog loaders also use [Scan] to read
// information_schema rows into structs.
//
// # Scan
//
//	err := scanner.Scan(ctx, rows, &user)   // single row; NOTFOUND if empty
//	err := scanner.Scan(ctx, rows, &users)  // []T appends each row
//	err := scanner.Scan(ctx, rows, &count)  // scalar / sql.Scanner
//
// [Scan] closes rows when finished. The destination must be a pointer (or a
// custom [ScanIter]). Empty rows with a single-value target return the package
// errors NOTFOUND code.
//
// # Destination shapes
//
//   - struct pointer: match row columns to `db` fields (and table__column
//     aliases from joins via [frag.Alias]); NULL columns are skipped via
//     nullable.NullIgnoreScanner
//   - [WithColumnReceivers]: explicit column-name → scan target map
//   - sql.Scanner / scalar pointer: rows.Scan into that value
//   - slice pointer: one New/Next cycle per row ([SliceScanIter])
//   - [ScanIter]: custom New + Next pipeline ([ScanIterFor])
//
// Unmatched columns are discarded ([nullable.EmptyScanner]).
//
// # nullable
//
// Subpackage nullable provides [nullable.NullIgnoreScanner] (nil src leaves
// the destination unchanged) and EmptyScanner for unused columns.
package scanner
