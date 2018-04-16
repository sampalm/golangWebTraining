package webapp

const (
	FirstNameEmpty          = 05
	FirstNameInvalid        = 06
	FirstNameInvalidOrEmpty = 07

	EmailEmpty       = 10
	EmailInvalid     = 11
	EmailAlreadyUsed = 12

	PasswordEmpty          = 20
	PasswordDoNotMatch     = 21
	PasswordInvalid        = 22
	PasswordInvalidOrEmpty = 23

	PasswordEmailErr = 25

	AccountDoNotExists = 30
	AccountInvalid     = 31

	UnauthorizedAccess = 40
)

var statusText = map[int]string{
	FirstNameEmpty:          "Firstname field is empty",
	FirstNameInvalid:        "Firstname field value is invalid",
	FirstNameInvalidOrEmpty: "Firstname field value is invalid or empty",

	EmailEmpty:       "Email field is empty",
	EmailInvalid:     "Email field value is invalid",
	EmailAlreadyUsed: "Email already been used",

	PasswordEmpty:          "Password field is empty",
	PasswordDoNotMatch:     "Password and Confirm Password field doesn't match",
	PasswordInvalid:        "Password field value is invalid",
	PasswordInvalidOrEmpty: "Password field is empty or invalid",

	PasswordEmailErr: "Password or Email is empty or invalid",

	AccountInvalid:     "This account reported is invalid",
	AccountDoNotExists: "This account reported doesn't exist",

	UnauthorizedAccess: "You don't have access to this page",
}

// StatusText returns a text for the given request code.
func StatusCodeText(code int) string {
	return statusText[code]
}
