package model

type User struct {
	ID int `json:"id" db:"id"`
	Sno string	`json:"sno" db:"sno" binding:"required"`
	Phone string `json:"phone" db:"phone" binding:"required"`
	Password string `json:"password" db:"password"`
	State string `json:"state" db:"state"`
	Ban string `json:"ban" db:"ban"`
	Username string `json:"username" db:"username" binding:"required"`
	Email string `json:"email" db:"email" binding:"required,email"`
}

const (
	VerifyUser = "verify_user"
	NormalUser = "normal_user"
	RefuseUser = "refuse_user"
	BlackList = "blacklist"
	VerifyAdmin = "verify_admin"
	NormalAdmin = "normal_admin"
	Root = "root"
)

func AllState() []string {
	return []string{"verify_user", "normal_user", "refuse_user", "blacklist", "verify_admin", "normal_admin"}
}

