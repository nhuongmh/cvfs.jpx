package tsunami

type Card struct {
	Question string
}

func (c *Card) Answer() string {
	return ""
}
