package status

import (
	"time"

	"github.com/go-openapi/strfmt"
)

type ContainerReport struct {
	MessageID   string          `json:"message-id"`
	ContainerID string          `json:"container-id"`
	Status      string          `json:"status"`
	Timestamp   strfmt.DateTime `json:"time"`
}

func (cr *ContainerReport) Time() time.Time {
	return time.Time(cr.Timestamp)
}
