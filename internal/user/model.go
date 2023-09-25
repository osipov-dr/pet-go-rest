package user

type User struct {
	Id           string `json:"id" bson:"_id,omitempty"`
	Name         string `json:"name" bson:"name"`
	PasswordHash string `json:"-" bson:"password"`
	Email        string `json:"email" bson:"email"`
}

type Dto struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
