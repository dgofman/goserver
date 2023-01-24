package constants

import (
	"net/http"
)

/**
 * Constant variables and definitions
 */

const (
	HTTP_TIMEOUT = 10
	SECRET_KEY   = "Em3rg1ngS3rv1c3s!"

	SPECIAL_CHARS = "~!@#$^&*()-_+=|:,./?"

	RESTRICT_KEYS = "'Support-Email', 'Referer', 'Accept-Language', 'Base-Url', 'Host', 'Upgrade-Insecure-Requests', 'X-Forwarded-For'"

	TOKEN_COOKIE_NAME = "EQXTEST.globaledge"
	JWT_NAME          = "ES-JWTAuth"
	APP_NAME          = "goserver"

	STATUS_NEW          = "new"          //Set at the moment of user creation.
	STATUS_INITIAL      = "initial"      //Set at the moment user accepted license agreement.
	STATUS_STEP1        = "step1"        //Set once the user has a single IP in whitelists
	STATUS_NO_USER      = "no_user"      //Set if the user doesn't  exist
	STATUS_FOREIGN_USER = "foreign_user" //Set if the user comes to us without Ping headers.

	URL_Static         = "/static/"
	URL_Login          = "/login"
	URL_Logout         = "/logout"
	URL_AddClient      = "/addclient"
	URL_ManageDevices  = "/manage_devices"
	URL_ManageLicenses = "/manage_licenses"
	URL_Map            = "/map"
	URL_Install        = "/installation"
	URL_UserProfile    = "/userProfile"
	URL_DeviceInfo     = "/device_info"
	URL_AjaxIpdata     = "/ajax_ipdata"
	URL_AjaxMetrics    = "/ajax_metrics"
	URL_AjaxStatus     = "/status"
	URL_LoadURL        = "/load"
	URL_GetStarted     = "/$"

	URL_CreateAccount   = "/public/createAccount"
	URL_ActivateAccount = "/public/activateAccount"
	URL_ForgotUsername  = "/public/forgotUsername"
	URL_ResetPassword   = "/public/resetPassword"
	URL_UpdateEmail     = "/public/updateEmail"
	URL_VerifyEmail     = "/public/verify_email"
	URL_LinkedInOauth   = "/public/oauth"
	URL_DeleteAccount   = "/public/deleteAccount"

	Authenticate      = "/api/auth_user"
	ApiStatus         = "/api/status"
	ApiCreateUser     = "/api/create_user"
	ApiDeviceInfo     = "/api/device_info"
	ApiIpData         = "/api/all_ips_by_device"
	ApiMetrics        = "/api/metrics"
	API_DeleteAccount = "/api/delete_user"

	ProxyInitial     = "/proxy/initial"
	ProxyVerifyEmail = "/proxy/verify"
	ProxyCreate      = "/proxy/create"
	ProxyUsername    = "/proxy/username"
	ProxyPassword    = "/proxy/password"
	ProxyActivate    = "/proxy/activate"
	ProxyProfile     = "/proxy/user_profile"

	DEFAULT_HTTP_ERROR                   = "An unexpected error has occurred."
	EQX_ESE_011_EMAIL_MISSING_IN_PROFILE = "EQX-ESE-011"

	LI_client_id              = "869884ejydr7k7"
	LI_client_secret          = "7RX3UfK2254hTv2P"
	LI_scope_permissions      = "r_basicprofile r_liteprofile r_emailaddress"
	LI_authorization_base_url = "https://www.linkedin.com/uas/oauth2/authorization"
	LI_token_url              = "https://www.linkedin.com/uas/oauth2/accessToken"
)

type Error struct {
	Code        string
	Message     string
	Status_Code float64
	Method      string
	Path        string
	Details     string
}

type Properties struct {
	PORT         int
	BASE_DIR     string
	STATIC_DIR   string
	TEMPLATE_DIR string
	VIEW         string
	LOGFILE      bool
	CONSOLE      bool
	ENV          string
	Envs         map[string]interface{}
}

type Response struct {
	Header     http.Header
	Body       map[string]interface{}
	StatusCode int
	Error      *Error
}

type Payload struct {
	Email          interface{} `json:"pawam_email,omitempty"`
	OrgKey         interface{} `json:"pawam_orgid,omitempty"`
	Locale         interface{} `json:"pawam_locale,omitempty"`
	LastName       interface{} `json:"pawam_lastname,omitempty"`
	FirstName      interface{} `json:"pawam_firstname,omitempty"`
	UserName       interface{} `json:"pawam_username,omitempty"`
	Authkey        interface{} `json:"pawam_oat,omitempty"`
	Avatar         interface{} `json:"pawam_avatar,omitempty"` //'pawam_avatar' or 'pawam_pftoken' (PING)
	Subject        interface{} `json:"sub,omitempty"`
	IssuedAt       int64       `json:"iat,omitempty"`
	Issuer         string      `json:"iss,omitempty"`
	Audience       interface{} `json:"aud,omitempty"`
	ExpirationTime int64       `json:"exp,omitempty"`
	AppName        interface{} `json:"acr,omitempty"`
}
