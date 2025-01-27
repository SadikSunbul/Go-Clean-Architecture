package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy_Struct(t *testing.T) {
	// Arrange
	type Person struct {
		Name string
		Age  int
	}
	src := Person{Name: "John", Age: 30}
	var dest Person

	// Act
	Copy(&dest, src)

	// Assert
	assert.Equal(t, src.Name, dest.Name)
	assert.Equal(t, src.Age, dest.Age)
}

func TestCopy_Map(t *testing.T) {
	// Arrange
	src := map[string]interface{}{
		"name": "John",
		"age":  30,
	}
	var dest map[string]interface{}

	// Act
	Copy(&dest, src)

	// Assert
	assert.Equal(t, src["name"], dest["name"])
	assert.Equal(t, float64(30), dest["age"])
}

func TestCopy_NestedStruct(t *testing.T) {
	// Arrange
	type Address struct {
		City    string
		Country string
	}
	type Person struct {
		Name    string
		Age     int
		Address Address
	}
	src := Person{
		Name: "John",
		Age:  30,
		Address: Address{
			City:    "Istanbul",
			Country: "Turkey",
		},
	}
	var dest Person

	// Act
	Copy(&dest, src)

	// Assert
	assert.Equal(t, src.Name, dest.Name)
	assert.Equal(t, src.Age, dest.Age)
	assert.Equal(t, src.Address.City, dest.Address.City)
	assert.Equal(t, src.Address.Country, dest.Address.Country)
}

func TestCopy_Slice(t *testing.T) {
	// Arrange
	src := []string{"one", "two", "three"}
	var dest []string

	// Act
	Copy(&dest, src)

	// Assert
	assert.Equal(t, len(src), len(dest))
	assert.Equal(t, src, dest)
}

func TestCopy_NilSource(t *testing.T) {
	// Arrange
	var src *struct{ Name string }
	var dest struct{ Name string }

	// Act
	Copy(&dest, src)

	// Assert
	assert.Empty(t, dest.Name)
}

func TestCopy_DifferentTypes(t *testing.T) {
	// Arrange
	type Source struct {
		Value int
	}
	type Dest struct {
		Value string
	}
	src := Source{Value: 123}
	var dest Dest

	// Act
	Copy(&dest, src)

	// Assert
	assert.Empty(t, dest.Value) // Farklı tipler kopyalanamaz
}
