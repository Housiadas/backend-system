package userbus

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/foundation/kafka"
)

const (
	Domain           = "users"
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
	eventData := kafka.EventData{
		Action: UserUpdatedEvent,
		Params: ActionUpdatedParms{
			UserID: userID,
			UpdateUser: UpdateUser{
				Enabled: uu.Enabled,
			},
		},
	}

	rawEventData, err := json.Marshal(eventData)
	if err != nil {
		panic(err)
	}

	return kafka.Event{
		Topic: Domain,
		Data:  rawEventData,
	}
}
