package function

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	handler "github.com/openfaas/templates-sdk/go-http"
)

type Level string

const (
	NOTICE  Level = "Notice"
	WARNING Level = "Warning"
)

type Alert struct {
	Output       string    `json:"output"`
	Priority     Level     `json:"priority"`
	Rule         string    `json:"rule"`
	Time         time.Time `json:"time"`
	OutputFields struct {
		ContainerID              string      `json:"container.id"`
		ContainerImageRepository interface{} `json:"container.image.repository"`
		ContainerImageTag        interface{} `json:"container.image.tag"`
		EvtTime                  int64       `json:"evt.time"`
		FdName                   string      `json:"fd.name"`
		K8SNsName                string      `json:"k8s.ns.name"`
		K8SPodName               string      `json:"k8s.pod.name"`
		ProcCmdline              string      `json:"proc.cmdline"`
	} `json:"output_fields"`
}

// Handle a function invocation
func Handle(req handler.Request) (handler.Response, error) {
	var alert Alert
	json.Unmarshal(req.Body, &alert)

	if alert.Priority == NOTICE {
		log.Println("Sending alert to notifier-fn")
		http.Post("http://gateway.openfaas:8080/function/notifier-fn", "application/json", bytes.NewBuffer(req.Body))
	} else if alert.Priority == WARNING && alert.OutputFields.K8SPodName != "" {
		log.Println("Sending alert to delete-pod-fn")
		http.Post("http://gateway.openfaas:8080/function/delete-pod-fn", "application/json", bytes.NewBuffer(req.Body))
	}

	return handler.Response{
		Body:       []byte(req.Body),
		StatusCode: http.StatusOK,
	}, nil
}
