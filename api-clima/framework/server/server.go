package server

import (
	"api-clima/app/services"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

type Server struct {
	Router *gin.Engine
	DB     *gorm.DB
}

func NewServer(db *gorm.DB) *Server {
	server := &Server{
		Router: gin.Default(),
		DB:     db,
	}

	server.SetupRoutes()
	return server
}

func (s *Server) SetupRoutes() {
	s.Router.GET("/cidades", s.Cidades)
}
func (s *Server) Cidades(c *gin.Context) {
	cidadeInput := c.DefaultQuery("nome", "")

	cidades, err := services.GetCidades(cidadeInput)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var cidadesResponse []services.CidadeResponse
	for _, cidade := range cidades {
		cidadesResponse = append(cidadesResponse, services.CidadeResponse{
			ID:   cidade.ID,
			Nome: cidade.Nome,
			UF:   cidade.UF,
		})
	}
	c.JSON(http.StatusOK, cidadesResponse)
}

func (s *Server) Start(port string) {
	log.Printf("Servidor rodando na porta %s", port)
	if err := s.Router.Run(port); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
