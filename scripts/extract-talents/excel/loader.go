package excel

import "fmt"

var resources = map[string]any{}

func load(name string, v any) {
	if _, ok := resources[name]; ok {
		panic(fmt.Errorf("%s is already loaded", name))
	}
	resources[name] = v
}

func LoadResources(loader func(name string, v any) error) error {
	for name, v := range resources {
		if err := loader(name, v); err != nil {
			return fmt.Errorf("failed to load %s: %w", name, err)
		}
	}
	return nil
}
