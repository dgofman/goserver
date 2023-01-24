package view

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"go/goserver.io/constants"
	"go/goserver.io/lib"
	"go/goserver.io/utils"
)

/**
 * Template handler for request: /public/createAccount
 */
func CreateAccount(w http.ResponseWriter, r *http.Request, app *lib.App) {
	log.Printf("view:CreateAccount: " + r.Method)
	ctx := lib.CreateContent(r)
	form := userForm(r, map[string]interface{}{
		"SPECIAL_CHARS": constants.SPECIAL_CHARS,
	})
	ctx["form"] = form.GetFields
	ctx["portal_oauth"] = constants.URL_LinkedInOauth

	if r.Method == "POST" {
		err := ValidatePassword(r.FormValue("Password"), r.FormValue("Confirm"))
		if err != nil {
			ctx["error"] = err.Error
		} else {
			profileType := r.FormValue("ProfileType")
			j := map[string]interface{}{
				"Email":     lib.FormValue(r, "Email"),
				"UserName":  lib.FormValue(r, "UserName"),
				"FirstName": lib.FormValue(r, "FirstName"),
				"LastName":  lib.FormValue(r, "LastName"),
				"Password":  lib.FormValue(r, "Password"),
				"Confirm":   lib.FormValue(r, "Confirm"),
			}
			if profileType == "LINKEDIN" {
				j["ExtendedUserAttribute"] = map[string]interface{}{
					"PROFILE_ID":  lib.FormValue(r, "PROFILE_ID"),
					"AVATAR":      lib.FormValue(r, "AVATAR"),
					"ProfileType": profileType,
				}
			}

			res, err := lib.ApiRequestImpl(nil, r, fmt.Sprintf("%s?%s", constants.ProxyCreate, r.URL.RawQuery), false, j)
			error := app.GetError(err, res)
			if error != nil {
				log.Println(err, res.Body)
				ctx["error"] = error
			} else {
				lib.BuildTemplate(w, "ping/createAccountSuccess.html", ctx)
				return
			}
		}
		ctx["isGET"] = "false"
	} else {
		ctx["isGET"] = "true"
		app.DeleteCookie(w, constants.TOKEN_COOKIE_NAME)
	}
	lib.BuildTemplate(w, "ping/createAccount.html", ctx)
}

/**
 * LinkedIn post form request handler: /public/oauth
 */
func CreateAccountOauth(w http.ResponseWriter, r *http.Request, app *lib.App) {
	code := r.URL.Query().Get("code")
	redirect_uri := r.URL.Query().Get("state")
	if code == "" { //linkedin oauth
		app.Redirect(w, r, constants.LI_authorization_base_url+
			"?response_type=code&client_id="+constants.LI_client_id+
			"&scope="+constants.LI_scope_permissions+
			"&state="+redirect_uri+
			"&redirect_uri="+redirect_uri+"/public/oauth")
	} else {
		j := getJson(constants.LI_token_url+
			"?grant_type=authorization_code"+
			"&client_id="+constants.LI_client_id+
			"&client_secret="+constants.LI_client_secret+
			"&code="+code+
			"&redirect_uri="+redirect_uri+"/public/oauth", "")
		if j != nil {
			access_token := j["access_token"]
			me := getJson("https://api.linkedin.com/v2/me?projection=(id,localizedFirstName,localizedLastName,profilePicture(displayImage~:playableStreams))", access_token)
			if me != nil {
				profile := map[string]interface{}{
					"id":        me["id"],
					"firstName": me["localizedFirstName"],
					"lastName":  me["localizedLastName"],
				}
				if me["profilePicture"] != nil {
					displayImage := me["profilePicture"].(map[string]interface{})["displayImage~"]
					if displayImage != nil {
						elements := displayImage.(map[string]interface{})["elements"]
						if elements != nil {
							elements, ok := elements.([]interface{})
							if ok && len(elements) > 0 {
								identifiers := elements[0].(map[string]interface{})["identifiers"]
								if identifiers != nil {
									profile["pictureUrl"] = identifiers.([]interface{})[0].(map[string]interface{})["identifier"]
								}
							}
						}
					}
				}

				email := getJson("https://api.linkedin.com/v2/emailAddress?q=members&projection=(elements*(handle~))", access_token)
				elements, ok := email["elements"].([]interface{})
				if ok && len(elements) > 0 {
					handle := elements[0].(map[string]interface{})["handle~"]
					if handle != nil {
						profile["emailAddress"] = handle.(map[string]interface{})["emailAddress"]
					}
				}
				jstr, err := json.Marshal(profile)
				if err == nil {
					w.Write([]byte("<!DOCTYPE html>\n<html><title>LinkedIn OAuth</title>\n<script>\ntry{localStorage.setItem('profile', JSON.stringify(" + string(jstr) + "));}catch(e){alert(e)}\nwindow.close();\n</script>\n</html>"))
					return
				}
			}
		}
		w.Write([]byte("<!DOCTYPE html><script>window.close()</script></html>"))
	}
}

