package oss

import (
	pcontract "github.com/prismgo/framework/contracts/provider"
	"github.com/prismgo/framework/filesystem"
)

// ServiceProvider installs OSS on an application's filesystem manager.
type ServiceProvider struct{}

// Name returns the provider's stable lifecycle identity.
func (ServiceProvider) Name() string { return "prismgo.extension.oss" }

// Register declares no eager bindings because OSS is installed during Boot.
func (ServiceProvider) Register(pcontract.Application) error { return nil }

// Boot registers the OSS driver without creating an OSS client or bucket.
func (ServiceProvider) Boot(app pcontract.Application) error {
	manager, err := filesystem.ManagerFrom(app.Container())
	if err != nil {
		return err
	}
	manager.Extend("oss", func(ctx filesystem.DriverFactoryContext) (filesystem.Driver, error) {
		return NewDriver(ctx.Config.OSS)
	})
	return nil
}

var _ pcontract.ServiceProvider = ServiceProvider{}
var _ pcontract.NamedProvider = ServiceProvider{}
