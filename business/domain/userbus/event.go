package userbus

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/pkg/kafka"
)

const (
	UserUpdatedEvent = "user-updated"
)

// ActionUpdatedParms represents the parameters for the updated action.
type ActionUpdatedParms struct {
	UserID uuid.UUID
	UpdateUser
}

// String returns a string representation of the action parameters.
func (au *ActionUpdatedParms) String() string {
	return fmt.Sprintf("&EventParamsUpdated{UserID:%v, Enabled:%v}", au.UserID, au.Enabled)
}

// ActionUpdatedData constructs the data for the updated action.
func ActionUpdatedData(uu UpdateUser, userID uuid.UUID) kafka.Event {
	params := ActionUpdatedParms{
		UserID:     userID,
		UpdateUser: uu,
	}

	rawData, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}

	return kafka.Event{
		Topic: UserUpdatedEvent,
		Data:  rawData,
	}
}
