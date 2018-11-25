package email

type (
	// UserNotification struct
	UserNotification struct{}
)

// NewUserNotification is a factory function
func NewUserNotification() *UserNotification {
	return &UserNotification{}
}

// ResetPasswordInstruction notification
func (n *UserNotification) ResetPasswordInstruction(name, email, token string) error {
	return nil
}
