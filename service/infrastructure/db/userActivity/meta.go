package useractivity

import (
	"fmt"
	"math"
)

func getPaginationRange(count, page, pageSize int) string {
	if count == 0 {
		return ""
	}

	var from int = (page-1)*pageSize + 1
	var to int = int(math.Min(float64(page*pageSize), float64(count)))

	return fmt.Sprintf("%d - %d", from, to)
}
