package gyverlampUDP

import (
	"fmt"
	"net"
)

const (
	DefaultGLKey string = "GL"
	DefaultGroup uint16 = 1
)

// GyverLamp describes the control methods of GyverLamp
type GyverLamp interface {
	CreateNewConnection() error
	TurnOn() error
	TurnOff() error
}

// GyverLampImpl structure for communicating with GyverLamp via Broadcast
//
// UDPAddress - Broadcast IP of your network with Port, formed from the Group number and GLKey
// PacketConnection - Connection for sending broadcast messages
type GyverLampImpl struct {
	UDPAddress       *net.UDPAddr
	PacketConnection net.PacketConn
}

// NewGyverLamp is constructor of GyverLamp
//
// IP - Broadcast IP of your network,
// UDPPort - Formed from the Group number and GLKey, you can use the ComputeUDPPort function to calculate
func NewGyverLamp(IP net.IP, UDPPort uint16) (gyverLamp GyverLamp) {
	UDPAddress := net.UDPAddr{
		IP:   IP,
		Port: int(UDPPort),
	}
	gyverLamp = &GyverLampImpl{
		UDPAddress: &UDPAddress,
	}

	return gyverLamp
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

// CreateNewConnection creates a Broadcast connection and adds to the GyverLampImpl PacketConnection
func (gyverLamp *GyverLampImpl) CreateNewConnection() error {
	conn, err := net.ListenPacket("udp4", fmt.Sprintf(":%d", gyverLamp.UDPAddress.Port))
	if err != nil {
		return fmt.Errorf("unable to create connection to Gyver Lamp %s, %w", gyverLamp.UDPAddress.String(), err)
	}

	gyverLamp.PacketConnection = conn

	return nil
}

func (gyverLamp *GyverLampImpl) sendBroadcastMessage(message []byte) error {
	_, err := gyverLamp.PacketConnection.WriteTo(message, gyverLamp.UDPAddress)
	if err != nil {
		return err
	}

	return nil
}

// TurnOff sends a command to turn off the light
func (gyverLamp *GyverLampImpl) TurnOff() error {
	commandTurnOff := []byte{0x47, 0x4c, 0x2c, 0x30, 0x2c, 0x30}

	err := gyverLamp.sendBroadcastMessage(commandTurnOff)
	if err != nil {
		return err
	}

	return nil
}

// TurnOn sends a command to turn on the light
func (gyverLamp *GyverLampImpl) TurnOn() error {
	commandTurnOn := []byte{0x47, 0x4c, 0x2c, 0x30, 0x2c, 0x31}

	err := gyverLamp.sendBroadcastMessage(commandTurnOn)
	if err != nil {
		return err
	}

	return nil
}
