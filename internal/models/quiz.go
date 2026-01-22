package models

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

// ENUM(public, private)
type QuizVisibility string

type Quiz struct {
	ID          string         `bun:"type:uuid,pk" json:"id"`
	OwnerID     string         `bun:"type:uuid,notnull" json:"owner_id"`
	Title       string         `json:"title"`
	Description *string        `bun:",nullzero" json:"description"`
	Visibility  QuizVisibility `bun:",nullzero" json:"visibility"`
	CoverImage  *string        `bun:",nullzero" json:"cover_image"`
	PublishedAt *time.Time     `bun:",nullzero" json:"published_at"`
	CreatedAt   time.Time      `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time      `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`

	Owner     *User      `bun:"rel:belongs-to,join:owner_id=id" json:"-"`
	Questions []Question `bun:"rel:has-many,join:id=quiz_id" json:"questions"`

	bun.BaseModel `bun:"table:quizzes" json:"-"`
}

type FindQuizOptions struct {
	ID      string
	OwnerID string
}

type ListQuizOptions struct {
	Paginator  Paginator
	Search     string
	OwnerID    string
	Visibility QuizVisibility
}

type QuizRepository interface {
	Create(context.Context, *Quiz) error
	Update(context.Context, *Quiz) error
	FindOne(context.Context, *FindQuizOptions) (*Quiz, error)
	FindAll(context.Context, *ListQuizOptions) ([]Quiz, int64, error)
}

type CreateOrEditQuizRequest struct {
	Title       string `json:"title" valid:"required~The title field is required,maxstringlength(70)"`
	Description string `json:"description" valid:"maxstringlength(500)"`
	Questions   []struct {
		Question          string     `json:"question"`
		TimeLimitDuration int        `json:"time_limit_duration" valid:"int"`
		Position          int        `json:"position" valid:"int,optional"`
		OptionType        OptionType `json:"option_type" valid:"in(single_choice|multiple_choice),optional"`
		Options           []struct {
			Option    string `json:"option"`
			IsCorrect bool   `json:"is_correct"`
		} `json:"options"`
	} `json:"questions"`
	Visibility QuizVisibility `json:"visibility" valid:"required~The visibility field is required,in(public|private)~The visibility field must be public or private"`
	CoverImage string         `json:"cover_image" valid:"optional"`
}
