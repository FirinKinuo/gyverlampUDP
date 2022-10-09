package gyverlampUDP

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func TestNewGyverLamp(t *testing.T) {
	payload := struct {
		IP      net.IP
		UDPPort uint16
		Group   uint16
		GLKey   []byte
	}{
		IP:      net.IPv4(127, 0, 0, 1),
		UDPPort: 61197,
		Group:   DefaultGroup,
		GLKey:   DefaultGLKey,
	}

	expectGyverLamp := &GyverLamp{
		UDPAddress: &net.UDPAddr{
			IP:   net.IPv4(127, 0, 0, 1),
			Port: 61197,
		},
		GLKey: DefaultGLKey,
	}

	actualGyverLamp := NewGyverLamp(payload.IP, payload.UDPPort, DefaultGLKey)
	require.Equal(t, expectGyverLamp, actualGyverLamp)
}

func TestComputeUDPPort(t *testing.T) {
	testCases := []struct {
		GLKey  []byte
		Group  uint16
		Expect uint16
	}{
		{
			GLKey:  DefaultGLKey,
			Group:  DefaultGroup,
			Expect: 61197,
		},
		{
			GLKey:  DefaultGLKey,
			Group:  7,
			Expect: 61203,
		},
	}

	for _, testCase := range testCases {
		actualUDPPort := ComputeUDPPort(testCase.GLKey, testCase.Group)

		assert.Equal(
			t,
			testCase.Expect,
			actualUDPPort,
			fmt.Sprintf(
				"The UDP port was not correctly calculated for GLKey: %s, Group: %d",
				testCase.GLKey,
				testCase.Group))
	}
}

func TestGyverLamp_CreateNewConnection(t *testing.T) {
	gl := GyverLamp{UDPAddress: &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 1}}

	err := gl.CreateNewConnection()
	require.NoError(t, err)
	require.NotNil(t, gl.PacketConnection)
}
