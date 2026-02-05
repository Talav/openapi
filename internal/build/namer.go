package build

import (
	"reflect"
	"strings"
	"unicode/utf8"
)

// schemaNamer provides schema names for types. It uses the type name
// when possible, ignoring the package name. If the type is generic, e.g.
// `MyType[SubType]`, then the brackets are removed like `MyTypeSubType`.
// If the type is unnamed, then the name hint is used.
// Note: if you plan to use types with the same name from different packages,
// you should implement your own namer function to prevent issues. Nested
// anonymous types can also present naming issues.
func schemaNamer(t reflect.Type, hint string) string {
	name := deref(t).Name()

	if name == "" {
		name = hint
	}

	// Better support for lists, so e.g. `[]int` becomes `ListInt`.
	name = strings.ReplaceAll(name, "[]", "List[")

	result := ""
	for _, part := range strings.FieldsFunc(name, func(r rune) bool {
		// Split on special characters. Note that `,` is used when there are
		// multiple inputs to a generic type.
		return r == '[' || r == ']' || r == '*' || r == ','
	}) {
		// Split fully qualified names like `github.com/foo/bar.Baz` into `Baz`.
		fqn := strings.Split(part, ".")
		base := fqn[len(fqn)-1]

		// Add to result, and uppercase for better scalar support (`int` -> `Int`).
		// Use unicode-aware uppercase to support non-ASCII characters.
		r, size := utf8.DecodeRuneInString(base)
		result += strings.ToUpper(string(r)) + base[size:]
	}
	name = result

	return name
}
