package listeningprocesses

import (
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

func CallSockToAppendProcessList(fn sockFn) (processes []netstat.SockTabEntry, err error) {
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

func ListListeningProcesses() (processes []types.Process, err error) {
	// TODO add protocol to list of processes
	var processSockList []netstat.SockTabEntry
	processSockList, err = CallSockToAppendProcessList(netstat.UDPSocks)
	if err != nil {
		return []types.Process{}, err
	}
	processSockList, err = CallSockToAppendProcessList(netstat.UDP6Socks)
	if err != nil {
		return []types.Process{}, err
	}
	processSockList, err = CallSockToAppendProcessList(netstat.TCPSocks)
	if err != nil {
		return []types.Process{}, err
	}
	processSockList, err = CallSockToAppendProcessList(netstat.TCP6Socks)
	if err != nil {
		return []types.Process{}, err
	}

	podName := common.GetPodName()
	podNamespace := environment.GetNamespace()
	podLabels, err := environment.GetPodLabels()
	if err != nil {
		return []types.Process{}, err
	}

	for _, p := range processSockList {
		// only include if on all interfaces
		if p.LocalAddr.IP.String() != string(types.IPAllInterfaces) {
			continue
		}
		env, err := environment.GetEnvForPid(p.Process.Pid)
		if err != nil {
			return []types.Process{}, err
		}

		allowedPorts := GetPortsListFromString(string(types.EnvironmentVariableNameSharingioPairIngressReconcilerAllowedPorts))
		disabledPorts := GetPortsListFromString(string(types.EnvironmentVariableNameSharingioPairIngressReconcilerDisabledPorts))
		disabled, _ := strconv.ParseBool(env[string(types.EnvironmentVariableNameSharingioPairExposerDisabled)])

		process := types.Process{
			Name:          p.Process.Name,
			Protocol:      "TCP",
			Pid:           p.Process.Pid,
			Uid:           p.UID,
			IP:            p.LocalAddr.IP,
			Port:          p.LocalAddr.Port,
			Hostname:      env[string(types.EnvironmentVariableNameSharingioPairSetHostname)],
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
		processes = append(processes, process)
	}
	return processes, nil
}
