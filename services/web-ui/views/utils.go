package views

import (
	"context"
	"web-ui/internal/models"
)

func GetContextData(ctx context.Context) models.ContextData {
	if data, ok := ctx.Value(models.DataKey).(models.ContextData); ok {
		return data
	}

	return models.ContextData{}
}
