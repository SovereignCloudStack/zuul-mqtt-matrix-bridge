package main

import (
	"flag"
	slogenv "github.com/cbrewster/slog-env"
	"github.com/SovereignCloudStack/zuul-mqtt-matrix-bridge/pkg"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	matrixHomeserver := flag.String("matrix-homeserver", "https://matrix.org", "Matrix homeserver")
	matrixToken := flag.String("matrix-token", "", "Matrix access token")
	matrixRoomID := flag.String("matrix-room-id", "", "Matrix room ID")
	matrixTemplateFile := flag.String("matrix-msg-template", "./templates/matrix_message.tmpl", "Matrix template file")

	mqttBroker := flag.String("mqtt-broker", "tcp://localhost:1883", "The MQTT Broker")
	mqttUser := flag.String("mqtt-user", "", "The MQTT user")
	mqttPassword := flag.String("mqtt-pass", "", "The MQTT password")
	mqttTopic := flag.String("mqtt-topic", "zuul/#", "The MQTT topic to subscribe")
	mqttTopicQOS := flag.Int("mqtt-topic-qos", 0, "The MQTT topic to subscribe")

	flag.Parse()

	logger := slog.New(slogenv.NewHandler(slog.NewTextHandler(os.Stdout, nil)))
	slog.SetDefault(logger)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go bridge.New(
		matrixHomeserver,
		matrixToken,
		matrixRoomID,
		matrixTemplateFile,
		mqttBroker,
		mqttUser,
		mqttPassword,
		mqttTopic,
		mqttTopicQOS,
	)
	<-c
}
