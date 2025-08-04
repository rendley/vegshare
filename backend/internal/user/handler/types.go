package handler

type UpdateProfileRequest struct {
	FullName  string `json:"full_name" validate:"omitempty,min=2,max=100"`
	AvatarURL string `json:"avatar_url" validate:"omitempty,url"`
}
