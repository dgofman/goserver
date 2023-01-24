package lib

import (
	"time"

	"go/goserver.io/constants"
	"go/goserver.io/utils"

	"github.com/gbrlsnchs/jwt"
)

/**
 * Get encryption key for encrypt/decrypt URL parameters
 */
func AuthKey(json map[string]interface{}) string {
	username := utils.Get(json, "UserName")
	if username != nil {
		return Encode(constants.SECRET_KEY, constants.SECRET_KEY+username.(string))
	}
	return constants.SECRET_KEY
}

/**
 * Create and sign JWT token
 */
func CreateJWT(json map[string]interface{}) string {
	pawam_orgid := utils.Get(json, "UserOrg")
	if pawam_orgid != nil {
		m := pawam_orgid.(map[string]interface{})
		pawam_orgid = m["OrgKey"]
	}
	avatar := utils.Get(json, "ExtendedUserAttribute")
	if avatar != nil {
		m := avatar.(map[string]interface{})
		avatar = m["AVATAR"]
	}

	now := time.Now()
	hs := jwt.NewHS256([]byte(Encode(constants.SECRET_TOKEN)))
	p := constants.Payload{
		Email:          utils.Get(json, "Email"),
		OrgKey:         pawam_orgid,
		Locale:         utils.Get(json, "Locale"),
		LastName:       utils.Get(json, "LastName"),
		FirstName:      utils.Get(json, "FirstName"),
		UserName:       utils.Get(json, "UserName"),
		Authkey:        AuthKey(json),
		Avatar:         avatar,
		Subject:        utils.Get(json, "UserName"),
		IssuedAt:       now.Unix(),
		Issuer:         constants.JWT_NAME,
		Audience:       "global",
		ExpirationTime: now.Add(30 * 60 * time.Hour).Unix(),
		AppName:        constants.APP_NAME,
	}
	token, err := jwt.Sign(p, hs)
	if err != nil {
		utils.LogError(err)
	}
	return string(token)
}

/**
 * Decode JWT cookie parameter to JSON/Payload object
 */
func JwtDecode(token string) (*constants.Payload, error) {
	hs := jwt.NewHS256([]byte(Encode(constants.SECRET_TOKEN)))
	p := &constants.Payload{}
	_, err := jwt.Verify([]byte(token), hs, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
