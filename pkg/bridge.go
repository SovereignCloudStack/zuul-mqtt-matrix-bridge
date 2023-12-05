package bridge

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/grokify/html-strip-tags-go"
	"github.com/matrix-org/gomatrix"
)

type jsonTime time.Time

type MqttMessage struct {
	Action      string      `json:"action"`
	Tenant      string      `json:"tenant"`
	Pipeline    string      `json:"pipeline"`
	Project     string      `json:"project"`
	Branch      string      `json:"branch"`
	ChangeURL   string      `json:"change_url"`
	Message     string      `json:"message"`
	TriggerTime *jsonTime   `json:"trigger_time"`
	EnqueueTime *jsonTime   `json:"enqueue_time"`
	Change      interface{} `json:"change"`
	Patchset    string      `json:"patchset"`
	CommitID    string      `json:"commit_id"`
	Owner       string      `json:"owner"`
	Ref         string      `json:"ref"`
	ZuulRef     string      `json:"zuul_ref"`
	Buildset    struct {
		UUID   string `json:"uuid"`
		Result string `json:"result,omitempty"`
		Builds []struct {
			JobName      string        `json:"job_name"`
			Voting       bool          `json:"voting"`
			UUID         string        `json:"uuid,omitempty"`
			ExecuteTime  *jsonTime     `json:"execute_time"`
			StartTime    *jsonTime     `json:"start_time,omitempty"`
			EndTime      *jsonTime     `json:"end_time,omitempty"`
			LogURL       string        `json:"log_url,omitempty"`
			WebURL       string        `json:"web_url,omitempty"`
			Result       string        `json:"result,omitempty"`
			Dependencies []interface{} `json:"dependencies,omitempty"`
			Artifacts    []interface{} `json:"artifacts,omitempty"`
		} `json:"builds"`
	} `json:"buildset"`
}

// https://spec.matrix.org/v1.8/client-server-api/#mtext
type MatrixRequestBody struct {
	Body          string `json:"body"`
	Format        string `json:"format"`
	FormattedBody string `json:"formatted_body"`
	Msgtype       string `json:"msgtype"`
}

func (t *jsonTime) UnmarshalJSON(s []byte) (err error) {
	q, err := strconv.ParseFloat(string(s), 64)

	if err != nil {
		return err
	}
	*(*time.Time)(t) = time.Unix(int64(q), 0)
	return
}

func (t jsonTime) String() string {
	return time.Time(t).String()
}

func ErrAttr(err error) slog.Attr {
	return slog.Any("error", err)
}

func connectMatrix(matrixHomeserver *string, matrixToken *string) *gomatrix.Client {
	matrixClient, err := gomatrix.NewClient(*matrixHomeserver, "", *matrixToken)
	if err != nil {
		slog.Error("Matrix: unable to connect", ErrAttr(err))
		os.Exit(1)
	}

	return matrixClient
}

func connectMqtt(mqttBroker *string, mqttUser *string, mqttPassword *string) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(*mqttBroker)
	opts.SetClientID(*mqttUser)
	opts.SetUsername(*mqttUser)
	opts.SetPassword(*mqttPassword)

	mqttClient := mqtt.NewClient(opts)
	token := mqttClient.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		slog.Error("MQTT: unable to connect", ErrAttr(token.Error()))
		os.Exit(1)
	}

	return mqttClient
}

func postToMatrix(matrixClient *gomatrix.Client, matrixRoomID *string, matrixTemplate *template.Template, topic string, message []byte) {
	var mqttMessage MqttMessage
	matrixMessage := new(bytes.Buffer)

	err := json.Unmarshal(message, &mqttMessage)
	if err != nil {
		slog.Error("Bridge: unable to unmarshal mqtt message", ErrAttr(err))
	}

	err = matrixTemplate.Execute(matrixMessage, mqttMessage)
	if err != nil {
		slog.Error("Bridge: unable to create matrix message", ErrAttr(err))
	}

	matrixRequest := MatrixRequestBody{
		strip.StripTags(matrixMessage.String()),
		"org.matrix.custom.html",
		matrixMessage.String(),
		"m.text",
	}

	_, err = matrixClient.SendMessageEvent(*matrixRoomID, "m.room.message", matrixRequest)
	if err != nil {
		slog.Error("Matrix: unable to send message", ErrAttr(err))
	}
}

func New(
	matrixHomeserver *string,
	matrixToken *string,
	matrixRoomID *string,
	matrixTemplateFile *string,
	mqttBroker *string,
	mqttUser *string,
	mqttPassword *string,
	mqttTopic *string,
	mqttTopicQOS *int) {

	matrixClient := connectMatrix(matrixHomeserver, matrixToken)
	mqttClient := connectMqtt(mqttBroker, mqttUser, mqttPassword)

	// Zuul's MQTT report does not contain a buildset URL; therefore, we use the following workaround:
	// Extract the base URL from the build URL, and then within the template file, we can construct and use the buildset URL.
	// The extract function is taken from: https://stackoverflow.com/questions/72387330/how-to-extract-base-url-using-golang
	funcMap := template.FuncMap{
		"getBaseUrl": func(rawUrl string) string {
			url, err := url.Parse(rawUrl)
			if err != nil {
				slog.Error("URLparse: unable to parse URL", slog.Any("error", err))
				return ""
			}
			url.Path = ""
			url.RawQuery = ""
			url.Fragment = ""
			return url.String()
		},
	}

	matrixTemplate, err := template.New(filepath.Base(*matrixTemplateFile)).Funcs(funcMap).ParseFiles(*matrixTemplateFile)

	if err != nil {
		slog.Error("Bridge: unable to load matrix message template", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info(
		"MQTT->Matrix bringe established",
		slog.String("MQTT broker", *mqttBroker),
		slog.String("Matrix homeserver", *matrixHomeserver),
	)

	mqttClient.Subscribe(*mqttTopic, byte(*mqttTopicQOS), func(client mqtt.Client, msg mqtt.Message) {
		slog.Debug(
			"Received message:",
			slog.String("topic", msg.Topic()),
			slog.String("message", string(msg.Payload())),
		)
		go postToMatrix(matrixClient, matrixRoomID, matrixTemplate, msg.Topic(), msg.Payload())
	})

}
