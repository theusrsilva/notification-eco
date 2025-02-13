# Ecossistema de notificação (notification-eco)

### To-Do
- [ ] Criar um cadastro básico de usuário na api-clima definindo que tipos de notificações ele aceita
- [ ] instalar RabbitMQ no docker compose e configurar
- [ ] Criar service para interagir com a mensageria
- [ ] Instalar scheduler para rodar cron
- [ ] Feature para rodar integrada com o cron e selecionar usuários que devem receber notificações naquele momento
- [ ] Feature de busca de clima por usuário
- [ ] Feature para lançar notificação para a mensageria
- [ ] Feature de cachear clima buscado no dia para aquela cidade especifica
- [ ] Criar base da api-notificações
- [ ] Criar service para interagir com o serviço de mensageria
- [ ] Feature para disparar as notificações para os lugares corretos
- [ ] Configurar docker para rodar sozinho main das api`s
- [ ] Configurar docker para clonar .env.example para .env

### Estrutura

- Api`s : Golang 1.23.5
- Mensageria : RabbitMq
- Banco de Dados : Mysql 8.0
- Cache : Redis
- Ambiente : Docker

### Estrutura Base
![estrutura-base.png](resources/estrutura-base.png)

### Api-clima estrutura
![estrutura-api-clima.png](resources/estrutura-api-clima.png)

### Futuro

- [ ] Integrar com keycloack para gerenciamento de usuários
- [ ] Integrar com serviço de socket para mensagens em real time
- [ ] Integrar com demais serviços para disparar os demais tipos de notificações