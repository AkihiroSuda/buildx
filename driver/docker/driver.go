package docker

import (
	"context"
	"net"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/moby/buildkit/client"
	"github.com/pkg/errors"
	"github.com/tonistiigi/buildx/driver"
	"github.com/tonistiigi/buildx/util/progress"
)

type Driver struct {
	factory driver.Factory
	driver.InitConfig
	version dockertypes.Version
}

func (d *Driver) Bootstrap(ctx context.Context, l progress.Logger) error {
	return nil
}

func (d *Driver) Info(ctx context.Context) (*driver.Info, error) {
	_, err := d.DockerAPI.ServerVersion(ctx)
	if err != nil {
		return nil, errors.Wrapf(driver.ErrNotConnecting, err.Error())
	}
	return &driver.Info{
		Status: driver.Running,
	}, nil
}

func (d *Driver) Stop(ctx context.Context, force bool) error {
	return nil
}

func (d *Driver) Rm(ctx context.Context, force bool) error {
	return nil
}

func (d *Driver) Client(ctx context.Context) (*client.Client, error) {
	return client.New(ctx, "", client.WithDialer(func(string, time.Duration) (net.Conn, error) {
		return d.DockerAPI.DialHijack(ctx, "/grpc", "h2c", nil)
	}))
}

func (d *Driver) Features() map[driver.Feature]bool {
	return map[driver.Feature]bool{
		driver.OCIExporter:    false,
		driver.DockerExporter: false,

		driver.CacheExport:   false,
		driver.MultiPlatform: false,
	}
}

func (d *Driver) Factory() driver.Factory {
	return d.factory
}

func (d *Driver) IsDefaultMobyDriver() {}
