package productbus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/foundation/kafka"
)

const (
	ProductUpdatedEvent = "product-updated"
	ProductDeletedEvent = "product-deleted"
)

// actionUserUpdated is executed by the user domain indirectly when a user is updated.
func (c *Core) actionUserUpdated(ctx context.Context, event kafka.Event) error {
	var params userbus.ActionUpdatedParms
	err := json.Unmarshal(event.Data, &params)
	if err != nil {
		return fmt.Errorf("expected an encoded %T: %w", params, err)
	}

	c.log.Info(ctx, "action user-update", "user_id", params.UserID, "enabled", params.Enabled)

	// Now we can see if this user has been disabled. If they have been, we will
	// want to disable to mark all these products as deleted. Right now we don't
	// have support for this, but you can see how we can process the event.

	return nil
}