package apiweb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/solution9th/NSBridge/oidc"
	"github.com/solution9th/NSBridge/service/cache"
	"github.com/solution9th/NSBridge/utils"
)

const (
	AppUrl       = "http://mf.sddeznsm.com"
	Domain       = "mf.sddeznsm.com"
	CallBackUrl  = "http://mf.sddeznsm.com/callback"
	AuthorizeUrl = "https://org.i1.dev.com/sso/oidc/authorize"
	I1TokenUrl   = "https://org.i1.dev.com/sso/oidc/token"

	DnsIdCookieKey = "dns_id"
)

//uid:       "hello"
//firstname: "三"
//lastname:  "张"
//name:      "张三"
//email:     "hello@mail.com"
//avatar:    "http://avatar.com/abc.png"
//token:     "ACCESS_TOKEN"
//mobile:    13800138000
//telephone: 010-12345678
type User struct {
	UserID    string `json:"uid"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Fullname  string `json:"fullname"`
	Avatar    string `json:"avatar"`
	Token     string `json:"token"`
	Email     string `json:"email"`
	Telephone string `json:"telephone"`
	Mobile    string `json:"mobile"`
}

type AuthToken struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
	Scope        string `json:"scope"`
	ExpiresIn    int64  `json:"expires_in"`
}

// Ping test router
func Ping(c *gin.Context) {

	c.JSON(http.StatusOK, utils.ParseSuccessWithData(map[string]interface{}{
		"msg":       "ok",
		"create_at": time.Now().Unix(),
	}))
	return
}

func OIDCLogin(c *gin.Context) {
	// parse jwt info
	id_token := c.Query("id_token")
	if id_token != "" {
		log.Println(id_token)
		payload, err := oidc.ParseJWT(id_token)
		var user User
		if err != nil {
			c.AbortWithError(400, err)
		}
		if err := json.Unmarshal(payload, &user); err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			c.AbortWithError(400, err)
			return
		}
	}

	OIDCAuth(c)
}

func OIDCAuth(c *gin.Context) {

	params := &struct {
		Location     string
		clientId     string
		redirectUri  string
		state        string
		responseType string
	}{
		// TODO
		Location:     AuthorizeUrl,
		clientId:     "",
		redirectUri:  fmt.Sprintf("%s/callback", AppUrl),
		state:        "",
		responseType: "code",
	}
	c.Redirect(302, fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&state=%s&response_type=%s",
		params.Location, params.clientId, params.redirectUri, params.state, params.responseType))
}

func CallBack(c *gin.Context) {

	code := c.Query("code")
	state := c.Query("state")

	client := http.Client{
		Timeout: 6 * time.Second,
	}
	// todo
	url := fmt.Sprintf("%s?grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&client_secret=%s",
		I1TokenUrl, code, CallBackUrl, state, "6rugz6ocxrou96lyc205wk37")
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	result, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(result))
	if err != nil {
		panic(err)
	}
	var authToken AuthToken

	err = json.Unmarshal(result, &authToken)
	if err != nil {
		panic(err)
	}
	payload, err := oidc.ParseJWT(authToken.IdToken)

	var user User
	if err != nil {
		c.AbortWithError(404, err)
	}
	if err := json.Unmarshal(payload, &user); err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		c.AbortWithError(404, err)
		return
	}
	session := sessions.Default(c)
	session.Set(UserSessionKey, user.UserID)
	session.Save()
	c.SetCookie(DnsIdCookieKey, state, 3600, "/", Domain, false, true)

	token := utils.GenToken(16)

	c.SetCookie(TokenCookieKey, token, 3600, "/", Domain, false, true)

	err = cache.DefaultCache.Set(fmt.Sprintf(passportTokenFormat, state, user.UserID, token), time.Now().Format(time.RFC3339), 1*time.Hour)
	if err != nil {

	}

	c.Redirect(302, "/web")

}
