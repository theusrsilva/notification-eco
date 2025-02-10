package repositories_test

import (
	"api-clima/app/repositories"
	"api-clima/domain"
	"api-clima/framework/database"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNotificacaoRepositoryDbInsertAndFind(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	usuario, err := domain.NewUsuario("Matheus", "Rocha", "matheusEmail@gmail.com", 1234)

	if err != nil {
		t.Fatal("Erro ao criar de usuario", err)
	}

	notificacao, err := domain.NewNotificacao(usuario, false, false, true, false, "22:00")

	if err != nil {
		t.Fatal("Erro ao criar notificação", err)
	}

	repository := repositories.UsuarioRepositoryDb{Db: db}

	repository.Insert(usuario, notificacao)
	UsuarioFind, err := repository.Find(usuario.Uid)

	repositoryNot := repositories.NotificacaoRepositoryDb{Db: db}
	notificacaoFind, err := repositoryNot.Find(UsuarioFind.Notificacao.Uid)

	require.NotEmpty(t, notificacaoFind.Uid)
	require.Nil(t, err)
	require.Equal(t, UsuarioFind.Notificacao.Uid, notificacaoFind.Uid)

}

func TestNotificacaoRepositoryDbInsertUpdateFind(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	usuario, err := domain.NewUsuario("Matheus", "Rocha", "matheusEmail@gmail.com", 1234)

	if err != nil {
		t.Fatal("Erro ao criar de usuario", err)
	}

	notificacao, err := domain.NewNotificacao(usuario, false, false, true, false, "22:00")

	if err != nil {
		t.Fatal("Erro ao criar notificação", err)
	}

	repository := repositories.UsuarioRepositoryDb{Db: db}

	repository.Insert(usuario, notificacao)
	UsuarioFind, err := repository.Find(usuario.Uid)

	repositoryNot := repositories.NotificacaoRepositoryDb{Db: db}
	notificacaoFind, err := repositoryNot.Find(UsuarioFind.Notificacao.Uid)
	notificacaoFind.Notificacao_Time = "10:00"
	repositoryNot.Update(notificacaoFind)

	notificacaoFind2, err := repositoryNot.Find(notificacaoFind.Uid)

	require.NotEmpty(t, notificacaoFind.Uid)
	require.Nil(t, err)
	require.Equal(t, "10:00", notificacaoFind2.Notificacao_Time)

}
