package lib

func Size(a int) string {
	switch {
	case a < 0:
		return "negative"
	case a == 0:
		return "zero"
	}

	return "positive"
}
