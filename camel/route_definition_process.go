package camel

import (
	"fmt"
)

// ==========================
//
// Extend RouteDefinition DSL
//
// ==========================

// Process --
func (definition *RouteDefinition) Process() *ProcessDefinition {
	d := ProcessDefinition{
		parent:   definition,
		children: nil,
	}

	definition.AddChild(&d)

	return &d
}

// ==========================
//
// ProcessDefinition
//
// ==========================

// ProcessDefinition --
type ProcessDefinition struct {
	parent   *RouteDefinition
	children []Definition

	processor    func(*Exchange)
	processorRef string
}

// Parent --
func (definition *ProcessDefinition) Parent() Definition {
	return definition.parent
}

// Children --
func (definition *ProcessDefinition) Children() []Definition {
	return definition.children
}

// Unwrap ---
func (definition *ProcessDefinition) Unwrap(context *Context, parent Processor) (Processor, Service, error) {
	if definition.processor != nil {
		p := NewProcessorWithParent(parent, func(e *Exchange, out chan<- *Exchange) {
			definition.processor(e)

			out <- e
		})

		return p, nil, nil
	}

	if definition.processorRef != "" {
		registry := context.Registry()
		ifc, err := registry.Lookup(definition.processorRef)

		if ifc != nil && err == nil {
			if processor, ok := ifc.(func(e *Exchange)); ok {
				p := NewProcessorWithParent(parent, func(e *Exchange, out chan<- *Exchange) {
					processor(e)

					out <- e
				})

				return p, nil, nil
			}
		}

		if err == nil {
			err = fmt.Errorf("Unsupported type for ref:%s, type=%T", definition.processorRef, ifc)
		}

		// TODO: error handling
		return nil, nil, err
	}

	return nil, nil, nil

}

// Fn --
func (definition *ProcessDefinition) Fn(processor func(*Exchange)) *RouteDefinition {
	definition.processor = processor
	return definition.parent
}

// Ref --
func (definition *ProcessDefinition) Ref(ref string) *RouteDefinition {
	definition.processorRef = ref
	return definition.parent
}
