package core

import (
	"fmt"
	"os"
	"path"
	"plugin"

	"github.com/lburgazzoli/camel-go/camel"
)

// ==========================
//
// DefaultRegistryLoader
//
// ==========================

// DefaultRegistryLoader --
type DefaultRegistryLoader struct {
	DefaultService
}

// Order --
func (loader *DefaultRegistryLoader) Order() int {
	return loader.order
}

// ==========================
//
// PluginRegistryLoader
//
//     Use Go's plugins to load objects
//
// ==========================

// NewPluginRegistryLoader --
func NewPluginRegistryLoader(searchPath string) camel.RegistryLoader {
	return &pluginRegistryLoader{
		DefaultRegistryLoader: DefaultRegistryLoader{
			DefaultService: DefaultService{
				order:  0,
				status: camel.ServiceStatusSTOPPED,
			},
		},
		searchPath: searchPath,
	}
}

type pluginRegistryLoader struct {
	DefaultRegistryLoader

	searchPath string
}

// Start --
func (loader *pluginRegistryLoader) Start() {
	// maybe here we should scan the search path to pre instantiate objects
}

// Stop --
func (loader *pluginRegistryLoader) Stop() {
}

// GetComponent --
func (loader *pluginRegistryLoader) Load(name string) (interface{}, error) {
	pluginPath := path.Join(loader.searchPath, fmt.Sprintf("%s.so", name))
	_, err := os.Stat(pluginPath)

	if os.IsNotExist(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	plug, err := plugin.Open(pluginPath)
	if err != nil {
		fmt.Printf("failed to open plugin %s: %v\n", name, err)
		return nil, err
	}

	symbol, err := plug.Lookup("Create")
	if err != nil {
		fmt.Printf("plugin %s does not export symbol \"Create\"\n", name)
		return nil, err
	}

	// Load the object from
	result := symbol.(func() interface{})()

	return result, nil
}
