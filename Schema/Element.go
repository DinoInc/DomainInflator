package Schema

type elementType string

const (
	Null    elementType = "null"
	Boolean elementType = "bool"
	Object  elementType = "object"
	Array   elementType = "array"
	Number  elementType = "i32"
	Str     elementType = "string"
)
