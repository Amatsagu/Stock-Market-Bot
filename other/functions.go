package other

func sif[T any](statement bool, truthValue T, elseValue T) T {
	if statement {
		return truthValue
	} else {
		return elseValue
	}
}
