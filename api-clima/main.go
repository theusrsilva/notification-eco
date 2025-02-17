package main

import (
	"api-clima/app/services"
	"api-clima/framework/database"
	"api-clima/framework/queue"
	"api-clima/framework/server"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"strconv"
	"time"
)

var SaoPauloTimeZone *time.Location

func init() {
	var err error
	SaoPauloTimeZone, err = time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		fmt.Println("Erro ao carregar fuso horário:", err)
		return
	}

}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro carregando .env")
	}

	autoMigrateDb, err := strconv.ParseBool(os.Getenv("AUTO_MIGRATE_DB"))
	if err != nil {
		log.Fatalf("Erro convertendo env AUTO_MIGRATE_DB")
	}

	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		log.Fatalf("Erro convertendo env DEBUG")
	}
	db := database.NewDb()
	db.AutoMigrateDb = autoMigrateDb
	db.Debug = debug
	db.DsnTest = os.Getenv("DSN_TEST")
	db.Dsn = os.Getenv("DSN")
	db.DbTypeTest = os.Getenv("DB_TYPE_TEST")
	db.DbType = os.Getenv("DB_TYPE")
	db.Env = os.Getenv("ENV")

	rabbitmq := queue.NewRabbitmq()
	rabbitmq.RabbitMQURL = os.Getenv("RABBITMQ_URL")
	rabbitmq.ExchangeName = os.Getenv("EXCHANGE_NAME")
	rabbitmq.RoutingKey = os.Getenv("ROUTING_KEY")
	rabbitmq.ContentType = os.Getenv("CONTENT_TYPE")

	dbConnection, err := db.Connect()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	log.Println("Conexão com o banco iniciada.")
	defer dbConnection.Close()

	c := cron.New()
	_, err = c.AddFunc("* * * * *", func() {
		horario := time.Now().In(SaoPauloTimeZone).Format("15:04")
		notificacaoService := services.NewNotificacaoService(dbConnection, rabbitmq)
		notificacaoService.ProcessaNotificacoes(horario)
	})
	if err != nil {
		log.Fatalf("Erro ao adicionar cron job: %v", err)
	}

	log.Println("Cron job iniciado! Verificando notificações a cada minuto.")
	c.Start()

	srv := server.NewServer(dbConnection)
	srv.Start(":8080")
}
