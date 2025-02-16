package services

import (
	"api-clima/app/repositories"
	"api-clima/domain"
	"github.com/jinzhu/gorm"
)

func InsertUsuario(usuario *domain.Usuario, time string, db *gorm.DB) (*domain.Usuario, error) {

	notificacao, err := domain.NewNotificacao(usuario, false, false, true, false, time)
	if err != nil {
		return nil, err
	}
	repository := repositories.UsuarioRepositoryDb{Db: db}
	return repository.Insert(usuario, notificacao)
}
