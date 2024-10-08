package handlers_v1

import (
	"github.com/gorilla/mux"
)

func RegisterV1Routes(router *mux.Router) {
	RegisterBookingRoutes(router)
}
