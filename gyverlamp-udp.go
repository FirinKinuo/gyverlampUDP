package gyverlampUDP

import (
	"net"
)

const (
	DefaultGLKey string = "GL"
	DefaultGroup uint16 = 1
)

// GyverLamp structure for communicating with GyverLamp via Broadcast
//
// # IP - Broadcast IP of your network,
//
// # UDPPort - Formed from the Group number and GLKey, you can use the ComputeUDPPort function to calculate
//
// # Group - The group in which the lamps are located
//
// # GLKey - Network key, you can find it in the firmware, define GL_KEY
type GyverLamp struct {
	IP      net.IP
	UDPPort uint16
	Group   uint16
	GLKey   string
}

// NewGyverLamp is constructor of GyverLamp
//
// # IP - Broadcast IP of your network,
//
// # UDPPort - Formed from the Group number and GLKey, you can use the ComputeUDPPort function to calculate
//
// # Group - The group in which the lamps are located
//
// # GLKey - Network key, you can find it in the firmware, define GL_KEY
//
// Return copy of GyverLamp
func NewGyverLamp(IP net.IP, UDPPort uint16, group uint16, GLKey string) GyverLamp {
	return GyverLamp{
		IP:      IP,
		UDPPort: UDPPort,
		Group:   group,
		GLKey:   GLKey,
	}
}

// ComputeUDPPort calculates the port for UDP Broadcast GyverLamp in the group
//
// Will return a port in the range 50001..65010
func ComputeUDPPort(GLKey string, group uint16) (port uint16) {
	port = 17 // Start port number

	// Multiply the ASCII code of each GLKey character and current port value
	for _, GLKeyRune := range GLKey {
		port *= uint16(GLKeyRune)
	}
	port = (port % 15000) + 50000 + group // Reduce to the range from 50001 to 65010

	return port
}
