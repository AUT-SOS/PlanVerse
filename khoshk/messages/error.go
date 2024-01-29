package messages

const InternalError = "Internal server error"
const InvalidRequestBody = "Invalid request body"
const TooManyRequests = "Too many requests"
const DuplicateEmail = "There is an other account existing with this email"
const FailedPasswordHashGeneration = "Failed to generate password hash"
const FailedToCreateUser = "Failed to create user"
const FailedToCreateProject = "Failed to create project"
const FailedToGenerateLink = "Failed to generate invitation link"
const PasswordIncorrect = "Password is incorrect"
const UserNotVerified = "User didn't verified"
const AlreadyMember = "Already member of project"
const Unauthorized = "Unauthorized"
const InvalidRefreshToken = "Invalid refresh token"
const InvalidAccessToken = "Access token is invalid or its expired"
const FailedToCreateCode = "Failed to create verify code"
const WrongOTP = "Wrong verification code"
const OTPExpired = "Verification code expired"
const FailedToCreateAccessToken = "Failed to create access token"
const FailedToCreateRefreshToken = "Failed to create refresh token"
const RefreshTokenExpired = "Refresh token expired"
const NewAccessToken = "New access token generated"
const WrongEmail = "Email doesn't exist"
const WrongUserID = "Wrong user ID"
const UserNoProject = "This user doesn't has any project"
const WrongProjectID = "Wrong project ID"
const WrongStateID = "Wrong state ID"
const Uninvited = "This user is not allowed to join the project"
const AlreadyJoined = "This user is the member of the project already"
const AdminAccess = "Only admins of project have access to this action"
const AdminChange = "Only admins of project can change this state"
const NotMember = "This user is not a member of project"
const OwnerAccess = "Only owner of project have access to this action"
const AlreadyAdmin = "This user is an admin of the project already"
const AlreadyMemberRole = "This user has member role already"
const OwnerChange = "Owner can't change its role"
const StateNotInProject = "This state is not in this project"
const TaskNotInProject = "This task is not in this project"
const WrongTaskID = "Wrong task ID"
const InTask = "This task has been assigned to this user already"
const NotInTask = "This task hasn't been assigned to this user"
