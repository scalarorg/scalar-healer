package auth

type Handler struct {
	domain string
}

func NewHandler(domain string) *Handler {
	return &Handler{
		domain: domain,
	}
}
