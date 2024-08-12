package forms

type UserSignupForm struct {
	Username string `schema:"username, required" validate:"required,max=100"`
	Email string `schema:"email, required" validate:"required,email"`
	Password string `schema:"password, required" validate:"required,min=8"`
	PasswordConfirm string `schema:"password_confirm, required" validate:"required,eqfield=Password"`
	BaseForm
}

type UserLoginForm struct {
	Email string `schema:"email, required" validate:"required,email"`
	Password string `schema:"password, required" validate:"required,min=8"`
	BaseForm
}

type UserPasswordUpdateForm struct {
	CurrentPassword string `schema:"current_password, required" validate:"required"`
	NewPassword string `schema:"new_password, required" validate:"required,min=8,nefield=CurrentPassword"`
	NewPasswordConfirm string `schema:"new_password_confirm, required" validate:"required,eqfield=NewPassword"`
	BaseForm
}