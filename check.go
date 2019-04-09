package health

type Check struct {
	Name   string
	Status State
	Data   map[string]interface{}
}
