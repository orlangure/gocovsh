package general

func Covered() string {
	return "covered"
}

func NotCovered() string {
	return "not covered"
}

func SecondCovered() string {
	switch true {
	default:
	}

	return "covered"
}

type useless struct{}
