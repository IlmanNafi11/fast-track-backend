package domain

type UpdateProfilRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

type ProfilResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
