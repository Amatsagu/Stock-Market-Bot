package other

import (
	"fmt"
	"index-bot/models"
)

func Sif[T any](statement bool, truthValue T, elseValue T) T {
	if statement {
		return truthValue
	} else {
		return elseValue
	}
}

func FormatDateParam(date models.DateParam) string {
	return fmt.Sprintf(
		"%d-%s%d-%s%d",
		date.Year,
		Sif(date.Month < 10, "0", ""),
		date.Month,
		Sif(date.Day < 10, "0", ""),
		date.Day,
	)
}
