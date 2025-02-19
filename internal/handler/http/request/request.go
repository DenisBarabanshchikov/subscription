package request

type CreateCustomer struct {
	Email string `json:"email"`
}

type SubscribeCustomer struct {
	Plan string `json:"plan"`
}
