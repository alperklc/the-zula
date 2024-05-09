package utils

import "regexp"

func AreArraysEqual(arrayOne, arrayTwo []string) bool {
	if len(arrayOne) != len(arrayTwo) {
		return false
	}

	equality := true
	for i, elementFromArrayOne := range arrayOne {
		if elementFromArrayOne != arrayTwo[i] {
			equality = false
		}
	}

	return equality
}

func ArrayContainsString(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func FilterFieldsOfObject(fields []string, obj interface{}) interface{} {
	input, _ := obj.(map[string]interface{})

	var objectWithAllowedFields = make(map[string]interface{})
	for nameOfField, valueOfField := range input {
		if fieldAllowed := ArrayContainsString(fields, nameOfField); fieldAllowed {
			objectWithAllowedFields[nameOfField] = valueOfField
		}
	}

	return objectWithAllowedFields
}

func ParseInternalLinksFromNote(noteContent string) []string {
	re := regexp.MustCompile(`\[(.*)]\((\/notes\/)(.*)\)`)
	match := re.FindAllStringSubmatch(noteContent, -1)

	var parsedIds []string
	for i := range match {
		parsedIds = append(parsedIds, match[i][3])
	}

	return parsedIds
}
