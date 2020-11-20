package model

type User struct {
	ID int `json:"id" db:"id"`
	Sno string	`json:"sno" db:"sno" binding:"required"`
	Phone string `json:"phone" db:"phone" binding:"required"`
	Password string `json:"password" db:"password"`
	State string `json:"state" db:"state"`
	Ban string `json:"ban" db:"ban"`
	Username string `json:"username" db:"username" binding:"required"`
	GroupList string `json:"group_list" db:"group_list"`
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

// StateOptions 不同角色可以管理的状态用户
// 1. 管理员可以管理普通用户
// 2. 只有root才可以管理管理员
func StateOptions(role string) []string {
	state := []string{"verify_user", "normal_user", "refuse_user", "blacklist"}
	if role == "root" {
		state = append(state, "verify_admin", "normal_admin")
	}
	return state
}

