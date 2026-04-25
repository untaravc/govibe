package postcontroller

import "time"

type createPostRequest struct {
	Title      string     `json:"title" validate:"required,min=2,max=191"`
	Subtitle   *string    `json:"subtitle" validate:"omitempty,max=255"`
	Content    *string    `json:"content" validate:"omitempty"`
	Status     uint8      `json:"status" validate:"required,oneof=0 1"`
	ImageURL   *string    `json:"image_url" validate:"omitempty,max=512"`
	ReleaseAt  *time.Time `json:"release_at" validate:"omitempty"`
	CategoryID *uint      `json:"category_id" validate:"omitempty"`
}

type updatePostRequest struct {
	Title      *string    `json:"title" validate:"omitempty,min=2,max=191"`
	Subtitle   *string    `json:"subtitle" validate:"omitempty,max=255"`
	Content    *string    `json:"content" validate:"omitempty"`
	Status     *uint8     `json:"status" validate:"omitempty,oneof=0 1"`
	ImageURL   *string    `json:"image_url" validate:"omitempty,max=512"`
	ReleaseAt  *time.Time `json:"release_at" validate:"omitempty"`
	CategoryID *uint      `json:"category_id" validate:"omitempty"`
}
