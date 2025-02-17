FROM rabbitmq:management

RUN apt-get update && apt-get install -y curl && \
    curl -LO https://raw.githubusercontent.com/rabbitmq/rabbitmq-management/v3.8.x/bin/rabbitmqadmin && \
    chmod +x rabbitmqadmin && \
    mv rabbitmqadmin /usr/local/bin/

CMD ["bash", "-c", " \
    rabbitmq-server & \
    sleep 10 && \
    rabbitmqadmin declare exchange name=notificacoes_exchange type=direct durable=true && \
    rabbitmqadmin declare exchange name=fila_dead_letter_notificacoes type=direct durable=true && \
    rabbitmqadmin declare queue name=fila_notificacoes durable=true arguments='{\"x-dead-letter-exchange\":\"fila_dead_letter_notificacoes\"}' && \
    rabbitmqadmin declare queue name=fila_dead_letter_notificacoes durable=true && \
    rabbitmqadmin declare binding source=notificacoes_exchange destination=fila_notificacoes routing_key=notificacao && \
    rabbitmqadmin declare binding source=fila_dead_letter_notificacoes destination=fila_dead_letter_notificacoes routing_key=dead && \
    wait"]