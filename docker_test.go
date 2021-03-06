package dockertest

import (
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParsePorts(t *testing.T) {

	c := Container{}

	c.parsePorts("")
	require.NotNil(t, c.ports)
	require.Len(t, c.ports, 0)

	c.parsePorts("6379/tcp -> 0.0.0.0:32815")
	require.NotNil(t, c.ports)
	require.Len(t, c.ports, 1)
	require.Equal(t, 32815, c.ports[6379])

	c.parsePorts("6379/tcp -> 0.0.0.0:32815\n6380/udp -> 0.0.0.0:32816")
	require.NotNil(t, c.ports)
	require.Len(t, c.ports, 2)
	require.Equal(t, 32815, c.ports[6379])
	require.Equal(t, 32816, c.ports[6380])
}

func TestRun(t *testing.T) {

	con := Run("redis")
	defer con.Close()

	require.Equal(t, "127.0.0.1", con.Host())

	port := con.ports[6379]

	require.Equal(t, net.JoinHostPort("127.0.0.1", strconv.Itoa(port)), con.Addr(6379))
	require.Equal(t, port, con.Port(6379))

}

func TestWaitPort(t *testing.T) {

	con := Run("redis")
	defer con.Close()

	p := con.WaitPort(6379, 1*time.Second)
	require.NotZero(t, p)

}

func TestWaitHTTP(t *testing.T) {

	con := Run("nginx")
	defer con.Close()

	p := con.WaitHTTP(80, "/", 1*time.Second)
	require.NotZero(t, p)

}

func TestRunEnvs(t *testing.T) {

	con := RunEnvs("mysql", map[string]string{"MYSQL_ROOT_PASSWORD": "test"})
	defer con.Close()

	p := con.WaitPort(3306, 1*time.Second)
	require.NotZero(t, p)

}
