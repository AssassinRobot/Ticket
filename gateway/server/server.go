package server

import (
	"gateway/events"
	router "gateway/server/routes"
)

func Start(requestHandler events.RequestSender, port string) error {
	fiberApp := router.InitRouters(requestHandler)

	err := fiberApp.Listen(":"+port)

	return err
}
