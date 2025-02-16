package domain

import (
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"time"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Usuario struct {
	Uid         string       `json:"uid" valid:"uuid" gorm:"type:varchar(36);primary_key"`
	Nome        string       `json:"nome" valid:"notnull" gorm:"type:varchar(50)"`
	Sobrenome   string       `json:"sobrenome" valid:"notnull" gorm:"type:varchar(60)"`
	Email       string       `json:"email" valid:"email" gorm:"type:varchar(100);unique_index"`
	Cidade      int          `json:"cidade_id" valid:"-" gorm:"index"`
	CreatedAt   time.Time    `json:"created_at" valid:"-"`
	UpdatedAt   time.Time    `json:"updated_at" valid:"-"`
	Notificacao *Notificacao `json:"notificacoes" valid:"-" gorm:"foreignKey:UsuarioId"`
}

func (Usuario) TableName() string {
	return "usuarios"
}
func NewUsuario(nome string, sobrenome string, email string, cidade int) (*Usuario, error) {
	usuario := Usuario{
		Nome:      nome,
		Sobrenome: sobrenome,
		Email:     email,
		Cidade:    cidade,
	}
	usuario.prepare()
	err := usuario.Validate()
	if err != nil {
		return nil, err
	}
	return &usuario, nil
}
func NewUsuarioEmpty() *Usuario {
	return &Usuario{}
}

func (usuario *Usuario) prepare() {
	usuario.Uid = uuid.New().String()
	usuario.CreatedAt = time.Now()
	usuario.UpdatedAt = time.Now()
}

func (usuario *Usuario) Validate() error {
	_, err := govalidator.ValidateStruct(usuario)
	if err != nil {
		return err
	}
	return nil
}
