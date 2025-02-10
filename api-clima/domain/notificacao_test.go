package domain_test

import (
	"api-clima/domain"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewNotificacao(t *testing.T) {

	usuario, err := domain.NewUsuario("matheus", "rocha", "theusrsilva@gmail.com", 2012)

	notificacao, err := domain.NewNotificacao(usuario, false, false, true, false, "23:30")

	require.NotNil(t, notificacao)
	require.Nil(t, err)
}
func TestNewNotificacaoHorarioInvalido(t *testing.T) {

	usuario, err := domain.NewUsuario("matheus", "rocha", "theusrsilva@gmail.com", 2012)

	notificacao, err := domain.NewNotificacao(usuario, false, false, true, false, "25:30")

	require.Error(t, err)
	require.Nil(t, notificacao)
}
