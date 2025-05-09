package addr

import (
	"net"
)

func GetFreeport() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, nil
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// func main() {
// 	port, _ := GetFreeport()
// 	fmt.Println(port)
// }
