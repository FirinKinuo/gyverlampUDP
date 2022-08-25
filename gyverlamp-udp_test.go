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
		GLKey   string
	}{
		IP:      net.IPv4(10, 0, 0, 255),
		UDPPort: 61197,
		Group:   DefaultGroup,
		GLKey:   DefaultGLKey,
	}

	expectGyverLamp := GyverLamp{
		IP:      payload.IP,
		UDPPort: payload.UDPPort,
		Group:   payload.Group,
		GLKey:   payload.GLKey,
	}

	actualGyverLamp := NewGyverLamp(payload.IP, payload.UDPPort, payload.Group, payload.GLKey)

	require.Equal(t, expectGyverLamp, actualGyverLamp)
}

func TestComputeUDPPort(t *testing.T) {
	testCases := []struct {
		GLKey  string
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
