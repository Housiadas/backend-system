package auditcore

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/internal/core/domain/audit"
	"github.com/Housiadas/backend-system/internal/core/domain/entity"
	"github.com/Housiadas/backend-system/internal/core/domain/name"
)

// TestNewAudits is a helper method for testing.
func TestNewAudits(n int, actorID uuid.UUID, objEntity entity.Entity, action string) []audit.NewAudit {
	newAudits := make([]audit.NewAudit, n)

	idx := rand.Intn(10000)
	for i := range n {
		idx++

		na := audit.NewAudit{
			ObjID:     uuid.New(),
			ObjEntity: objEntity,
			ObjName:   name.MustParse(fmt.Sprintf("ObjName%d", idx)),
			ActorID:   actorID,
			Action:    action,
			Data:      struct{ Name string }{Name: fmt.Sprintf("Name%d", idx)},
			Message:   fmt.Sprintf("Message%d", idx),
		}

		newAudits[i] = na
	}

	return newAudits
}

// TestSeedAudits is a helper method for testing.
func TestSeedAudits(
	ctx context.Context,
	n int,
	actorID uuid.UUID,
	objEntity entity.Entity,
	action string,
	api *Core,
) ([]audit.Audit, error) {
	newAudits := TestNewAudits(n, actorID, objEntity, action)

	audits := make([]audit.Audit, len(newAudits))
	for i, na := range newAudits {
		adt, err := api.Create(ctx, na)
		if err != nil {
			return nil, fmt.Errorf("seeding audit: idx: %d : %w", i, err)
		}

		audits[i] = adt
	}

	return audits, nil
}
