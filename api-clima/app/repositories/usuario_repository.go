package repositories

import (
	"api-clima/domain"
	"fmt"
	"github.com/jinzhu/gorm"
)

type UsuarioRepository interface {
	Insert(usuario *domain.Usuario) (*domain.Usuario, error)
	Find(usuario *domain.Usuario) (*domain.Usuario, error)
	Update(usuario *domain.Usuario) error
}

type UsuarioRepositoryDb struct {
	Db *gorm.DB
}

func (repository UsuarioRepositoryDb) Insert(usuario *domain.Usuario, notificacao *domain.Notificacao) (*domain.Usuario, error) {
	err := repository.Db.Create(usuario).Error

	if err != nil {
		return nil, err
	}

	notificacao.UsuarioId = usuario.Uid

	err = repository.Db.Create(notificacao).Error

	if err != nil {
		return nil, err
	}

	return usuario, nil
}

func (repository UsuarioRepositoryDb) Find(usuario_id string) (*domain.Usuario, error) {
	var usuario domain.Usuario

	repository.Db.Preload("Notificacao").First(&usuario, "uid = ?", usuario_id)
	if usuario.Uid == "" {
		return nil, fmt.Errorf("usuario não existe")
	}
	return &usuario, nil

}
func (repository UsuarioRepositoryDb) Update(usuario *domain.Usuario) error {
	err := repository.Db.Save(&usuario).Error

	if err != nil {
		return err
	}

	return nil

}
