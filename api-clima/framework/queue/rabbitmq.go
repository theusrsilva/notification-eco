package queue

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type Rabbitmq struct {
	RabbitMQURL  string
	ExchangeName string
	RoutingKey   string
	ContentType  string
}

func NewRabbitmq() *Rabbitmq {
	return &Rabbitmq{}
}

func (r *Rabbitmq) PublicaMensagem(mensagem []byte) error {
	conn, err := amqp.Dial(r.RabbitMQURL)
	if err != nil {
		return fmt.Errorf("erro ao conectar ao RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Criar um canal de comunicação com o RabbitMQ
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("erro ao criar canal: %v", err)
	}
	defer ch.Close()

	// Declarar a exchange (do tipo direct)
	err = ch.ExchangeDeclare(
		r.ExchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("erro ao declarar exchange: %v", err)
	}

	// Publicar a mensagem na exchange com a chave de roteamento
	err = ch.Publish(
		r.ExchangeName, // Nome da exchange
		r.RoutingKey,   // Chave de roteamento
		false,          // Não obrigatório
		false,          // Não urgente
		amqp.Publishing{
			ContentType: r.ContentType,
			Body:        []byte(mensagem),
		},
	)
	if err != nil {
		return fmt.Errorf("erro ao publicar mensagem: %v", err)
	}

	log.Printf("Mensagem enviada: %s", mensagem)
	return nil
}
