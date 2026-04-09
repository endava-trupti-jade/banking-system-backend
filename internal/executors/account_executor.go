package executors

import (
	"banking-system-backend/constants"
	_ "banking-system-backend/internal/models"
	"banking-system-backend/internal/repositories/mongorepo"
	"context"
	_ "encoding/json"

	_ "go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountExecutor struct {
	accountRepo *mongorepo.AccountRepository
}

func (e *AccountExecutor) Execute(ctx context.Context, action constants.Action, payload []byte) error {
	// switch action {
	// case constants.ActionUpdate:
	// 	var account models.Account
	// 	if err := json.Unmarshal(payload, &account); err != nil {
	// 		return err
	// 	}
	// 	return nil, //e.accountRepo.Update(ctx, account.ID, account)
	// case constants.ActionDelete:
	// 	var data =  struct {
	// 		ID primitive.ObjectID `json:"id"`
	// 	}
	// 	if err := json.Unmarshal(payload, &data); err != nil {
	// 		return err
	// 	}
	// 	return nil //e.accountRepo.Delete(ctx, data.ID)
	// default:
	// 	return constants.ErrInvalidAction
	// }
	return nil
}
