package models

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

// ENUM(active, inactive)
type QuestionTypeStatus string

type QuestionType struct {
	ID        string             `bun:"type:uuid,pk" json:"id"`
	Name      string             `json:"name"`
	Slug      string             `json:"slug"`
	Status    QuestionTypeStatus `json:"status"`
	CreatedAt time.Time          `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time          `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`

	bun.BaseModel `bun:"table:question_types" json:"-"`
}

// ENUM(single_choice, multiple_choice)
type OptionType string

type Question struct {
	ID                string     `bun:"type:uuid,pk" json:"id"`
	QuizID            string     `bun:"type:uuid,notnull" json:"quiz_id"`
	QuestionTypeID    string     `bun:"type:uuid,notnull" json:"question_type_id"`
	Question          string     `json:"question"`
	TimeLimitDuration int        `json:"time_limit_duration"`
	Position          int        `json:"position"`
	OptionType        OptionType `json:"option_type"`
	CreatedAt         time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt         time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`

	QuestionType    *QuestionType    `bun:"rel:belongs-to,join:question_type_id=id" json:"question_type"`
	QuestionOptions []QuestionOption `bun:"rel:has-many,join:id=question_id" json:"options"`

	bun.BaseModel `bun:"table:questions" json:"-"`
}

type QuestionOption struct {
	ID         string    `bun:"type:uuid,pk" json:"id"`
	QuestionID string    `bun:"type:uuid,notnull" json:"question_id"`
	Option     string    `json:"option"`
	IsCorrect  bool      `json:"is_correct"`
	CreatedAt  time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt  time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`

	bun.BaseModel `bun:"table:question_options" json:"-"`
}

type QuestionTypeRepository interface {
	FindAll(context.Context) ([]QuestionType, error)
}
