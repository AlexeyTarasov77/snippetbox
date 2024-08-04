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