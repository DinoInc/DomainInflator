// Optional Extension Element - found in all resources.
struct Extension {
	1 : required string url,
	2 : optional string valueId
}

// Base definition for all elements in a resource.
struct Element {
    1: required string id,
    2: list<Extension> extension
}