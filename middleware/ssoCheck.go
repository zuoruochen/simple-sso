package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"simple-sso/data"
	"simple-sso/util"
	"time"
)

func SSOCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("check user status begin")
		defer log.Println("check user status end")
		originUrl := util.GetFullUrl(c.Request)
		ticket := c.Query("ticket")
		if len(ticket) == 0 {
			if sessionId, err := c.Cookie("JSESSIONID"); err != nil {
				if err == http.ErrNoCookie {
					log.Printf("no user login status, go to cas to login")
					c.Redirect(http.StatusFound, util.GetHttpUrl(util.HTTP, util.CAS_HOST) + "/login?service=" + originUrl)
					return
				} else {
					c.JSON(http.StatusInternalServerError, data.Response{http.StatusInternalServerError, "Internal server error", nil})
					return
				}
			} else {
				//  check user status in session, because of sample code, we skip this step
				if !checkUserStatusInSession(sessionId) {
					log.Println("check user failed, go to cas to login")
					c.Redirect(http.StatusFound, util.GetHttpUrl(util.HTTP, util.CAS_HOST) + "/login?service=" + originUrl)
					return
				} else {
					c.Next()
				}
			}
		} else {
			res := validateTicket(util.GetBaseUrl(c.Request), ticket)
			if res.Status == http.StatusOK {
				ret := res.Data.(map[string]interface{})
				expireAt, err := time.Parse(time.RFC3339, ret["expireAt"].(string))
				if err != nil {
					panic(err)
				}
				maxAge := int(expireAt.Sub(time.Now()).Seconds())
				if maxAge > 0 {
					log.Printf("origin url : %s", originUrl)
					c.SetCookie("JSESSIONID", "TEST", util.SESSION_TIME_TO_LIVE, "/", c.Request.Host, false, false)
					c.Redirect(http.StatusFound, util.GetBaseUrl(c.Request))
					return
				} else {
					log.Println("cas session expired, go to cas login again")
					c.Redirect(http.StatusFound, util.GetHttpUrl(util.HTTP, util.CAS_HOST) + "/login?service=" + util.GetBaseUrl(c.Request))
					return
				}
			} else {
				log.Println("validate ticket failed, go to cas login again")
				c.Redirect(http.StatusFound, util.GetHttpUrl(util.HTTP, util.CAS_HOST) + "/login?service=" + util.GetBaseUrl(c.Request))
				return
			}
		}
	}
}

func checkUserStatusInSession(sessionId string) bool {
	return true
}

func validateTicket(service, ticket string) *data.Response {
	url := fmt.Sprintf("%s://%s%s?service=%s&ticket=%s", util.HTTP, util.CAS_HOST, "/serviceValidate", service, ticket)
	log.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Println(body)
	if err != nil {
		panic(err)
	}
	ret := new(data.Response)
	if err = json.Unmarshal(body, ret); err != nil {
		panic(err)
	}
	return ret
}