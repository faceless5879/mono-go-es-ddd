package world

// Define an interface with the methods you wish to expose
type Exposer interface {
	SetData(s string)
	PublicMethod() string
}

// Define a struct with both public and private methods
type myStruct struct {
	data string
}

// Implement the interface method
func (m *myStruct) PublicMethod() string {
	return m.privateMethod()
}

func (m *myStruct) SetData(s string) {
	m.data = s
}

// Define a private method that is not exposed through the interface
func (m *myStruct) privateMethod() string {
	if m.data != "" {
		return m.data
	}
	return "Empty"
}

// Define a function to create an instance of the struct and return it as the interface
func NewExposer(data string) Exposer {
	return &myStruct{data: data}
}
