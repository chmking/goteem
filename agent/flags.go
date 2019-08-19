package agent

import "flag"

const (
	defaultAgentHost = "127.0.0.1"
	defaultAgentPort = "5558"

	defaultManagerHost = "127.0.0.1"
	defaultManagerPort = "5557"
)

var (
	agentHost string
	agentPort string

	managerHost string
	managerPort string
)

func Flags() {
	flag.StringVar(&agentHost, "host", defaultAgentHost, "the agent address")
	flag.StringVar(&agentPort, "port", defaultAgentPort, "the agent port")

	flag.StringVar(&managerHost, "manager_host", defaultManagerHost, "the manager address")
	flag.StringVar(&managerPort, "manager_port", defaultManagerPort, "the manager port")
}
