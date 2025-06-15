package entity

import "fmt"

// The set of roles that can be used.
var (
	User    = newEntity("USER")
	Product = newEntity("PRODUCT")
)

// Set of known entities.
var entities = make(map[string]Entity)

// Entity represents a domain in the system.
type Entity struct {
	value string
}

func newEntity(entity string) Entity {
	e := Entity{entity}
	entities[entity] = e
	return e
}

// String returns the name of the role.
func (e Entity) String() string {
	return e.value
}

// Equal provides support for the go-cmp package and testing.
func (e Entity) Equal(d2 Entity) bool {
	return e.value == d2.value
}

// MarshalText provides support for logging and any marshal needs.
func (e Entity) MarshalText() ([]byte, error) {
	return []byte(e.value), nil
}

// Parse parses the string value and returns a role if one exists.
func Parse(value string) (Entity, error) {
	entity, exists := entities[value]
	if !exists {
		return Entity{}, fmt.Errorf("invalid domain %q", value)
	}

	return entity, nil
}

// MustParse parses the string value and returns a role if one exists.
// If an error occurs, the function panics.
func MustParse(value string) Entity {
	entity, err := Parse(value)
	if err != nil {
		panic(err)
	}

	return entity
}
