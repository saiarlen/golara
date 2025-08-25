package container

import (
	"fmt"
	"reflect"
	"sync"
)

// Container provides dependency injection functionality
type Container struct {
	bindings  map[string]binding
	instances map[string]interface{}
	mutex     sync.RWMutex
}

type binding struct {
	factory   interface{}
	singleton bool
	instance  interface{}
}

// NewContainer creates a new service container
func NewContainer() *Container {
	return &Container{
		bindings:  make(map[string]binding),
		instances: make(map[string]interface{}),
	}
}

// Bind registers a binding in the container
func (c *Container) Bind(name string, factory interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.bindings[name] = binding{
		factory:   factory,
		singleton: false,
	}
}

// Singleton registers a singleton binding
func (c *Container) Singleton(name string, factory interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.bindings[name] = binding{
		factory:   factory,
		singleton: true,
	}
}

// Instance registers an existing instance
func (c *Container) Instance(name string, instance interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.instances[name] = instance
}

// Make resolves a binding from the container
func (c *Container) Make(name string) (interface{}, error) {
	c.mutex.RLock()
	
	// Check for existing instance
	if instance, exists := c.instances[name]; exists {
		c.mutex.RUnlock()
		return instance, nil
	}
	
	// Check for binding
	binding, exists := c.bindings[name]
	if !exists {
		c.mutex.RUnlock()
		return nil, fmt.Errorf("binding not found: %s", name)
	}
	
	c.mutex.RUnlock()
	
	// Create instance
	instance, err := c.createInstance(binding)
	if err != nil {
		return nil, err
	}
	
	// Store singleton instance
	if binding.singleton {
		c.mutex.Lock()
		c.instances[name] = instance
		c.mutex.Unlock()
	}
	
	return instance, nil
}

// MustMake resolves a binding or panics
func (c *Container) MustMake(name string) interface{} {
	instance, err := c.Make(name)
	if err != nil {
		panic(err)
	}
	return instance
}

// Call invokes a function with dependency injection
func (c *Container) Call(fn interface{}, args ...interface{}) ([]reflect.Value, error) {
	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()
	
	if fnType.Kind() != reflect.Func {
		return nil, fmt.Errorf("not a function")
	}
	
	// Prepare arguments
	var callArgs []reflect.Value
	
	// Add provided arguments first
	for _, arg := range args {
		callArgs = append(callArgs, reflect.ValueOf(arg))
	}
	
	// Resolve remaining parameters from container
	for i := len(args); i < fnType.NumIn(); i++ {
		paramType := fnType.In(i)
		
		// Try to resolve by type name
		instance, err := c.Make(paramType.String())
		if err != nil {
			// Try to resolve by interface name if it's an interface
			if paramType.Kind() == reflect.Interface {
				instance, err = c.Make(paramType.Name())
			}
			
			if err != nil {
				return nil, fmt.Errorf("cannot resolve parameter %d (%s): %w", i, paramType.String(), err)
			}
		}
		
		callArgs = append(callArgs, reflect.ValueOf(instance))
	}
	
	// Call function
	return fnValue.Call(callArgs), nil
}

func (c *Container) createInstance(b binding) (interface{}, error) {
	factoryValue := reflect.ValueOf(b.factory)
	factoryType := factoryValue.Type()
	
	switch factoryType.Kind() {
	case reflect.Func:
		// Call factory function with dependency injection
		results, err := c.Call(b.factory)
		if err != nil {
			return nil, err
		}
		
		if len(results) == 0 {
			return nil, fmt.Errorf("factory function returned no values")
		}
		
		// Return first result (ignore error for now)
		return results[0].Interface(), nil
		
	default:
		// Return the value directly
		return b.factory, nil
	}
}

// Bound checks if a binding exists
func (c *Container) Bound(name string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	_, exists := c.bindings[name]
	if exists {
		return true
	}
	
	_, exists = c.instances[name]
	return exists
}

// Remove removes a binding
func (c *Container) Remove(name string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	delete(c.bindings, name)
	delete(c.instances, name)
}

// Flush removes all bindings and instances
func (c *Container) Flush() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.bindings = make(map[string]binding)
	c.instances = make(map[string]interface{})
}

// GetBindings returns all binding names
func (c *Container) GetBindings() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	var names []string
	for name := range c.bindings {
		names = append(names, name)
	}
	for name := range c.instances {
		names = append(names, name)
	}
	
	return names
}

// ServiceProvider interface for service providers
type ServiceProvider interface {
	Register(container *Container)
	Boot(container *Container) error
}

// RegisterProviders registers multiple service providers
func (c *Container) RegisterProviders(providers ...ServiceProvider) error {
	// Register all providers first
	for _, provider := range providers {
		provider.Register(c)
	}
	
	// Boot all providers
	for _, provider := range providers {
		if err := provider.Boot(c); err != nil {
			return fmt.Errorf("failed to boot provider: %w", err)
		}
	}
	
	return nil
}

// Example service providers

// DatabaseServiceProvider provides database services
type DatabaseServiceProvider struct{}

func (p *DatabaseServiceProvider) Register(container *Container) {
	container.Singleton("database.manager", func() interface{} {
		// Return database manager instance
		return "database_manager_instance"
	})
}

func (p *DatabaseServiceProvider) Boot(container *Container) error {
	// Boot logic here
	return nil
}

// CacheServiceProvider provides cache services
type CacheServiceProvider struct{}

func (p *CacheServiceProvider) Register(container *Container) {
	container.Singleton("cache.manager", func() interface{} {
		// Return cache manager instance
		return "cache_manager_instance"
	})
}

func (p *CacheServiceProvider) Boot(container *Container) error {
	// Boot logic here
	return nil
}