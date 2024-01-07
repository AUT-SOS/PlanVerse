package messages

const InternalError = "Internal server error"
const InvalidRequestBody = "Invalid request body"
const TooManyRequests = "Too many requests"

const UserExistence = "This username doesn't exist"
const DuplicateUsername = "Username has been taken already"
const DuplicateEmail = "There is an other account existing with this email"
const FailedToSendEmail = "Failed to send email"
const FailedPasswordHashGeneration = "Failed to generate password hash"
const FailedToCreateUser = "Failed to create user"
const UsernameOrPasswordIncorrect = "Username or password is incorrect"
const FailedToCreateToken = "Failed to create token"
const Unauthorized = "Unauthorized"
const InvalidToken = "Invalid token"
const FailedToCreateCode = "Failed to create verify code"
const EmailNotVerified = "Email didn't verified"
const WrongCode = "Wrong verification code"
