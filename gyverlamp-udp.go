package gyverlampUDP

import (
	"errors"
	"fmt"
	"net"
)

type (
	command byte // Type describing the GyverLamp protocol command, more details in xlsx: https://shorturl.at/cftwX
)

const (
	controlCommand command = 0x30
	settingCommand command = 0x31
	modeCommand    command = 0x32
	sunriseCommand command = 0x33
	paletteCommand command = 0x35
)

const (
	controlParamTurnOffLight byte = 0x30
	controlParamTurnOnLight  byte = 0x31
)

// Send sends a command to GyverLamp
//
// data - command data bytes, more about command data in xlsx: https://shorturl.at/cftwX
func (c command) Send(gl *GyverLamp, data ...byte) error {
	commandMessage := gl.getCommandMessage(c, data...)
	err := gl.sendBroadcastMessage(commandMessage)
	if err != nil {
		return fmt.Errorf("unable to send command, %w", err)
	}

	return nil
}

const (
	DefaultGroup  uint16 = 1    // Default group for lamps
	dataSeparator byte   = 0x2c // Data separator (comma character)
)

var DefaultGLKey = []byte{0x47, 0x4c} // Default GL Key for lamps

// GyverLamp structure for communicating with GyverLamp via Broadcast
// https://github.com/AlexGyver/GyverLamp2/blob/main/firmware/GyverLamp2/GyverLamp2.ino in #define GL_KEY
type GyverLamp struct {
	GLKey            []byte         // The GyverLamp lamp network key, by default DefaultGLKey key is used, the custom is written in the firmware code
	UDPAddress       *net.UDPAddr   // Broadcast IP of your network with Port, formed from the Group number and GLKey
	PacketConnection net.PacketConn // Connection for sending broadcast messages
}

// NewGyverLamp is constructor of GyverLamp
//
// IP - Broadcast IP of your network,
// UDPPort - Formed from the Group number and GLKey, you can use the ComputeUDPPort function to calculate,
// GLKey - The GyverLamp lamp network key, by default DefaultGLKey key is used, the custom is written in the firmware code
// https://github.com/AlexGyver/GyverLamp2/blob/main/firmware/GyverLamp2/GyverLamp2.ino in #define GL_KEY.
func NewGyverLamp(IP net.IP, UDPPort uint16, GLKey []byte) (gyverLamp *GyverLamp) {
	UDPAddress := net.UDPAddr{
		IP:   IP,
		Port: int(UDPPort),
	}
	gyverLamp = &GyverLamp{
		UDPAddress: &UDPAddress,
		GLKey:      GLKey,
	}

	return gyverLamp
}

// ComputeUDPPort calculates the port for UDP Broadcast GyverLamp in the group
//
// Will return a port in the range 50001..65010
func ComputeUDPPort(GLKey []byte, group uint16) (port uint16) {
	port = 17 // Start port number

	// Multiply the ASCII code of each GLKey character and current port value
	for _, GLKeyRune := range GLKey {
		port *= uint16(GLKeyRune)
	}
	port = (port % 15000) + 50000 + group // Reduce to the range from 50001 to 65010

	return port
}

// CreateNewConnection creates a Broadcast connection and adds to the GyverLamp PacketConnection
func (gyverLamp *GyverLamp) CreateNewConnection() error {
	conn, err := net.ListenPacket("udp4", fmt.Sprintf(":%d", gyverLamp.UDPAddress.Port))
	if err != nil {
		return fmt.Errorf("unable to create connection to Gyver Lamp %s, %w", gyverLamp.UDPAddress.String(), err)
	}

	gyverLamp.PacketConnection = conn

	return nil
}

func (gyverLamp *GyverLamp) sendBroadcastMessage(message []byte) error {
	if gyverLamp.PacketConnection == nil {
		return errors.New("before sending the command, you need to create a connection by CreateNewConnection()")
	}
	_, err := gyverLamp.PacketConnection.WriteTo(message, gyverLamp.UDPAddress)
	if err != nil {
		return err
	}

	return nil
}

func (gyverLamp *GyverLamp) getCommandMessage(commandType command, data ...byte) []byte {
	var commandMessage []byte

	commandMessage = append(commandMessage, gyverLamp.GLKey...)
	commandMessage = append(commandMessage, dataSeparator, byte(commandType))

	// Adds data bytes to the commandMessage slice, delimited by dataSeparator
	for _, dataByte := range data {
		commandMessage = append(commandMessage, dataSeparator, dataByte)
	}

	return commandMessage
}

// TurnOff sends a command to turn off the light
func (gyverLamp *GyverLamp) TurnOff() error {
	return controlCommand.Send(gyverLamp, controlParamTurnOffLight)
}

// TurnOn sends a command to turn on the light
func (gyverLamp *GyverLamp) TurnOn() error {
	return controlCommand.Send(gyverLamp, controlParamTurnOnLight)
}
