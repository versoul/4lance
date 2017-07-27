package delivery

import (
	"fmt"
	"versoul/4lance/socket"
)

func Deliver() {
	fmt.Println("Deliver")
	fmt.Println(socket.Data)
}
