package modules

// registry stores registered task handlers by module name.
var registry = make(map[string]TaskHandler)

// Register adds a new module handler.
func Register(name string, h TaskHandler) { registry[name] = h }

// GetHandler retrieves a handler by name.
func GetHandler(name string) (TaskHandler, bool) {
	h, ok := registry[name]
	return h, ok
}
