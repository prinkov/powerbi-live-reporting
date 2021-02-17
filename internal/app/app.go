package app

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"log"
	"os"
	"powerbi-live-reporting/config"
	"powerbi-live-reporting/internal/publisher"
	"powerbi-live-reporting/internal/type"
	"sync"
)

type App struct {
	config *config.Config
	logger *zap.Logger
	amqpConnection *amqp.Connection
	waitGroup *sync.WaitGroup
}

func CreateApp() (app *App) {
	app = &App{}
	c, err := config.ReadConfig()
	if err != nil {
		log.Fatal("Config error", zap.Error(err))
	}
	app.config = c

	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.Level = zap.NewAtomicLevel()
	err = zapConfig.Level.UnmarshalText([]byte(c.LogLevel))
	if err != nil {
		log.Fatal("Error when set log level", zap.Error(err), zap.String("log level", c.LogLevel))
	}

	app.logger, err = zapConfig.Build()
	if err != nil {
		log.Fatal("Error when create logger", zap.Error(err))
	}

	app.waitGroup = new(sync.WaitGroup)
	amqpUri := amqp.URI{
		Scheme:   "amqp",
		Vhost:    "/",
		Host:     app.config.Host,
		Port:     app.config.Port,
		Username: app.config.Username,
		Password: app.config.Password,
	}
	app.amqpConnection, err = amqp.Dial(amqpUri.String())

	if err != nil {
		app.logger.Error("AMQP connection error", zap.Error(err), zap.String("AMQP URI", amqpUri.String()))
	}

	return app
}

func (app *App) Run() {
	defer app.amqpConnection.Close()
	channel, err := app.amqpConnection.Channel()
	if err != nil {
		app.logger.Fatal("Error when opening channel", zap.Error(err))
	}

	queue, err := channel.QueueDeclare("sales_events", true, false, false, false, nil)
	if err != nil {
		app.logger.Fatal("Error when declare queue", zap.Error(err))
	}

	messageChannel, err := channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	for i := 0; i < app.config.ConsumersCount; i++ {
		app.waitGroup.Add(1)
		go app.ConsumeChannel(messageChannel, app.config.PBIUrl)
	}
	app.waitGroup.Wait()
}

func (app *App) ConsumeChannel(messageChannel <-chan amqp.Delivery, url string)  {
	app.logger.Info("Consumer ready", zap.Int("pid", os.Getpid()))

	p := publisher.PBIPublisher{Url: app.config.PBIUrl, Logger: app.logger}
	for d := range messageChannel {
		app.logger.Debug("Received a message: %s", zap.ByteString("Body", d.Body))

		sale := &_type.SaleItem{}
		err := json.Unmarshal(d.Body, sale)
		if err != nil {
			app.logger.Error("Error decoding message from JSON", zap.Error(err))
		}

		if err := d.Ack(false); err != nil {
			app.logger.Error("Error acknowledging message:", zap.Error(err))
		} else {
			var data [1]*_type.SaleItem
			data[0] = sale
			app.logger.Info("Info", zap.Any("raw data", data))

			dataArr, err := json.Marshal(data)
			if err != nil {
				app.logger.Error("Error decoding message from JSON", zap.Error(err))
			}

			p.Publish(dataArr)

			app.logger.Info("Acknowledged message", zap.String("send data", string(dataArr)))
		}
	}

	defer app.waitGroup.Done()
}