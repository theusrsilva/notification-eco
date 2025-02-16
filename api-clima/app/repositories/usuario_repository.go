package repositories

import (
	"api-clima/domain"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type UsuarioRepository interface {
	Insert(usuario *domain.Usuario, notificacao *domain.Notificacao) (*domain.Usuario, error)
	Find(usuario *domain.Usuario) (*domain.Usuario, error)
	Update(usuario *domain.Usuario) error
}

type UsuarioRepositoryDb struct {
	Db *gorm.DB
}

func (repository UsuarioRepositoryDb) Insert(usuario *domain.Usuario, notificacao *domain.Notificacao) (*domain.Usuario, error) {
	err := repository.Db.Create(usuario).Error

	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		if mysqlErr.Number == 1062 {
			return nil, fmt.Errorf("o email '%s' já está em uso", usuario.Email)
		}
	}

	notificacao.UsuarioId = usuario.Uid

	err = repository.Db.Create(notificacao).Error

	if err != nil {
		return nil, err
	}

	var usuarioComNotificacao domain.Usuario
	err = repository.Db.Preload("Notificacao").First(&usuarioComNotificacao, "uid = ?", usuario.Uid).Error
	if err != nil {
		return nil, err
	}
	return &usuarioComNotificacao, nil
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
