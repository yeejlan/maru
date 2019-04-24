package maru

var (
	//global registry object, NOT thread-safe
	Registry *registry
)

type registry struct {
	storage map[string]interface{}
}

func init() {
	Registry = newRegistry()
}

func newRegistry() *registry {
	return &registry {
		storage: make(map[string]interface{}),
	}
}

func (this *registry) get(key string) interface{} {
	return this.storage[key]
}

func (this *registry) set(key string, val interface{}) {
	this.storage[key] = val
}