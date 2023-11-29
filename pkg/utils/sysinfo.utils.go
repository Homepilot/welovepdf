package utils

import (
	"github.com/elastic/go-sysinfo"
)

func getSysInfo() map[string]any {
	goInfo := sysinfo.Go()
	host, _ := sysinfo.Host()
	hostInfo := host.Info()
	osInfo := hostInfo.OS

	return map[string]any{
		"sysinfo": map[string]any{
			"OS":       goInfo.OS,
			"Arch":     goInfo.Arch,
			"MaxProcs": goInfo.MaxProcs,
			"Version":  goInfo.Version,
		},
		"hostinfo": map[string]any{
			"Architecture":  hostInfo.Architecture,
			"Hostname":      hostInfo.Hostname,
			"KernelVersion": hostInfo.KernelVersion,
			"UniqueID":      hostInfo.UniqueID,
		},
		"osinfo": map[string]any{
			"Type":     osInfo.Type,
			"Family":   osInfo.Family,
			"Platform": osInfo.Platform,
			"Name":     osInfo.Name,
			"Version":  osInfo.Version,
			"Major":    osInfo.Major,
			"Minor":    osInfo.Minor,
			"Patch":    osInfo.Patch,
			"Build":    osInfo.Build,
		},
	}
}
