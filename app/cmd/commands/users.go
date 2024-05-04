package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Housiadas/backend-system/business/data/sqldb"
	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/business/domain/userbus/stores/userdb"
	"github.com/Housiadas/backend-system/foundation/logger"
)

// Users retrieves all users from the database.
func Users(log *logger.Logger, cfg sqldb.Config, pageNumber string, rowsPerPage string) error {
	db, err := sqldb.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	page, err := strconv.Atoi(pageNumber)
	if err != nil {
		return fmt.Errorf("converting page number: %w", err)
	}

	rows, err := strconv.Atoi(rowsPerPage)
	if err != nil {
		return fmt.Errorf("converting rows per page: %w", err)
	}

	userBus := userbus.NewCore(log, userdb.NewStore(log, db), nil)

	users, err := userBus.Query(ctx, userbus.QueryFilter{}, userbus.DefaultOrderBy, page, rows)
	if err != nil {
		return fmt.Errorf("retrieve users: %w", err)
	}

	stringify, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}
	fmt.Println(stringify)
	return nil
}
