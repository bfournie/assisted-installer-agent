package util

import (
	"net"
	"os"
	"path/filepath"

	"github.com/jaypipes/ghw"
	"github.com/openshift/assisted-installer-agent/src/config"
	"github.com/vishvananda/netlink"
)

//go:generate mockery --name IDependencies --inpackage
type IDependencies interface {
	Execute(command string, args ...string) (stdout string, stderr string, exitCode int)
	ExecutePrivileged(command string, args ...string) (stdout string, stderr string, exitCode int)
	ReadFile(fname string) ([]byte, error)
	Stat(fname string) (os.FileInfo, error)
	Hostname() (string, error)
	Interfaces() ([]Interface, error)
	Block(opts ...*ghw.WithOption) (*ghw.BlockInfo, error)
	Product(opts ...*ghw.WithOption) (*ghw.ProductInfo, error)
	ReadDir(dirname string) ([]os.FileInfo, error)
	Abs(path string) (string, error)
	EvalSymlinks(path string) (string, error)
	LinkByName(name string) (netlink.Link, error)
	RouteList(link netlink.Link, family int) ([]netlink.Route, error)
	GPU(opts ...*ghw.WithOption) (*ghw.GPUInfo, error)
	Memory(opts ...*ghw.WithOption) (*ghw.MemoryInfo, error)
	Chassis(opts ...*ghw.WithOption) (*ghw.ChassisInfo, error)
	PCI(opts ...*ghw.WithOption) (*ghw.PCIInfo, error)
	GetGhwChrootRoot() string
}

type Dependencies struct {
	NetlinkRouteFinder
	GhwChrootRoot string
	dryRunConfig  *config.DryRunConfig
}

func (d *Dependencies) GetGhwChrootRoot() string {
	return d.GhwChrootRoot
}

func (d *Dependencies) Execute(command string, args ...string) (stdout string, stderr string, exitCode int) {
	return Execute(command, args...)
}

func (d *Dependencies) ExecutePrivileged(command string, args ...string) (stdout string, stderr string, exitCode int) {
	return ExecutePrivileged(command, args...)
}

func (d *Dependencies) ReadFile(fname string) ([]byte, error) {
	return os.ReadFile(fname)
}

func (d *Dependencies) Stat(fname string) (os.FileInfo, error) {
	return os.Stat(fname)
}

func (d *Dependencies) Hostname() (string, error) {
	if d.dryRunConfig.DryRunEnabled {
		return d.dryRunConfig.ForcedHostname, nil
	}

	return os.Hostname()
}

func (d *Dependencies) Interfaces() ([]Interface, error) {
	ins, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	ret := make([]Interface, 0)
	for _, in := range ins {
		ret = append(ret, &NetworkInterface{netInterface: in, dependencies: d})
	}
	return ret, nil
}

func (d *Dependencies) Block(opts ...*ghw.WithOption) (*ghw.BlockInfo, error) {
	return ghw.Block(opts...)
}

func (d *Dependencies) Product(opts ...*ghw.WithOption) (*ghw.ProductInfo, error) {
	return ghw.Product(opts...)
}

func (d *Dependencies) ReadDir(dirname string) ([]os.FileInfo, error) {
	dirEntry, err := os.ReadDir(dirname)
	if dirEntry != nil {
		var filesInfo []os.FileInfo
		for _, entry := range dirEntry {
			info, _ := entry.Info()
			filesInfo = append(filesInfo, info)
		}
		return filesInfo, err
	}
	return nil, err
}

func (d *Dependencies) Abs(path string) (string, error) {
	return filepath.Abs(path)
}

func (d *Dependencies) EvalSymlinks(path string) (string, error) {
	return filepath.EvalSymlinks(path)
}

func (d *Dependencies) GPU(opts ...*ghw.WithOption) (*ghw.GPUInfo, error) {
	return ghw.GPU(opts...)
}

func (d *Dependencies) Memory(opts ...*ghw.WithOption) (*ghw.MemoryInfo, error) {
	return ghw.Memory(opts...)
}

func (d *Dependencies) Chassis(opts ...*ghw.WithOption) (*ghw.ChassisInfo, error) {
	return ghw.Chassis(opts...)
}

func (d *Dependencies) PCI(opts ...*ghw.WithOption) (*ghw.PCIInfo, error) {
	return ghw.PCI(opts...)
}

func NewDependencies(dryRunConfig *config.DryRunConfig, ghwChrootRoot string) IDependencies {
	return &Dependencies{
		GhwChrootRoot: ghwChrootRoot,
		dryRunConfig:  dryRunConfig,
	}
}
