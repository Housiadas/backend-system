package userbus

import (
	"encoding/json"
	"fmt"
	"github.com/Housiadas/backend-system/business/sys/delegate"

	"github.com/google/uuid"
)

// Domain represents the name of this domain.
const Domain = "user"

// ActionUpdated Set of delegate actions.
const (
	ActionUpdated = "updated"
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

// Marshal returns the event parameters encoded as JSON.
func (au *ActionUpdatedParms) Marshal() ([]byte, error) {
	return json.Marshal(au)
}

// ActionUpdatedData constructs the data for the updated action.
func ActionUpdatedData(uu UpdateUser, userID uuid.UUID) delegate.Data {
	params := ActionUpdatedParms{
		UserID: userID,
		UpdateUser: UpdateUser{
			Enabled: uu.Enabled,
		},
	}

	rawParams, err := params.Marshal()
	if err != nil {
		panic(err)
	}

	return delegate.Data{
		Domain:    Domain,
		Action:    ActionUpdated,
		RawParams: rawParams,
	}
}
