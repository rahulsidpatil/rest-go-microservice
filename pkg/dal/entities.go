package dal

//Message ... message object
type Message struct {
	ID      int    `json:"id,omitempty" example:"1"`
	Message string `json:"message,omitempty" example:"it is what is it"`
}
