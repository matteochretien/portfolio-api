package dig

import "go.uber.org/dig"

type Provider struct {
	// Constructor is the Constructor function that will be called to create the provider.
	Constructor    any
	ProvideOptions []dig.ProvideOption
}

func NewProvider(constructor any, options ...dig.ProvideOption) Provider {
	p := Provider{
		Constructor:    constructor,
		ProvideOptions: options,
	}

	return p
}
