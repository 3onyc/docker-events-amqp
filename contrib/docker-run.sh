main() {
    /go/bin/app \
        -docker="unix:///docker.sock" \
        -amqp="amqp://${AMQP_USER}:${AMQP_PASSWORD}@${AMQP_PORT_5672_TCP_ADDR}:${AMQP_PORT_5672_TCP_PORT}" \
        -exchange="${AMQP_EXCHANGE}"
}

main "$@"
