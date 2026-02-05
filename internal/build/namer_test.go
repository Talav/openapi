package build

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchemaNamer(t *testing.T) {
	// Test types
	type User struct {
		Name string
	}
	type Product struct {
		ID int
	}

	tests := []struct {
		name string
		typ  reflect.Type
		hint string
		want string
	}{
		{
			name: "named struct type",
			typ:  reflect.TypeOf(User{}),
			hint: "",
			want: "User",
		},
		{
			name: "named struct type with hint",
			typ:  reflect.TypeOf(User{}),
			hint: "CustomHint",
			want: "User", // type name takes priority
		},
		{
			name: "pointer to named type",
			typ:  reflect.TypeOf((*User)(nil)),
			hint: "",
			want: "User", // deref removes pointer
		},
		{
			name: "double pointer",
			typ:  reflect.TypeOf((**User)(nil)),
			hint: "",
			want: "User", // deref removes all pointers
		},
		{
			name: "primitive int",
			typ:  reflect.TypeOf(0),
			hint: "",
			want: "Int",
		},
		{
			name: "primitive string",
			typ:  reflect.TypeOf(""),
			hint: "",
			want: "String",
		},
		{
			name: "primitive bool",
			typ:  reflect.TypeOf(false),
			hint: "",
			want: "Bool",
		},
		{
			name: "pointer to int",
			typ:  reflect.TypeOf((*int)(nil)),
			hint: "",
			want: "Int",
		},
		{
			name: "slice of int with hint",
			typ:  reflect.TypeOf([]int{}),
			hint: "[]int",
			want: "ListInt",
		},
		{
			name: "slice of string with hint",
			typ:  reflect.TypeOf([]string{}),
			hint: "[]string",
			want: "ListString",
		},
		{
			name: "slice of User with hint",
			typ:  reflect.TypeOf([]User{}),
			hint: "[]User",
			want: "ListUser",
		},
		{
			name: "pointer to slice with hint",
			typ:  reflect.TypeOf((*[]int)(nil)),
			hint: "[]int",
			want: "ListInt",
		},
		{
			name: "slice without hint",
			typ:  reflect.TypeOf([]int{}),
			hint: "",
			want: "", // unnamed type without hint returns empty
		},
		{
			name: "anonymous struct with hint",
			typ:  reflect.TypeOf(struct{ Name string }{}),
			hint: "CreateUserRequest",
			want: "CreateUserRequest",
		},
		{
			name: "anonymous struct without hint",
			typ:  reflect.TypeOf(struct{ Name string }{}),
			hint: "",
			want: "",
		},
		{
			name: "map type with hint",
			typ:  reflect.TypeOf(map[string]int{}),
			hint: "map[string]int",
			want: "MapStringInt",
		},
		{
			name: "map without hint",
			typ:  reflect.TypeOf(map[string]int{}),
			hint: "",
			want: "", // unnamed type without hint returns empty
		},
		{
			name: "channel type with hint",
			typ:  reflect.TypeOf((chan int)(nil)),
			hint: "chan int",
			want: "Chan int", // space is preserved, not split
		},
		{
			name: "channel without hint",
			typ:  reflect.TypeOf((chan int)(nil)),
			hint: "",
			want: "", // unnamed type without hint returns empty
		},
		{
			name: "function type with hint",
			typ:  reflect.TypeOf(func() {}),
			hint: "func()",
			want: "Func()", // parentheses preserved, not split
		},
		{
			name: "function without hint",
			typ:  reflect.TypeOf(func() {}),
			hint: "",
			want: "", // unnamed type without hint returns empty
		},
		{
			name: "interface type",
			typ:  reflect.TypeOf((*interface{})(nil)).Elem(),
			hint: "",
			want: "",
		},
		{
			name: "array type with hint",
			typ:  reflect.TypeOf([5]int{}),
			hint: "[5]int",
			want: "5Int",
		},
		{
			name: "array without hint",
			typ:  reflect.TypeOf([5]int{}),
			hint: "",
			want: "", // unnamed type without hint returns empty
		},
		{
			name: "fully qualified name simulation",
			typ:  reflect.TypeOf(User{}),
			hint: "github.com/example.User",
			want: "User", // type name takes priority, but if hint was used it would extract "User"
		},
		{
			name: "hint with package path",
			typ:  reflect.TypeOf(struct{}{}),
			hint: "github.com/example.User",
			want: "User", // extracts base name after last dot
		},
		{
			name: "hint with multiple dots",
			typ:  reflect.TypeOf(struct{}{}),
			hint: "github.com/example/v2.User",
			want: "User",
		},
		{
			name: "hint with brackets",
			typ:  reflect.TypeOf(struct{}{}),
			hint: "List[Int]",
			want: "ListInt", // brackets removed, parts concatenated
		},
		{
			name: "hint with generic-like syntax",
			typ:  reflect.TypeOf(struct{}{}),
			hint: "MyType[SubType]",
			want: "MyTypeSubType",
		},
		{
			name: "hint with multiple generic params",
			typ:  reflect.TypeOf(struct{}{}),
			hint: "Map[Key,Value]",
			want: "MapKeyValue",
		},
		{
			name: "hint with asterisk",
			typ:  reflect.TypeOf(struct{}{}),
			hint: "*User",
			want: "User", // asterisk removed
		},
		{
			name: "hint with slice brackets",
			typ:  reflect.TypeOf(struct{}{}),
			hint: "[]int",
			want: "ListInt", // [] becomes List[
		},
		{
			name: "hint lowercase first letter",
			typ:  reflect.TypeOf(struct{}{}),
			hint: "user",
			want: "User", // first letter uppercased
		},
		{
			name: "hint already uppercase",
			typ:  reflect.TypeOf(struct{}{}),
			hint: "User",
			want: "User",
		},
		{
			name: "hint with unicode",
			typ:  reflect.TypeOf(struct{}{}),
			hint: "пользователь", // Russian for "user"
			want: "Пользователь", // unicode-aware uppercase
		},
		{
			name: "empty hint with named type",
			typ:  reflect.TypeOf(Product{}),
			hint: "",
			want: "Product",
		},
		{
			name: "complex hint",
			typ:  reflect.TypeOf(struct{}{}),
			hint: "github.com/example.List[*User]",
			want: "ListUser", // extracts base, removes brackets and asterisk
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := schemaNamer(tt.typ, tt.hint)
			assert.Equal(t, tt.want, got)
		})
	}
}
