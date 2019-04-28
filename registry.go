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

func (this *registry) Get(key string) interface{} {
	return this.storage[key]
}

func (this *registry) Set(key string, val interface{}) {
	this.storage[key] = val
}

//get storage map
func (this *registry) GetMap() map[string]interface{} {
	return this.storage
}