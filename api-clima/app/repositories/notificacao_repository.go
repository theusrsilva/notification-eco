package repositories

import (
	"api-clima/domain"
	"fmt"
	"github.com/jinzhu/gorm"
)

type NotificacaoRepository interface {
	Find(notificacao *domain.Notificacao) (*domain.Notificacao, error)
	Update(notificacao *domain.Notificacao) (*domain.Notificacao, error)
}

type NotificacaoRepositoryDb struct {
	Db *gorm.DB
}

func (repository NotificacaoRepositoryDb) Find(notificacao_uid string) (*domain.Notificacao, error) {
	var notificacao domain.Notificacao

	repository.Db.Preload("Usuario").First(&notificacao, "uid = ?", notificacao_uid)
	if notificacao.Uid == "" {
		return nil, fmt.Errorf("notificação não existe")
	}
	return &notificacao, nil

}
func (repository NotificacaoRepositoryDb) Update(notificacao *domain.Notificacao) error {
	err := repository.Db.Save(&notificacao).Error

	if err != nil {
		return err
	}

	return nil

}
