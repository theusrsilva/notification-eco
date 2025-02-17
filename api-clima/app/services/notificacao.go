package services

import (
	"api-clima/domain"
	"api-clima/framework/database"
	"api-clima/framework/queue"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/jinzhu/gorm"
	"golang.org/x/net/html/charset"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type NotificacaoService struct {
	DB       *gorm.DB
	Rabbitmq *queue.Rabbitmq
}

type NotificacaoRabbitMq struct {
	Nome      string         `json:"nome"`
	Sobrenome string         `json:"sobrenome"`
	Email     string         `json:"email"`
	Uid       string         `json:"uid"`
	Tipo      string         `json:"tipo"`
	Previsoes CidadePrevisao `json:"mensagem"`
}

type CidadePrevisao struct {
	Nome        string     `xml:"nome"`
	Uf          string     `xml:"uf"`
	Atualizacao string     `xml:"atualizacao"`
	Previsao    []Previsao `xml:"previsao"`
}

type Previsao struct {
	Dia    string  `xml:"dia"`
	Tempo  string  `xml:"tempo"`
	Maxima int     `xml:"maxima"`
	Minima int     `xml:"minima"`
	Iuv    float64 `xml:"iuv"`
}

func NewNotificacaoService(db *gorm.DB, rabbitmq *queue.Rabbitmq) *NotificacaoService {
	return &NotificacaoService{
		DB:       db,
		Rabbitmq: rabbitmq,
	}
}
func (ns *NotificacaoService) ProcessaNotificacoes(horario string) {

	log.Printf("Checando notificações para %s", horario)

	notificacoes, err := ns.BuscarNotificacoes(horario)
	if err != nil {
		log.Printf("Erro ao buscar notificações: %v", err)
		return
	}
	var wg sync.WaitGroup
	for _, n := range notificacoes {
		wg.Add(1)
		go ns.EnviarNotificacao(n, &wg)
	}

	wg.Wait()
	log.Println("Todas as notificações foram processadas!")

}
func (ns *NotificacaoService) BuscarNotificacoes(horario string) ([]domain.Notificacao, error) {

	var notificacoes []domain.Notificacao

	err := ns.DB.Preload("Usuario").Where("notificacao_time = ?", horario).Find(&notificacoes).Error
	if err != nil {
		return nil, err
	}
	return notificacoes, nil
}

func (s *NotificacaoService) EnviarNotificacao(n domain.Notificacao, wg *sync.WaitGroup) {
	defer wg.Done()
	if n.Sms == false && n.Push == false && n.Web == false && n.Email == false {
		return
	}
	clima, err := BuscaClima(n)
	if err != nil {
		log.Println(err)
	}
	mensagem := NotificacaoRabbitMq{
		Nome:      n.Usuario.Nome,
		Sobrenome: n.Usuario.Sobrenome,
		Email:     n.Usuario.Email,
		Uid:       n.Usuario.Uid,
		Previsoes: clima,
	}

	if n.Sms == true {
		mensagem.Tipo = "SMS"
		body, err := json.Marshal(mensagem)
		err = s.Rabbitmq.PublicaMensagem(body)
		if err != nil {
			log.Printf("Erro ao publicar mensagem no RabbitMQ: %v", err)
			return
		}
		log.Printf("Publicando mensagem SMS.")
	}
	if n.Push == true {
		mensagem.Tipo = "PUSH"
		body, err := json.Marshal(mensagem)
		err = s.Rabbitmq.PublicaMensagem(body)
		if err != nil {
			log.Printf("Erro ao publicar mensagem no RabbitMQ: %v", err)
			return
		}
		log.Printf("Publicando mensagem PUSH.")
	}
	if n.Web == true {
		mensagem.Tipo = "WEB"
		body, err := json.Marshal(mensagem)
		err = s.Rabbitmq.PublicaMensagem(body)
		if err != nil {
			log.Printf("Erro ao publicar mensagem no RabbitMQ: %v", err)
			return
		}
		log.Printf("Publicando mensagem WEB.")
	}
	if n.Email == true {
		mensagem.Tipo = "EMAIL"
		body, err := json.Marshal(mensagem)
		err = s.Rabbitmq.PublicaMensagem(body)
		if err != nil {
			log.Printf("Erro ao publicar mensagem no RabbitMQ: %v", err)
			return
		}
		log.Printf("Publicando mensagem EMAIL.")
	}
}

func BuscaClima(n domain.Notificacao) (CidadePrevisao, error) {
	url := fmt.Sprintf(os.Getenv("CPTEC_URL")+"cidade/%v/previsao.xml", n.Usuario.Cidade)

	redis := database.NewRedisClient(0)

	valoresCache, err := redis.Find(strconv.Itoa(n.Usuario.Cidade) + time.Now().Format("2006-01-02"))

	if valoresCache != "" {
		log.Printf("Pegando dados salvos em cache")
		var previsao CidadePrevisao
		err := json.Unmarshal([]byte(valoresCache), &previsao)
		if err != nil {
			return CidadePrevisao{}, fmt.Errorf("erro ao desserializar dados do cache: %v", err)
		}
		return previsao, nil
	}
	resp, err := http.Get(url)
	if err != nil {
		return CidadePrevisao{}, fmt.Errorf("erro ao fazer requisição: %v", err)
	}
	defer resp.Body.Close()

	reader, err := charset.NewReader(resp.Body, "ISO-8859-1")
	if err != nil {
		return CidadePrevisao{}, fmt.Errorf("erro ao configurar leitor de charset: %v", err)
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return CidadePrevisao{}, fmt.Errorf("erro ao ler resposta: %v", err)
	}

	decoder := xml.NewDecoder(strings.NewReader(string(body)))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if charset == "ISO-8859-1" {
			return input, nil
		}
		return nil, fmt.Errorf("não suportado charset: %s", charset)
	}
	var previsaoResponse CidadePrevisao
	err = decoder.Decode(&previsaoResponse)

	cidadesCacheJSON, err := json.Marshal(previsaoResponse)
	if err != nil {

	}
	err = redis.Insert(strconv.Itoa(n.Usuario.Cidade)+time.Now().Format("2006-01-02"), string(cidadesCacheJSON), 24*time.Hour)
	log.Printf("Salvando dados em cache")
	return previsaoResponse, nil
}
