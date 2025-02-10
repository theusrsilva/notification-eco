package domain_test

import (
	"api-clima/domain"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewUsuario(t *testing.T) {

	usuario, err := domain.NewUsuario("matheus", "rocha", "theusrsilva@gmail.com", 2000)

	require.NotNil(t, usuario)
	require.Nil(t, err)
}
