// Package frag builds and flattens SQL fragments with bound arguments.
//
// A [Fragment] yields query text chunks and args via [Fragment.Frag].
// [Collect] concatenates them into a final query string and argument list
// for database/sql execution. [builder] statements are Fragments; this package
// is the low-level composition layer underneath.
//
//	q, args := frag.Collect(ctx, frag.Query("SELECT * FROM t WHERE id = ?", 1))
//	// q = "SELECT * FROM t WHERE id = ?"  args = [1]
//
// # Core types
//
//   - [Fragment]: IsNil + Frag(ctx) → [Iter]
//   - [Iter]: yields (query piece, args) pairs
//   - [Collect] / [Stringify]: flatten a Fragment for execution or debug
//   - [IsNil] / [NonNil]: nil checks and filtering
//
// # Constructors
//
//   - [Query]: template with `?` positional and `@name` named placeholders
//     ([NamedArg] / [NamedArgs]). Nested Fragments, slices, and sequences
//     expand in place (e.g. `IN (?)` with []int → `IN (?,?,?)`).
//
//   - [Lit] / [Empty]: literal SQL with no args / empty fragment
//
//   - [Arg] / [ArgIter] / [ArgIterFunc] / [Values]: expand a value into
//     placeholders; supports [CustomValueArg], driver.Valuer, slices, seqs
//
//   - [Func]: wrap a context-aware iterator factory as a Fragment
//
//     frag.Query("IN (?)", []int{1, 2, 3})
//     // IN (?,?,?)  [1 2 3]
//
//     frag.Query("time > @min AND time < @max", sql.Named("min", 10), sql.Named("max", 20))
//     // time > ? AND time < ?  [10, 20]
//
// # Composition
//
//   - [Compose] / [ComposeSeq]: join fragments with a separator (skips nil)
//   - [Block]: wrap a fragment in parentheses
//   - [BlockWithoutBrackets]: join fragments with commas, no surrounding ()
//
// # Other
//
//   - [Alias]: stable short alias for table__column names
package frag
