package manager

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"
	"github.com/shipyard/shipyard"
	"github.com/shipyard/shipyard/utils"
)

type (
	EventHandler struct {
		Manager Manager
	}
)

func (h *EventHandler) Handle(e *dockerclient.Event) error {
	log.Infof("日志事件: 日期=%d 状态=%s 容器=%s", e.Time, e.Status, e.ID[:12])
	h.logDockerEvent(e)
	return nil
}

func (h *EventHandler) logDockerEvent(e *dockerclient.Event) error {
	info, err := h.Manager.Container(e.ID)
	if err != nil {
		return err
	}

	ts, err := utils.FromUnixTimestamp(e.Time)
	if err != nil {
		return err
	}

	evt := &shipyard.Event{
		Type: e.Status,
		Message: fmt.Sprintf("action=%s container=%s",
			e.Status, e.ID[:12]),
		Time:          *ts,
		ContainerInfo: info,
		Tags:          []string{"docker"},
	}
	if err := h.Manager.SaveEvent(evt); err != nil {
		return err
	}
	return nil
}
