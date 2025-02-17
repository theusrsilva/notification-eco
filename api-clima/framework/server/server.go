package server

import (
	"api-clima/app/services"
	"api-clima/domain"
	"fmt"
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
	s.Router.POST("/usuarios", s.StoreUsuario)
	s.Router.GET("/usuarios/:uid", s.ShowUsuario)
	s.Router.GET("/usuarios", s.IndexUsuario)
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
func (s *Server) StoreUsuario(c *gin.Context) {
	var req struct {
		Nome      string `json:"nome" binding:"required"`
		Sobrenome string `json:"sobrenome" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		Cidade    int    `json:"cidade" binding:"required"`
		Horario   string `json:"horario" binding:"omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usuario, err := domain.NewUsuario(req.Nome, req.Sobrenome, req.Email, req.Cidade)
	if err != nil {
		log.Println("Erro: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ""})
		return
	}
	if req.Horario == "" {
		req.Horario = "10:00"
	}
	usuarioNew, err := services.InsertUsuario(usuario, req.Horario, s.DB)

	if err != nil {
		if err.Error() == fmt.Sprintf("erro: o email '%s' já está em uso", usuario.Email) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		log.Println("Erro: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ""})
		return
	}
	c.JSON(http.StatusOK, usuarioNew)
}

func (s *Server) IndexUsuario(c *gin.Context) {
	nomeInput := c.DefaultQuery("nome", "")
	sobrenomeInput := c.DefaultQuery("sobrenome", "")
	cidadeInput := c.DefaultQuery("cidade", "")
	emailInput := c.DefaultQuery("email", "")

	usuarios, err := services.IndexUsuario(nomeInput, sobrenomeInput, cidadeInput, emailInput, s.DB)

	if err != nil {
		log.Println("Erro: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": ""})
		return
	}

	c.JSON(http.StatusOK, usuarios)
}
func (s *Server) ShowUsuario(c *gin.Context) {
	uid := c.Param("uid")
	usuario, err := services.ShowUsuario(uid, s.DB)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, usuario)
}

func (s *Server) Start(port string) {
	log.Printf("Servidor rodando na porta %s", port)
	if err := s.Router.Run(port); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
