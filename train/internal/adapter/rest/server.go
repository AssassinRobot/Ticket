package rest

import "train/internal/ports"

func Start(api ports.APIPort, port string) error {
	fiberApp := InitRouters(api)

	err := fiberApp.Listen(port)

	return err
}
