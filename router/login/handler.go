package login

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	url2 "net/url"
	"simple-sso/data"
	"simple-sso/model"
	"simple-sso/util"
	"time"
)

func getLoginHtmlHandler(c *gin.Context) {
	serviceUrl := c.Query("service")
	url, err := url2.Parse(serviceUrl)
	if err != nil {
		c.JSON(http.StatusBadRequest, data.Response{http.StatusBadRequest,err.Error(), nil})
		return
	}
	if sessionId, err:= c.Cookie("CASTGC"); err != nil {
		if err == http.ErrNoCookie {
			log.Print(err)
			c.HTML(http.StatusOK, "login.html", gin.H{"loginUrl": "/login?service=" + serviceUrl})
			return
		} else {
			c.JSON(http.StatusInternalServerError, data.Response{http.StatusInternalServerError, "Internal server error", nil})
			return
		}
	} else {
		if _, err := data.ValidateSession(sessionId, url.Host); err != nil {
			c.JSON(http.StatusUnauthorized, data.Response{http.StatusUnauthorized, err.Error(), nil})
			return
		} else {
			newTicket, err := data.GenerateServiceTicket(data.RedisCli, sessionId, url.Host)
			if err != nil {
				c.JSON(http.StatusInternalServerError, data.Response{http.StatusInternalServerError, "Internal server error", nil})
				return
			}
			if err = data.AddTicketToSession(data.RedisCli, sessionId, newTicket); err != nil {
				c.JSON(http.StatusInternalServerError, data.Response{http.StatusInternalServerError, "Internal server error", nil})
				return
			}
			c.Redirect(http.StatusFound, util.GetHttpUrl(url.Scheme, url.Host)+"?ticket="+newTicket.Ticket)
			return
		}
	}

}

func loginUserHandler(c *gin.Context) {
	serviceUrl := c.Query("service")
	log.Printf("serviceUrl: %s", serviceUrl)
	url, err := url2.Parse(serviceUrl)
	if err != nil {
		c.JSON(http.StatusBadRequest, data.Response{http.StatusBadRequest,err.Error(), nil})
		return
	}
	var form model.Login
	if err := c.Bind(&form); err != nil {
		c.JSON(http.StatusBadRequest, data.Response{http.StatusBadRequest,err.Error(), nil})
		return
	}
	if _, err := data.AuthUser(form.Username, form.Password, url.Host); err != nil {
		c.JSON(http.StatusUnauthorized, data.Response{http.StatusUnauthorized, err.Error(), nil})
		return
	}
	// grant user, this function will create tgt session and service ticket
	if session, ticket, err := data.GrantUser(form.Username, url.Host); err != nil {
		c.JSON(http.StatusInternalServerError, data.Response{http.StatusInternalServerError, "Internal server error", nil})
		return
	} else {
		log.Println(session)
		log.Println(ticket)
		c.SetCookie("CASTGC", session.Id, util.SESSION_TIME_TO_LIVE, "/", util.CAS_HOST, false, false)
		log.Printf("sessionId: %s, service %s ticket :%s", session.Id, serviceUrl, ticket.Ticket)
		c.Redirect(http.StatusFound, util.GetHttpUrl(url.Scheme, url.Host)+"?ticket="+ticket.Ticket)
		return
	}

}

func serviceValidateHandler(c *gin.Context) {
	service := c.Query("service")
	ticket := c.Query("ticket")
	log.Printf("service: %s, ticket: %s", service, ticket)
	if len(service) == 0 || len(ticket) == 0 {
		c.JSON(http.StatusBadRequest, data.Response{http.StatusBadRequest,"service or ticket is null", nil})
		return
	}
	url, err := url2.Parse(service)
	if err != nil {
		c.JSON(http.StatusBadRequest, data.Response{http.StatusBadRequest,err.Error(), nil})
		return
	}
	if session, err := data.ValidateTicket(url.Host, ticket); err != nil {
		c.JSON(http.StatusForbidden, data.Response{http.StatusForbidden,err.Error(), nil})
		return
	} else {
		res := make(map[string]interface{})
		res["user"] = session.Username
		res["expireAt"] = session.ExpiredAt.Format(time.RFC3339)
		res["ticketValidatedAt"] = time.Now().Format(time.RFC3339)
		res["redirectUrl"] = service
		c.JSON(http.StatusOK,data.Response{http.StatusOK,"ticket validate success", res})
		return
	}
}