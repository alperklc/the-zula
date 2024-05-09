package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAreArraysEqual(t *testing.T) {
	t.Run("it returns true if given two arrays of string are identical", func(t *testing.T) {
		// arrange
		array1 := []string{"test", "lorem", "ipsum"}
		array2 := []string{"test", "lorem", "ipsum"}

		// act
		result := AreArraysEqual(array1, array2)

		// assert
		assert.Equal(t, true, result)
	})

	t.Run("it returns false if given two arrays of string are not identical", func(t *testing.T) {
		// arrange
		array1 := []string{"test", "lorem", "ipsum", "dolor"}
		array2 := []string{"test", "lorem", "ipsum"}

		// act
		result := AreArraysEqual(array1, array2)

		// assert
		assert.Equal(t, false, result)
	})
}

func TestFilterFieldsOfObject(t *testing.T) {
	t.Run("it returns an object that contains selected and existing fields", func(t *testing.T) {
		// arrange
		obj := map[string]interface{}{
			"Firstname": "name",
			"Lastname":  "lastname",
			"username":  "namename",
			"age":       29,
			"eyes":      "brown",
		}
		fields := []string{"Firstname", "Lastname", "height"}

		// act
		filteredObject := FilterFieldsOfObject(fields, obj)

		// assert
		assert.Contains(t, filteredObject, "Firstname")
		assert.Contains(t, filteredObject, "Lastname")
		assert.NotContains(t, filteredObject, "height")
		assert.NotContains(t, filteredObject, "age")
	})
}

func TestParseInternalLinksFromNote(t *testing.T) {
	t.Run("it parses internal links from note content", func(t *testing.T) {
		// arrange
		content := `
		  this is a markdown-like text for testing this regex parser function

		  this one is to an external page: [link text](www.google.com)

		  this one is to an internal page: [internal link](/notes/uid-of-note)
		  and this one is to an another internal page: [internal link](/notes/uid-of-note-2)
		`

		// act
		parsedLinks := ParseInternalLinksFromNote(content)

		// assert
		assert.Contains(t, parsedLinks, "uid-of-note")
		assert.Contains(t, parsedLinks, "uid-of-note-2")
		assert.NotContains(t, parsedLinks, "www.google.com")
	})
}
