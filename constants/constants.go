package constants

// import pbt "github.com/Allen-Career-Institute/common-protos/user_management/v1/types"
//
// const AwsRegionKey = "AWS_REGION"
// const StudentsDocumentS3BucketNameKey = "STUDENTS_DOCUMENT_S3_BUCKET_NAME"
// const PreSignedURLTTLMinutesKey = "PRE_SIGNED_URL_TTL_MINUTES"
// const OtelCollectorEndpoint = "otel-collector.default.svc.cluster.local:4317"
// const UniqueConstrainstError = "Unique constraint violated for "
// const ExternalIDUniqueIndex = "users.external_user_id"
// const RetryAttempts = 3
//
// const InvalidIdentityErrorMessage = "Invalid identity or credentials."
//
// const DefaultSortBy = "created_at"
// const DefaultOrderBy = "desc"
//
// const PhoneNumber = "phone_number"
// const Email = "email"
// const EmployeeID = "employee_id"
// const IdentityValue = "identity_value"
// const IdentityType = "identity_type"
//
// const IdentityTypePhone = "PHONE"
// const PersonaTypeParent = "PARENT"
//
// const UserIdIdentityMapKey = "%s#%s#%s" // userId#ersonaType#identityType
//
// const ID = "id"
// const PersonaType = "persona_type"
// const UserStatus = "status"
// const TenantID = "tenant_id"
// const UserStage = "stage"
// const UserID = "user_id"
// const SystemUser = "system"
// const UserInfoCacheTTLInHoursKey = "USER_INFO_CACHE_TTL_IN_HOURS"
// const UserInfoCacheEnabledKey = "USER_INFO_CACHE_ENABLED"
// const DefaultUserInfoCacheTTLInHours = 3
// const IsVisible = "is_visible"
// const EmptyString = ""
//
//	var SortByMap = map[string]string{
//		"CreatedAt": "created_at",
//		"UpdatedAt": "updated_at",
//		"Asc":       "asc",
//		"Desc":      "desc",
//	}
//
// var ValidStudentDocumentFileTypes = []string{".jpeg", ".png", ".pdf", ".jpg", "image/jpeg", "image/png", "application/pdf"}
//
// var CommunicableIdentityTypes = []string{pbt.IdentityType_PHONE.String(), pbt.IdentityType_EMAIL.String(), pbt.IdentityType_WHATSAPP.String()}
//
// const StudentDocumentsMinSizeInBytes = 5
// const StudentDocumentsMaxSizeInBytes = 5 * 1024 * 1024
//
// const UserEntityName = "User"
// const UserDeleteEvent = "DELETE_USER"
// const Delete = "DELETE"
//
// const MinCharsFirstName = 3
//
// const UnimplementedMethodError = "implement me"
//
// const UserNotFoundError = "users not found with due to error :: %s"
//
// const (
//
//	ErrUnauthorized    = "Unauthorized"
//	ErrUserNotLoggedIn = "User not logged in"
//
// )
//
// const (
//
//	UnauthorizedErrorCode = 401
//
// )
//
// const (
//
//	DbEmailKey             = "email"
//	DbStatusKey            = "status"
//	DbFirstNameKey         = "first_name"
//	DbMiddleNameKey        = "middle_name"
//	DbLastNameKey          = "last_name"
//	DbGenderKey            = "gender"
//	DbProfilePhotoUrlIdKey = "profile_photo_url_id"
//	DbDobKey               = "date_of_birth"
//	DbPhoneNumberKey       = "phone_number"
//	DbUpdateAtKey          = "updated_at"
//	DbUpdateByKey          = "updated_by"
//
// )
//
// const (
//
//	DbParentUserIDKey = "parent_user_id"
//	DbFatherUserIDKey = "father_user_id"
//	DbMotherUserIDKey = "mother_user_id"
//
// )
//
// const (
//
//	ErrLocationNotFound          = "Location with pincode %s not found"
//	ErrGetLocationFailed         = "GetLocationByPincode failed due to error %v"
//	ErrBulkCreateLocationsFailed = "BulkCreateLocations failed due to error %v"
//	ErrRecordNotFound            = "record not found"
//	ErrDeleteIdentities          = "Error deleting identity for user %v, Err - %v"
//	ErrDeleteCreds               = "Error deleting credentials for user : %v, Err - %v"
//	ErrUserInfoCacheDisabled     = "User info cache is disabled"
//	UserIdMissing                = "User ID is missing"
//	OnlyActiveAddress            = "Only 1 active address present,This can not be deleted"
//
// )
//
// const (
//
//	ShouldFetchDataFromIdentitiesKey = "SHOULD_FETCH_DATA_FROM_IDENTITIES"
//
// )
//
// const (
//
//	FormIDRegex = `^[a-zA-Z0-9_]+$`
//	PhoneRegex  = `^[1-9]{1}[0-9]{9}$`
//	EmailRegex  = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`
//
// )
//
// const (
//
//	InvalidIdentityError = "invalid identity value"
//
// )
//
// const JwtSecretKey = "jwtSecret"
//
// const ServiceName = "user_management"
// const CacheEntityUserKey = "user"
// const Profile = "profile"
// const Address = "address"
// const DefaultCacheTTLInHours = 168
const DefaultCacheEntryCount = 5

//
//const CredentialHashingCost = 11
//
//const SpecialCharsSpace = " "
//
//const MysqlUniqueConstraintViolationErrorCode = 1062
//
//const AddressHasActiveOrderWarningMessage = "Address is associated with %v active orders."
//
//const DefaultDateOfBirth = "01-01-2024"
