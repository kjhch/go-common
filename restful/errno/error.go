package errno

type Error struct {
	Code    string
	Message string
	Status  int
}

var ErrMap = map[string]Error{
	OK:        {Code: "00000", Message: "OK", Status: 200},
	ClientErr: {Code: "A0001", Message: "Client Err", Status: 400},

	RegistrationErr:            {Code: "A0100", Message: "Registration Err", Status: 400},
	PrivacyPolicyNotAccepted:   {Code: "A0101", Message: "Privacy Policy Not Accepted", Status: 400},
	RegistrationAreaRestricted: {Code: "A0102", Message: "Registration Area Restricted", Status: 400},

	InvalidUsername:       {Code: "A0110", Message: "Invalid Username", Status: 400},
	ExistedUsername:       {Code: "A0111", Message: "Existed Username", Status: 400},
	SensitiveUsername:     {Code: "A0112", Message: "Sensitive Username", Status: 400},
	IllegalCharInUsername: {Code: "A0113", Message: "Illegal Char In Username", Status: 400},

	InvalidPassword: {Code: "A0120", Message: "Invalid Password", Status: 400},
	ShortPassword:   {Code: "A0121", Message: "Short Password", Status: 400},
	WeakPassword:    {Code: "A0122", Message: "Weak Password", Status: 400},

	IncorrectVerificationCode:      {Code: "A0130", Message: "Incorrect Verification Code", Status: 400},
	IncorrectSMSVerificationCode:   {Code: "A0131", Message: "Incorrect SMS Verification Code", Status: 400},
	IncorrectEmailVerificationCode: {Code: "A0132", Message: "Incorrect Email Verification Code", Status: 400},
	IncorrectVoiceVerificationCode: {Code: "A0133", Message: "Incorrect Voice Verification Code", Status: 400},

	InvalidDocumentation:   {Code: "A0140", Message: "Invalid Documentation", Status: 400},
	DocumentationNotChosen: {Code: "A0141", Message: "Documentation Not Chosen", Status: 400},
	InvalidIDCardNum:       {Code: "A0142", Message: "Invalid ID Card Num", Status: 400},
	InvalidPassportNum:     {Code: "A0143", Message: "Invalid Passport Num", Status: 400},
	InvalidMilitaryID:      {Code: "A0144", Message: "Invalid Military ID", Status: 400},

	InvalidBasicInfo: {Code: "A0150", Message: "Invalid Basic Info", Status: 400},
	InvalidPhoneNum:  {Code: "A0151", Message: "Invalid Phone Num", Status: 400},
	InvalidAddress:   {Code: "A0152", Message: "Invalid Address", Status: 400},
	InvalidEmail:     {Code: "A0153", Message: "Invalid Email", Status: 400},

	AbnormalLogin: {Code: "A0200", Message: "Abnormal Login", Status: 400},
	UserNotFound:  {Code: "A0201", Message: "User Not Found", Status: 404},
	UserFrozen:    {Code: "A0202", Message: "User Frozen", Status: 400},
	UserInvalid:   {Code: "A0203", Message: "User Invalid", Status: 400},

	IncorrectPassword:        {Code: "A0210", Message: "Incorrect Password", Status: 400},
	TooManyIncorrectPassword: {Code: "A0211", Message: "Too Many Incorrect Password", Status: 403},

	IdentityUnrecognized:        {Code: "A0220", Message: "Identity Unrecognized", Status: 400},
	FingerprintUnrecognized:     {Code: "A0221", Message: "Fingerprint Unrecognized", Status: 400},
	FaceUnrecognized:            {Code: "A0222", Message: "Face Unrecognized", Status: 400},
	ThirdPartyLoginUnauthorized: {Code: "A0223", Message: "Third Party Login Unauthorized", Status: 401},

	LoginExpired: {Code: "A0230", Message: "Login Expired", Status: 401},

	WrongCaptcha:        {Code: "A0240", Message: "Wrong Captcha", Status: 400},
	TooManyWrongCaptcha: {Code: "A0241", Message: "Too Many Wrong Captcha", Status: 403},

	AccessRightAbnormal: {Code: "A0300", Message: "Access Right Abnormal", Status: 400},
	AccessUnauthorized:  {Code: "A0301", Message: "Access Unauthorized", Status: 401},
	Authorizing:         {Code: "A0302", Message: "Authorizing", Status: 401},
	AuthorizationDenied: {Code: "A0303", Message: "Authorization Denied", Status: 403},

	AccessRestrictedPrivacy: {Code: "A0310", Message: "Access Restricted Privacy", Status: 403},
	AuthorizationExpired:    {Code: "A0311", Message: "Authorization Expired", Status: 401},
	APIForbidden:            {Code: "A0312", Message: "API Forbidden", Status: 403},

	AccessDenied:            {Code: "A0320", Message: "Access Denied", Status: 403},
	BlacklistUser:           {Code: "A0321", Message: "Blacklist User", Status: 403},
	AccountFrozen:           {Code: "A0322", Message: "Account Frozen", Status: 403},
	IllegalIP:               {Code: "A0323", Message: "Illegal IP", Status: 403},
	GatewayAccessRestricted: {Code: "A0324", Message: "Gateway Access Restricted", Status: 403},
	BlacklistArea:           {Code: "A0325", Message: "Blacklist Area", Status: 403},

	OutOfCredit: {Code: "A0330", Message: "Out Of Credit", Status: 400},

	SignatureAbnormal:    {Code: "A0340", Message: "Signature Abnormal", Status: 400},
	RSASignatureAbnormal: {Code: "A0341", Message: "RSA Signature Abnormal", Status: 400},

	WrongParam:         {Code: "A0400", Message: "Wrong Param", Status: 400},
	MaliciousLinkFound: {Code: "A0401", Message: "Malicious Link Found", Status: 400},
	InvalidInput:       {Code: "A0402", Message: "Invalid Input", Status: 400},

	RequiredParamEmpty: {Code: "A0410", Message: "Required Param Empty", Status: 400},
	OrderNumEmpty:      {Code: "A0411", Message: "Order Num Empty", Status: 400},
	OrderQuantityEmpty: {Code: "A0412", Message: "Order Quantity Empty", Status: 400},
	TimeParamEmpty:     {Code: "A0413", Message: "Time Param Empty", Status: 400},
	IllegalTimeParam:   {Code: "A0414", Message: "Illegal Time Param", Status: 400},

	ParamOutOfRange:       {Code: "A0420", Message: "Param Out Of Range", Status: 400},
	ParamFormatUnexpected: {Code: "A0421", Message: "Param Format Unexpected", Status: 400},
	AddressOutOfRange:     {Code: "A0422", Message: "Address Out Of Range", Status: 400},
	TimeOutOfRange:        {Code: "A0423", Message: "Time Out Of Range", Status: 400},
	AmountOutOfRange:      {Code: "A0424", Message: "Amount Out Of Range", Status: 400},
	NumberOutOfRange:      {Code: "A0425", Message: "Number Out Of Range", Status: 400},
	TotalBatchOutOfRange:  {Code: "A0426", Message: "Total Batch Out Of Range", Status: 400},
	JSONParseErr:          {Code: "A0427", Message: "JSON Parse Err", Status: 400},

	IllegalContent:        {Code: "A0430", Message: "Illegal Content", Status: 400},
	IllegalWord:           {Code: "A0431", Message: "Illegal Word", Status: 400},
	IllegalPicture:        {Code: "A0432", Message: "Illegal Picture", Status: 400},
	CopyrightInfringement: {Code: "A0433", Message: "Copyright Infringement", Status: 400},

	OperationAbnormal:         {Code: "A0440", Message: "Operation Abnormal", Status: 400},
	PaymentOvertime:           {Code: "A0441", Message: "Payment Overtime", Status: 400},
	OrderConfirmationOvertime: {Code: "A0442", Message: "Order Confirmation Overtime", Status: 400},
	OrderCancelled:            {Code: "A0443", Message: "Order Cancelled", Status: 400},

	UserRequestAbnormal:       {Code: "A0500", Message: "User Request Abnormal", Status: 400},
	TooManyRequests:           {Code: "A0501", Message: "Too Many Requests", Status: 429},
	TooManyRequestConcurrency: {Code: "A0502", Message: "Too Many Request Concurrency", Status: 429},
	OperationWait:             {Code: "A0503", Message: "Operation Wait", Status: 400},
	WebsocketConnAbnormal:     {Code: "A0504", Message: "Websocket Conn Abnormal", Status: 400},
	WebsocketDisconnected:     {Code: "A0505", Message: "Websocket Disconnected", Status: 400},
	DuplicateRequest:          {Code: "A0506", Message: "Duplicate Request", Status: 400},

	UserResourceAbnormal: {Code: "A0600", Message: "User Resource Abnormal", Status: 400},
	InsufficientBalance:  {Code: "A0601", Message: "Insufficient Balance", Status: 400},
	InsufficientDisk:     {Code: "A0602", Message: "Insufficient Disk", Status: 400},
	InsufficientRAM:      {Code: "A0603", Message: "Insufficient RAM", Status: 400},
	InsufficientOSSSpace: {Code: "A0604", Message: "Insufficient OSS Space", Status: 400},
	OutOfQuota:           {Code: "A0605", Message: "Out Of Quota", Status: 400},

	FileUploadAbnormal:     {Code: "A0700", Message: "File Upload Abnormal", Status: 400},
	UnexpectedFileType:     {Code: "A0701", Message: "Unexpected File Type", Status: 400},
	FileTooLarge:           {Code: "A0702", Message: "File Too Large", Status: 413},
	PictureTooLarge:        {Code: "A0703", Message: "Picture Too Large", Status: 413},
	VideoTooLarge:          {Code: "A0704", Message: "Video Too Large", Status: 413},
	CompressedFileTooLarge: {Code: "A0705", Message: "Compressed File Too Large", Status: 413},

	VersionAbnormal:      {Code: "A0800", Message: "Version Abnormal", Status: 400},
	VersionMismatched:    {Code: "A0801", Message: "Version Mismatched", Status: 400},
	VersionTooLow:        {Code: "A0802", Message: "Version Too Low", Status: 400},
	VersionTooHigh:       {Code: "A0803", Message: "Version Too High", Status: 400},
	VersionOutdated:      {Code: "A0804", Message: "Version Outdated", Status: 400},
	UnexpectedAPIVersion: {Code: "A0805", Message: "Unexpected API Version", Status: 400},
	APIVersionTooHigh:    {Code: "A0806", Message: "API Version Too High", Status: 400},
	APIVersionTooLow:     {Code: "A0807", Message: "API Version Too Low", Status: 400},

	UserPrivacyUnauthorized: {Code: "A0900", Message: "User Privacy Unauthorized", Status: 401},
	UserPrivacyUnsigned:     {Code: "A0901", Message: "User Privacy Unsigned", Status: 401},
	WebcamUnauthorized:      {Code: "A0902", Message: "Webcam Unauthorized", Status: 401},
	CameraUnauthorized:      {Code: "A0903", Message: "Camera Unauthorized", Status: 401},
	AlbumUnauthorized:       {Code: "A0904", Message: "Album Unauthorized", Status: 401},
	FileUnauthorized:        {Code: "A0905", Message: "File Unauthorized", Status: 401},
	LocationUnauthorized:    {Code: "A0906", Message: "Location Unauthorized", Status: 401},
	ContactsUnauthorized:    {Code: "A0907", Message: "Contacts Unauthorized", Status: 401},

	UserDeviceAbnormal: {Code: "A1000", Message: "User Device Abnormal", Status: 400},
	CameraAbnormal:     {Code: "A1001", Message: "Camera Abnormal", Status: 400},
	MicrophoneAbnormal: {Code: "A1002", Message: "Microphone Abnormal", Status: 400},
	ReceiverAbnormal:   {Code: "A1003", Message: "Receiver Abnormal", Status: 400},
	SpeakerAbnormal:    {Code: "A1004", Message: "Speaker Abnormal", Status: 400},
	GPSAbnormal:        {Code: "A1005", Message: "GPS Abnormal", Status: 400},

	SystemErr: {Code: "B0001", Message: "System Err", Status: 500},

	SystemExecTimeout:   {Code: "B0100", Message: "System Exec Timeout", Status: 500},
	OrderProcessTimeout: {Code: "B0101", Message: "Order Process Timeout", Status: 500},

	DisasterRecovery: {Code: "B0200", Message: "Disaster Recovery", Status: 503},

	SystemRateLimited: {Code: "B0210", Message: "System Rate Limited", Status: 503},

	SystemFunctionDowngrade: {Code: "B0220", Message: "System Function Downgrade", Status: 503},

	SystemResourceAbnormal: {Code: "B0300", Message: "System Resource Abnormal", Status: 500},

	SystemResourceExhausted: {Code: "B0310", Message: "System Resource Exhausted", Status: 500},
	DiskSpaceExhausted:      {Code: "B0311", Message: "Disk Space Exhausted", Status: 500},
	RAMExhausted:            {Code: "B0312", Message: "RAM Exhausted", Status: 500},
	FileHandleExhausted:     {Code: "B0313", Message: "File Handle Exhausted", Status: 500},
	ConnPoolExhausted:       {Code: "B0314", Message: "Conn Pool Exhausted", Status: 500},
	ThreadPoolExhausted:     {Code: "B0315", Message: "Thread Pool Exhausted", Status: 500},

	SystemResourceAccessAbnormal: {Code: "B0320", Message: "System Resource Access Abnormal", Status: 500},
	FileAccessErr:                {Code: "B0321", Message: "File Access Err", Status: 500},

	ThirdPartyErr: {Code: "C0001", Message: "Third Party Err", Status: 500},

	MiddlewareServiceErr: {Code: "C0100", Message: "Middleware Service Err", Status: 500},

	RPCServiceErr:           {Code: "C0110", Message: "RPC Service Err", Status: 500},
	RPCServiceNotFound:      {Code: "C0111", Message: "RPC Service Not Found", Status: 500},
	RPCServiceNotRegistered: {Code: "C0112", Message: "RPC Service Not Registered", Status: 500},
	InterfaceNotFound:       {Code: "C0113", Message: "Interface Not Found", Status: 500},

	MQServiceErr:      {Code: "C0120", Message: "MQ Service Err", Status: 500},
	MQDeliveryErr:     {Code: "C0121", Message: "MQ Delivery Err", Status: 500},
	MQConsumptionErr:  {Code: "C0122", Message: "MQ Consumption Err", Status: 500},
	MQSubscriptionErr: {Code: "C0123", Message: "MQ Subscription Err", Status: 500},
	MQGroupNotFound:   {Code: "C0124", Message: "MQ Group Not Found", Status: 500},

	CacheServiceErr:            {Code: "C0130", Message: "Cache Service Err", Status: 500},
	CacheKeyLengthTooLong:      {Code: "C0131", Message: "Cache Key Length Too Long", Status: 500},
	CacheValueLengthTooLong:    {Code: "C0132", Message: "Cache Value Length Too Long", Status: 500},
	FullCacheStorage:           {Code: "C0133", Message: "Full Cache Storage", Status: 500},
	UnsupportedCacheDataFormat: {Code: "C0134", Message: "Unsupported Cache Data Format", Status: 500},

	ConfigServiceErr: {Code: "C0140", Message: "Config Service Err", Status: 500},

	NetworkServiceErr: {Code: "C0150", Message: "Network Service Err", Status: 500},
	VPNServiceErr:     {Code: "C0151", Message: "VPN Service Err", Status: 500},
	CDNServiceErr:     {Code: "C0152", Message: "CDN Service Err", Status: 500},
	DNSServiceErr:     {Code: "C0153", Message: "DNS Service Err", Status: 500},
	GatewayServiceErr: {Code: "C0154", Message: "Gateway Service Err", Status: 500},

	ThirdPartyServiceTimeout: {Code: "C0200", Message: "Third Party Service Timeout", Status: 500},

	RPCServiceTimeout: {Code: "C0210", Message: "RPC Service Timeout", Status: 500},

	MQDeliveryTimeout: {Code: "C0220", Message: "MQ Delivery Timeout", Status: 500},

	CacheServiceTimeout: {Code: "C0230", Message: "Cache Service Timeout", Status: 500},

	ConfigServiceTimeout: {Code: "C0240", Message: "Config Service Timeout", Status: 500},

	DatabaseServiceTimeout: {Code: "C0250", Message: "Database Service Timeout", Status: 500},

	DatabaseServiceErr: {Code: "C0300", Message: "Database Service Err", Status: 500},

	DBTableNotFound:  {Code: "C0311", Message: "DB Table Not Found", Status: 500},
	DBColumnNotFound: {Code: "C0312", Message: "DB Column Not Found", Status: 500},

	DuplicateColumnJoined: {Code: "C0321", Message: "Duplicate Column Joined", Status: 500},

	DBDeadlock: {Code: "C0331", Message: "DB Deadlock", Status: 500},

	DBPrimaryKeyConflict: {Code: "C0341", Message: "DB Primary Key Conflict", Status: 500},

	ThirdPartyDisasterRecovery:  {Code: "C0400", Message: "Third Party Disaster Recovery", Status: 500},
	ThirdPartyRateLimited:       {Code: "C0401", Message: "Third Party Rate Limited", Status: 500},
	ThirdPartyFunctionDowngrade: {Code: "C0402", Message: "Third Party Function Downgrade", Status: 500},

	NotificationServiceErr:      {Code: "C0500", Message: "Notification Service Err", Status: 500},
	SMSNotificationServiceErr:   {Code: "C0501", Message: "SMS Notification Service Err", Status: 500},
	VoiceNotificationServiceErr: {Code: "C0502", Message: "Voice Notification Service Err", Status: 500},
	EmailNotificationServiceErr: {Code: "C0503", Message: "Email Notification Service Err", Status: 500},
}
