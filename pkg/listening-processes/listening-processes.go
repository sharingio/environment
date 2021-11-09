package listeningprocesses

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/cakturk/go-netstat/netstat"

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
			env, err := environment.GetEnvForPid(s.Process.Pid)
			if err != nil {
				log.Println(err.Error())
			} else {
				log.Printf("%#v\n", env)
			}
		}
		return p
	}()
	return processes, nil
}

func ListListeningProcesses() (processes []types.Process, err error) {
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

	for _, p := range processSockList {
		env, err := environment.GetEnvForPid(p.Process.Pid)
		if err != nil {
			return []types.Process{}, err
		}

		allowedPorts := GetPortsListFromString(string(types.EnvironmentVariableNameSharingioPairIngressReconcilerAllowedPorts))
		disabledPorts := GetPortsListFromString(string(types.EnvironmentVariableNameSharingioPairIngressReconcilerDisabledPorts))

		process := types.Process{
			Name:          p.Process.Name,
			Pid:           p.Process.Pid,
			Uid:           p.UID,
			LocalAddr:     p.LocalAddr,
			Hostname:      env[string(types.EnvironmentVariableNameSharingioPairSetHostname)],
			AllowedPorts:  allowedPorts,
			DisabledPorts: disabledPorts,
		}
		processes = append(processes, process)
	}
	return processes, nil
}
