package entity

import (
	"fmt"

	"github.com/konnovK/superchat/internal/utils"
)

func (ccr *CreateChatRequest) Validate() (*utils.ValidationFields, error) {
	fields := utils.ValidationFields{}

	if ccr.Title == "" {
		fields = append(fields, utils.ValidationField{
			Name:  "title",
			Error: "field shouldn't be empty",
		})
	}

	if ccr.Creator == "" {
		fields = append(fields, utils.ValidationField{
			Name:  "creator",
			Error: "field shouldn't be empty",
		})
	}

	if ccr.TTL == 0 {
		fields = append(fields, utils.ValidationField{
			Name:  "ttl",
			Error: "field shouldn't be empty",
		})
	}

	if len(fields) == 0 {
		return &fields, nil
	}
	return &fields, fmt.Errorf("validation error")
}
