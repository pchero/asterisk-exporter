package amihandler

import (
	"github.com/ivahaev/amigo"
	"github.com/ivahaev/amigo/uuid"
	log "github.com/sirupsen/logrus"
)

// AMIHandler interface
type AMIHandler interface {
	Connect()
}

type amiHandler struct {
	host     string
	username string
	password string
	port     string

	sock *amigo.Amigo
}

// NewAMIHandler returns AMIHandler
func NewAMIHandler(host, port, username, password string) AMIHandler {
	h := &amiHandler{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}

	return h
}

func (h *amiHandler) Connect() {
	settings := &amigo.Settings{
		Username: h.username,
		Password: h.password,
		Host:     h.host,
		Port:     h.port,
	}

	a := amigo.New(settings)

	a.Connect()
	h.sock = a

	a.On("connect", func(message string) {
		log.WithFields(log.Fields{
			"Username": h.username,
			"Host":     h.host,
		}).Infof("Connected to Asterisk. message: %s", message)

		h.resetMetrics()
	})

	a.On("error", func(message string) {
		log.WithFields(log.Fields{
			"Username": h.username,
			"Host":     h.host,
		}).Errorf("Connection error. message: %s", message)
	})

	a.Connected()
}

func (h *amiHandler) resetMetrics() error {
	id := uuid.NewV4()
	res, err := h.sock.Action(map[string]string{
		"Action":   "CoreShowChannels",
		"ActionID": id,
	})
	if err != nil {
		return err
	}
	log.Debugf("res: %s", res)

	return nil
}
