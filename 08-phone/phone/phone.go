package phone

import "fmt"

type Phone struct {
	Id     int
	Number string
}

func (p Phone) String() string {
	return fmt.Sprintf("%d, %s", p.Id, p.Number)
}
