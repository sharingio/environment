package listeningprocesses

import (
	"fmt"
	"log"
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
	portsStrings := strings.Split(portsString, " ")
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
		return types.Process{}, fmt.Errorf("Warning: IP (%v) must match all interfaces", sock.LocalAddr.IP.String())
	}
	env, err := environment.GetEnvForPid(sock.Process.Pid)
	if err != nil {
		return types.Process{}, err
	}

	// use new env var or deprecated
	allowedPorts := GetPortsListFromString(env[string(types.EnvVarNameSharingioPairExposerAllowedPorts)])
	if len(allowedPorts) == 0 {
		allowedPorts = GetPortsListFromString(env[string(types.EnvVarNameSharingioPairIngressReconcilerAllowedPorts)])
	}
	disabledPorts := GetPortsListFromString(env[string(types.EnvVarNameSharingioPairExposerDisabledPorts)])
	if len(disabledPorts) == 0 {
		disabledPorts = GetPortsListFromString(env[string(types.EnvVarNameSharingioPairIngressReconcilerDisabledPorts)])
	}
	hostname := env[string(types.EnvVarNameSharingioPairExposerHostname)]
	if hostname == "" {
		hostname = env[string(types.EnvVarNameSharingioPairSetHostname)]
	}
	disabled, _ := strconv.ParseBool(env[string(types.EnvVarNameSharingioPairExposerDisabled)])

	process = types.Process{
		Name:          sock.Process.Name,
		Protocol:      overrides.Protocol,
		Pid:           sock.Process.Pid,
		Uid:           sock.UID,
		IP:            sock.LocalAddr.IP,
		Port:          sock.LocalAddr.Port,
		Hostname:      hostname,
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
	protocol types.Protocol
}

func ListListeningProcesses() (processes []types.Process, err error) {
	// TODO UDP filter for UNCONN as well as LISTEN
	var processSockList []netstat.SockTabEntry
	socks := []netstatSock{
		{
			fn:       netstat.UDPSocks,
			protocol: types.ProtocolUDP,
		},
		{
			fn:       netstat.UDP6Socks,
			protocol: types.ProtocolUDP,
		},
		{
			fn:       netstat.TCPSocks,
			protocol: types.ProtocolTCP,
		},
		{
			fn:       netstat.TCP6Socks,
			protocol: types.ProtocolTCP,
		},
	}
	for _, s := range socks {
		processSockList, err = GetProcessFromSockFn(s.fn)
		if err != nil {
			log.Println(err)
			continue
		}
	process:
		for _, p := range processSockList {
			process, err := NewProcessForSockTabEntry(p, types.Process{Protocol: s.protocol})
			if err != nil {
				log.Println(err)
				continue process
			}
			if len(process.DisabledPorts) > 0 && common.IntInListOfInts(process.DisabledPorts, int(process.Port)) == true {
				log.Printf("Ignoring port '%v' as it's in disabled ports list '%#v' for process '%v'(%v)", process.Port, process.DisabledPorts, process.Name, process.Pid)
				continue process
			}
			if len(process.AllowedPorts) > 0 && common.IntInListOfInts(process.AllowedPorts, int(process.Port)) == false {
				log.Printf("Ignoring port '%v' as it's not in allowed ports list '%#v' for process '%v'(%v)", process.Port, process.DisabledPorts, process.Name, process.Pid)
				continue process
			}
			processes = append(processes, process)
		}
	}
	return processes, nil
}
