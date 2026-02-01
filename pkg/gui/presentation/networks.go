package presentation

import "github.com/warl0ck/lazycontainer/pkg/commands"

// GetNetworkDisplayStrings returns the display strings for a network
func GetNetworkDisplayStrings(network *commands.Network) []string {
	return []string{
		network.ID,
		network.Config.Mode,
		network.Status.IPv4Subnet,
	}
}
