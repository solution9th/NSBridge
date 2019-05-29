package apiweb

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/solution9th/NSBridge/config"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/solution9th/go-saml"

	"github.com/solution9th/NSBridge/oidc"
	"github.com/solution9th/NSBridge/service/cache"
	"github.com/solution9th/NSBridge/utils"
)

const (
	UserSessionKey = "user_id:user_name"
	TokenCookieKey   = "token_id"
)

var sp saml.ServiceProviderSettings

// SAMLLogin saml 登录处理
func SAMLLogin(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(UserSessionKey)
	if user == nil {
		SAMLAuth(c)
		return
	}

	token, err := c.Cookie(TokenCookieKey)
	if err != nil {
		SAMLAuth(c)
		return
	}
	exist, err := cache.DefaultCache.Exist(fmt.Sprintf(passportTokenFormat, user, token))
	if err != nil {
		utils.Error("Exist Session Err: ", err.Error())
		c.AbortWithStatusJSON(500, gin.H{
			"msg": errors.New("server err"),
		})
		return
	}
	if !exist {
		SAMLAuth(c)
		return
	}
	// parse jwt info
	idToken := c.Query("id_token")
	if idToken != "" {
		payload, err := oidc.ParseJWT(idToken)
		var user User
		if err != nil {
			utils.Error("ParseJWT Err: ", err.Error())
			c.AbortWithStatusJSON(500, gin.H{
				"msg": errors.New("server err"),
			})
		}
		if err := json.Unmarshal(payload, &user); err != nil {
			utils.Error("json.Unmarshal Err: ", err.Error())
			c.AbortWithStatusJSON(500, gin.H{
				"msg": errors.New("server err"),
			})
			return
		}
	}

	SAMLAuth(c)
}

func SAMLAuth(c *gin.Context) {
	RelayState := c.Query("RelayState")
	if RelayState == "" {
		RelayState = "/"
	}
	sp = saml.ServiceProviderSettings{
		IDPSSOURL:                   config.SamlConfig.IDPSSOURL,
		IDPSSODescriptorURL:         config.SamlConfig.IDPSSODescriptorURL,
		IDPPublicCertPath:           "/etc/ns_bridge/apiweb/certs/idp.crt",
		PublicCertPath:              "/etc/ns_bridge/apiweb/certs/idp.crt",
		PrivateKeyPath:              "/etc/ns_bridge/apiweb/certs/idp.key",
		SPSignRequest:               true,
		AssertionConsumerServiceURL: config.SamlConfig.AssertionConsumerServiceURL,
	}
	err := sp.Init()
	if err != nil {
		utils.Error("sp.Init Err: ", err.Error())
		c.AbortWithStatusJSON(500, gin.H{
			"msg": errors.New("sp init server err"),
		})
		return
	}

	// generate the AuthnRequest and then get a base64 encoded string of the XML
	authnRequest := sp.GetAuthnRequest()
	b64XML, err := authnRequest.EncodedSignedString(sp.PrivateKeyPath)
	if err != nil {
		utils.Error("SAMLAuth Err: ", err.Error())
		c.AbortWithStatusJSON(500, gin.H{
			"msg": errors.New("saml init server err"),
		})
		return
	}

	// for convenience, get a URL formed with the SAMLRequest parameter
	url, err := saml.GetAuthnRequestURL(sp.IDPSSOURL, b64XML, "")
	if err != nil {
		utils.Error("SAMLAuth Err:", err.Error())
		c.AbortWithStatusJSON(500, gin.H{
			"msg": errors.New("SAMLAuth server err"),
		})
		return
	}

	// below is bonus for how you might respond to a request with a form that POSTs to the IdP
	data := struct {
		Base64AuthRequest string
		URL               string
		RelayState        string
	}{
		Base64AuthRequest: b64XML,
		URL:               url,
		RelayState:        RelayState,
	}

	t := template.New("saml")
	t, err = t.Parse("<html><body style=\"display: none\" onload=\"document.frm.submit()\"><form method=\"post\" name=\"frm\" action=\"{{.URL}}\"><input type=\"hidden\" name=\"SAMLRequest\" value=\"{{.Base64AuthRequest}}\" /><input type=\"hidden\" name=\"RelayState\" value=\"{{.RelayState}}\" /><input type=\"submit\" value=\"Submit\" /></form></body></html>")

	// how you might respond to a request with the templated form that will auto post
	t.Execute(c.Writer, data)

}

func SAMLAcs(c *gin.Context) {

	encodedXML := c.Request.FormValue("SAMLResponse")
	RelayState := c.Request.FormValue("RelayState")

	if encodedXML == "" {
		utils.Error("SAMLResponse form value missing")
		c.AbortWithStatusJSON(500, gin.H{
			"msg": errors.New("SAMLResponse form value missing"),
		})
		return
	}

	response, err := saml.ParseEncodedResponse(encodedXML)
	if err != nil {
		utils.Error("SAMLResponse parse: " + err.Error())
		c.AbortWithStatusJSON(500, gin.H{
			"msg": errors.New("SAMLResponse parse err"),
		})
		return
	}

	err = response.Validate(&sp)
	if err != nil {
		utils.Error("SAMLResponse validation: " + err.Error())
		c.AbortWithStatusJSON(500, gin.H{
			"msg": errors.New("SAMLResponse validation err"),
		})
		return
	}

	userID := response.GetAttributeValue("username")
	if userID == "" {
		utils.Error("SAML attribute identifier uid missing", err.Error())
		c.AbortWithStatusJSON(500, gin.H{
			"msg": errors.New("SAML attribute identifier uid missing"),
		})
		return
	}

	userName := response.GetAttributeValue("name")
	if userName == "" {
		utils.Error("SAML attribute identifier name missing", err.Error())
		c.AbortWithStatusJSON(500, gin.H{
			"msg": errors.New("SAML attribute identifier name missing"),
		})
		return
	}
	session := sessions.Default(c)
	session.Set(UserSessionKey, userID + ":" + userName)
	session.Save()

	token := utils.GenToken(16)
	c.SetCookie(TokenCookieKey, token, 3600, "/", config.SamlConfig.Domain, false, true)

	err = cache.DefaultCache.Set(fmt.Sprintf(passportTokenFormat, userID + ":" + userName, token), time.Now().Format(time.RFC3339), 1*time.Hour)
	if err != nil {
		utils.Error("DefaultCache Set Err", err.Error())
		c.AbortWithStatusJSON(500, gin.H{
			"msg": errors.New("DefaultCache Set Err"),
		})
		return
	}

	c.Redirect(302, RelayState)
}

func GetI1Notice(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"receive" : "yes",
	})
}