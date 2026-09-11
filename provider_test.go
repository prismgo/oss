package oss

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/prismgo/framework/filesystem"
	"github.com/prismgo/framework/foundation"
)

func TestNewDriverConstructsLazily(t *testing.T) {
	driver, err := NewDriver(filesystem.OSSConfig{
		Bucket:   "example",
		Endpoint: "http://127.0.0.1:1",
	})
	if err != nil {
		t.Fatalf("NewDriver() error = %v, want nil", err)
	}
	if driver == nil {
		t.Fatal("NewDriver() driver = nil, want non-nil")
	}
	if err := driver.Close(); err != nil {
		t.Fatalf("Close() error = %v, want nil", err)
	}
}

func TestServiceProviderRegistersLazyApplicationLocalDriver(t *testing.T) {
	manager := newTestManager(t)
	otherManager := newTestManager(t)
	otherErr := errors.New("other application driver")
	otherManager.Extend("oss", func(filesystem.DriverFactoryContext) (filesystem.Driver, error) {
		return nil, otherErr
	})
	app := foundation.NewApplication()
	t.Cleanup(func() { _ = app.Close() })
	if err := app.Instance("filesystem.manager", manager); err != nil {
		t.Fatalf("bind filesystem manager: %v", err)
	}

	provider := ServiceProvider{}
	if got, want := provider.Name(), "prismgo.extension.oss"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
	if err := provider.Register(app); err != nil {
		t.Fatalf("Register() error = %v, want nil", err)
	}
	if err := provider.Boot(app); err != nil {
		t.Fatalf("Boot() error = %v, want nil", err)
	}

	err := manager.Default().Put(context.Background(), "probe.txt", "probe")
	if err == nil || !strings.Contains(err.Error(), "oss bucket or endpoint is empty") {
		t.Fatalf("extension Put() error = %v, want lazy OSS configuration error", err)
	}
	err = otherManager.Default().Put(context.Background(), "probe.txt", "probe")
	if !errors.Is(err, otherErr) {
		t.Fatalf("other application Put() error = %v, want %v", err, otherErr)
	}
}

func newTestManager(t *testing.T) *filesystem.Manager {
	t.Helper()
	manager, err := filesystem.NewManager(filesystem.Config{
		Default: "cloud",
		Disks: map[string]filesystem.DiskConfig{
			"cloud": {Driver: "oss"},
		},
	})
	if err != nil {
		t.Fatalf("NewManager() error = %v, want nil", err)
	}
	t.Cleanup(func() { _ = manager.Close() })
	return manager
}
