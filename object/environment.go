package object

func NewEnvironment() Environment {
	return Environment{Store: make(map[string]Object)}
}

type Environment struct {
	Store map[string]Object
}

func (e Environment) Get(name string) (Object, bool) {
	obj, ok := e.Store[name]
	return obj, ok
}

func (e Environment) Set(name string, val Object) Object {
	e.Store[name] = val
	return val
}
