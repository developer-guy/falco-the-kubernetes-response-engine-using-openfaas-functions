package function

import (
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Alert struct {
	Output       string    `json:"output"`
	Priority     string    `json:"priority"`
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

func Handle(w http.ResponseWriter, r *http.Request) {
	var alert Alert

	if r.Body != nil {
		defer r.Body.Close()

		bodyBytes, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(bodyBytes, &alert)

		slack.PostWebhook(os.Getenv("SLACK_WEBHOOK_URL"), &slack.WebhookMessage{
			Username: "Falco",
			IconURL:  "https://branding.cncf.io/img/projects/falco/icon/color/falco-icon-color.png",
			Text:     fmt.Sprintf("*Output:* %s, *Rule:* %s, *Time:* %s, *Detail:* %+v", alert.Output, alert.Rule, alert.Time, alert.OutputFields),
		})
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
