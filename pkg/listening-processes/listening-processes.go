package listeningprocesses

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cakturk/go-netstat/netstat"

	"github.com/sharingio/environment/pkg/common"
	"github.com/sharingio/environment/pkg/environment"
	"github.com/sharingio/environment/pkg/types"
)

type sockFn = func(netstat.AcceptFn) ([]netstat.SockTabEntry, error)

func GetPortsListFromString(portsString string) (ports []int) {
	portsStrings := strings.Split(portsString, ",")
	for _, p := range portsStrings {
		port, _ := strconv.Atoi(p)
		if port == 0 {
			continue
		}
		ports = append(ports, port)
	}
	return ports
}

func GetProcessFromSockFn(fn sockFn) (processes []netstat.SockTabEntry, err error) {
	processes, err = fn(func(s *netstat.SockTabEntry) bool {
		return s.State == netstat.Listen
	})
	if err != nil {
		return processes, err
	}
	processes = func() (p []netstat.SockTabEntry) {
		for _, s := range processes {
			if s.Process.Pid == os.Getpid() {
				continue
			}
			p = append(p, s)
		}
		return p
	}()
	return processes, nil
}

func NewProcessForSockTabEntry(sock netstat.SockTabEntry, overrides types.Process) (process types.Process, err error) {
	podName := common.GetPodName()
	podNamespace := environment.GetNamespace()
	podLabels, err := environment.GetPodLabels()
	if err != nil {
		return types.Process{}, err
	}
	// only include if on all interfaces
	if !(sock.LocalAddr.IP.String() == string(types.IPAllInterfaces) ||
		sock.LocalAddr.IP.String() == string(types.IPv6AllInterfaces)) {
		return types.Process{}, fmt.Errorf("Error: IP (%v) must match all interfaces", sock.LocalAddr.IP.String())
	}
	env, err := environment.GetEnvForPid(sock.Process.Pid)
	if err != nil {
		return types.Process{}, err
	}

	allowedPorts := GetPortsListFromString(string(types.EnvVarNameSharingioPairIngressReconcilerAllowedPorts))
	disabledPorts := GetPortsListFromString(string(types.EnvVarNameSharingioPairIngressReconcilerDisabledPorts))
	disabled, _ := strconv.ParseBool(env[string(types.EnvVarNameSharingioPairExposerDisabled)])

	process = types.Process{
		Name:          sock.Process.Name,
		Protocol:      overrides.Protocol,
		Pid:           sock.Process.Pid,
		Uid:           sock.UID,
		IP:            sock.LocalAddr.IP,
		Port:          sock.LocalAddr.Port,
		Hostname:      env[string(types.EnvVarNameSharingioPairSetHostname)],
		AllowedPorts:  allowedPorts,
		DisabledPorts: disabledPorts,
		Disabled:      disabled,
		PodName:       podName,
		PodNamespace:  podNamespace,
		PodLabels:     podLabels,
		ServiceName:   "",
		ServicePort:   0,
		ExternalIP:    "",
		IngressHost:   "",
	}
	return process, nil
}

type netstatSock struct {
	fn       func(netstat.AcceptFn) ([]netstat.SockTabEntry, error)
	protocol string
}

func ListListeningProcesses() (processes []types.Process, err error) {
	// TODO UDP filter for UNCONN as well as LISTEN
	var processSockList []netstat.SockTabEntry
	socks := []netstatSock{
		{
			fn:       netstat.UDPSocks,
			protocol: "UDP",
		},
		{
			fn:       netstat.UDP6Socks,
			protocol: "UDP",
		},
		{
			fn:       netstat.TCPSocks,
			protocol: "TCP",
		},
		{
			fn:       netstat.TCP6Socks,
			protocol: "TCP",
		},
	}
	for _, s := range socks {
		processSockList, err = GetProcessFromSockFn(s.fn)
		if err != nil {
			return []types.Process{}, err
		}
		for _, p := range processSockList {
			process, err := NewProcessForSockTabEntry(p, types.Process{Protocol: s.protocol})
			if err != nil {
				return []types.Process{}, err
			}
			processes = append(processes, process)
		}
	}
	return processes, nil
}
