package domain

import (
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"regexp"
	"time"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
	govalidator.TagMap["notificacao-horario"] = govalidator.Validator(func(str string) bool {
		matched, _ := regexp.MatchString(`^(?:[01]\d|2[0-3]):[0-5]\d`, str)
		return matched
	})
}

type Notificacao struct {
	Uid              string    `json:"uid" valid:"uuid" gorm:"type:varchar(36);primary_key"`
	Usuario          *Usuario  `json:"usuario" valid:"-" gorm:"foreignKey:UsuarioId;references:Uid"`
	UsuarioId        string    `json:"-" valid:"-" gorm:"column:usuario_id;type:varchar(36);notnull"`
	Sms              bool      `json:"aceita_sms" valid:"-"`
	Push             bool      `json:"aceita_push" valid:"-"`
	Web              bool      `json:"aceita_web" valid:"-"`
	Email            bool      `json:"aceita_email" valid:"-"`
	Notificacao_Time string    `json:"horario" valid:"notificacao-horario" gorm:"type:varchar(10)"`
	CreatedAt        time.Time `json:"created_at" valid:"-"`
	UpdatedAt        time.Time `json:"updated_at" valid:"-"`
}

func (Notificacao) TableName() string {
	return "notificacoes"
}

func NewNotificacao(usuario *Usuario, sms bool, push bool, web bool, email bool, time string) (*Notificacao, error) {
	//alterar conforme liberar novas features
	notificacao := Notificacao{
		Sms:              false,
		Push:             false,
		Web:              web,
		Email:            false,
		Notificacao_Time: time,
	}

	notificacao.prepare()
	err := notificacao.Validate()
	if err != nil {
		return nil, err
	}
	return &notificacao, nil
}
func NewNotificacaoEmpty() *Notificacao {
	return &Notificacao{}
}

func (notificacao *Notificacao) prepare() {
	notificacao.Uid = uuid.New().String()
	notificacao.CreatedAt = time.Now()
	notificacao.UpdatedAt = time.Now()
}

func (notificacao *Notificacao) Validate() error {
	_, err := govalidator.ValidateStruct(notificacao)
	if err != nil {
		return err
	}
	return nil
}
