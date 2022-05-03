package object

func NewEnvironment() Environment {
	return Environment{store: make(map[string]Object)}
}

func NewEnclosedEnvironment(outer Environment) Environment {
	return Environment{store: make(map[string]Object), outer: &outer}
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return obj, ok
}

func (e Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
