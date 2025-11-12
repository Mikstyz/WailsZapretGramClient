package Tools

import (
	"fmt"
	"net"
)

func Ping(ip string, port string) error {
	_, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ip, port))

	if err != nil {
		return err
	}

	return nil
}
