package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/razorpay/bankProject/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	server.router = router
	return server

}
func errResponse(err error) map[string]interface{} {
	return map[string]interface{}{
		"error": err.Error(),
	}
}
func (s *Server) StartServer(address string) {
	s.router.Run(address)
}
