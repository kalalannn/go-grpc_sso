package messages

const InvalidPassword = "Invalid password format (8 <= password <= 72 AND one <lowercase> AND one <UPPERCASE> AND one <specsym>)"

const InvalidEmail = "Invalid email format (valid email AND email lenght <= 64)"

const InvalidUsername = "Invalid username format a..z A..Z 0..9 -_"

const EmailAlreadyExists = "User with this email already exists"

const UsernameAlreadyExists = "User with this username already exists"

const UnknownError = "Unknown error"

const UsernameOrEmailIsRequired = "Either username or email is required"

const AuthNotFound = "User with this credentials not found"

const Unauthenticated = "Invalid token"
