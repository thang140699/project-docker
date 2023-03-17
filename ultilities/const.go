package ultilities

const (
	REGPASS        = `?=^.{8,}$)(?=.\d)(?=.[a-z])(?=.[A-Z])(?!.\s)[0-9a-zA-Z!@#$%^&()]*$`
	RegPhoneNumber = `(?m)(84|0[0-9])+([0-9]{8})\b`
	ReggName       = `^[a-zA-Z]*$`
)
