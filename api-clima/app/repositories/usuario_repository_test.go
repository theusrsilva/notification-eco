package repositories_test

import (
	"api-clima/app/repositories"
	"api-clima/domain"
	"api-clima/framework/database"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUsuarioRepositoryDbInsertAndFind(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	usuario, err := domain.NewUsuario("Matheus", "Rocha", "matheusEmail@gmail.com", 1234)

	if err != nil {
		t.Fatal("Erro ao criar de usuario", err)
	}

	notificacao, err := domain.NewNotificacao(usuario, false, false, false, false, "22:00")

	if err != nil {
		t.Fatal("Erro ao criar notificação", err)
	}

	repository := repositories.UsuarioRepositoryDb{Db: db}

	repository.Insert(usuario, notificacao)

	userFind, err := repository.Find(usuario.Uid)

	require.NotEmpty(t, userFind.Uid)
	require.Nil(t, err)
	require.Equal(t, usuario.Uid, userFind.Uid)

}

func TestUsuarioRepositoryDbUpdate(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()
	usuario, err := domain.NewUsuario("Matheus", "Rocha", "matheusEmail@gmail.com", 1234)

	if err != nil {
		t.Fatal("Erro ao criar de usuario", err)
	}

	notificacao, err := domain.NewNotificacao(usuario, false, false, false, false, "22:00")

	if err != nil {
		t.Fatal("Erro ao criar notificação", err)
	}

	repository := repositories.UsuarioRepositoryDb{Db: db}

	repository.Insert(usuario, notificacao)

	usuario.Sobrenome = "teste"

	repository.Update(usuario)

	userFind, err := repository.Find(usuario.Uid)

	require.NotEmpty(t, userFind.Uid)
	require.Nil(t, err)
	require.Equal(t, usuario.Sobrenome, userFind.Sobrenome)
}