/**
 * AJAX handler for request: /public/verify_email
 */
func CreateAccountVerifyEmail(w http.ResponseWriter, r *http.Request, app *lib.App) {
	lib.ApiRequest(w, r, fmt.Sprintf("%s?%s", constants.ProxyVerifyEmail, r.URL.RawQuery))
}

/**
 * Selenium automation handler for request: /public/deleteAccount
 */
func DeleteAccount(w http.ResponseWriter, r *http.Request, app *lib.App) {
	lib.ApiRequest(w, r, fmt.Sprintf("%s?%s", constants.API_DeleteAccount, r.URL.RawQuery))
}

func getJson(url string, access_token interface{}) map[string]interface{} {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		utils.LogError(err)
		return nil
	}
	defer res.Body.Close()
	var j interface{}
	err = json.NewDecoder(res.Body).Decode(&j)
	if err != nil {
		utils.LogError(err)
		return nil
	}
	return j.(map[string]interface{})
}

func userForm(r *http.Request, dict map[string]interface{}) *lib.Form {
	profileType := lib.FormValue(r, "ProfileType")
	if profileType == "" {
		profileType = "FORM"
	}

	form := lib.InitForm(r, dict)
	form.CharField("Email", "Email Address:", "email", `focusIn title="Please enter a valid email address." required`)
	form.CharField("UserName", "Username:", "search", `minlength="3" required`)
	form.CharField("FirstName", "First Name:", "text", `required`)
	form.CharField("LastName", "Last Name:", "text", `required`)
	form.CharField("Password", "Create a password:", "password", `minlength="8" required`, nil)
	form.CharField("Confirm", "Confirm your password:", "password", `minlength="8" required`, nil)
	form.CharField("PROFILE_ID", "", "hidden", ``)
	form.CharField("AVATAR", "", "hidden", ``)
	form.CharField("ProfileType", "", "hidden", ``, profileType)
	return form
}

func ValidatePassword(password string, confirm_password string) error {
	if password != confirm_password {
		return errors.New("The passwords do not match. Please re-enter.")
	}

	invalidCheck := false
	numCheck := false
	lowerCheck := false
	upperCheck := false
	charCheck := true
	for _, ch := range password {
		i := int(ch)
		if i >= 48 && i <= 57 {
			numCheck = true
		} else if i >= 65 && i <= 90 {
			upperCheck = true
		} else if i >= 97 && i <= 122 {
			lowerCheck = true
		} else if strings.ContainsRune(constants.SPECIAL_CHARS, ch) {
			invalidCheck = true
		} else {
			charCheck = false
			break
		}
	}
	if !charCheck || !invalidCheck || !numCheck || !lowerCheck || !upperCheck {
		return errors.New("The password entered does not meet the criteria listed below. Please try entering another password.")
	}
	return nil
}
