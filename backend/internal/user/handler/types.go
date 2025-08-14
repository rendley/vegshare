// package handler объявляет, что код принадлежит пакету 'handler'.
package handler

type UpdateProfileRequest struct {
	FullName  string `json:"full_name" validate:"omitempty,min=2,max=100"`
	AvatarURL string `json:"avatar_url" validate:"omitempty,url"`
}

type ProfileResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	AvatarURL string `json:"avatar_url"`
}
