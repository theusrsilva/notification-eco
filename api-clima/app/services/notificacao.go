package services

import (
	"api-clima/domain"
	"github.com/jinzhu/gorm"
	"log"
	"sync"
	"time"
)

type NotificacaoService struct {
	DB *gorm.DB
}

func NewNotificacaoService(db *gorm.DB) *NotificacaoService {
	return &NotificacaoService{DB: db}
}
func (ns *NotificacaoService) ProcessaNotificacoes(horario string) {

	log.Printf("Checando notificações para %s", horario)

	notificacoes, err := ns.BuscarNotificacoes(horario)
	if err != nil {
		log.Printf("Erro ao buscar notificações: %v", err)
		return
	}
	var wg sync.WaitGroup
	for _, n := range notificacoes {
		wg.Add(1)
		go ns.EnviarNotificacao(n, &wg)
	}

	wg.Wait()
	log.Println("Todas as notificações foram processadas!")

}
func (ns *NotificacaoService) BuscarNotificacoes(horario string) ([]domain.Notificacao, error) {

	var notificacoes []domain.Notificacao

	err := ns.DB.Preload("Usuario").Where("notificacao_time = ?", horario).Find(&notificacoes).Error
	if err != nil {
		return nil, err
	}
	return notificacoes, nil
}

func (s *NotificacaoService) EnviarNotificacao(n domain.Notificacao, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Enviando notificação para %s %s (%s): %s", n.Usuario.Nome, n.Usuario.Sobrenome, n.Usuario.Uid, n.Notificacao_Time)
	time.Sleep(2 * time.Second) // Simula tempo de envio
}
