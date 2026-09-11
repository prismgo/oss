# PrismGo OSS Extension

`github.com/prismgo/oss` adds the `oss` filesystem driver to one PrismGo Application.

## Installation

```bash
go get github.com/prismgo/oss
```

Register the extension between PrismGo's default providers and your application providers:

```go
import (
    "github.com/prismgo/framework/foundation"
    "github.com/prismgo/oss"
)

app := foundation.Configure().
    WithExtensionProviders(oss.ServiceProvider{}).
    WithProviders(applicationProviders...).
    Create()
```

The provider installs an application-local driver factory during `Boot`. It does not create an OSS client or bucket until the application's filesystem manager first resolves an OSS disk.

Your regular `filesystem.disks.oss` configuration continues to use `driver: "oss"`, with `bucket`, `endpoint`, `access_key`, `secret_key`, `prefix`, `url`, `visibility`, and `timeout` options.

## Direct manager use

Tests and tools that construct a filesystem manager without an Application can install the driver directly:

```go
manager, err := filesystem.NewManager(cfg)
if err != nil {
    return err
}
manager.Extend("oss", func(ctx filesystem.DriverFactoryContext) (filesystem.Driver, error) {
    return oss.NewDriver(ctx.Config.OSS)
})
```

## Tests

Unit tests use an in-process OSS-compatible HTTP server:

```bash
go test -race -coverprofile=coverage.out ./...
```
