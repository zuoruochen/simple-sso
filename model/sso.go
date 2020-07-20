package model

import "time"

// Session id has prefix util.TGT_PREFIX
// Ticket has prefix util.ST_PREFIX
type Session struct {
	Username string `json:"username""`
	Id string `json:"id"`
	Ticket []string `json:"ticket"`
	ExpiredAt time.Time
}


type ServiceTicket struct {
	ServiceDomain string
	SessionId string
	Ticket string
}

type UserServiceList struct {
	Username string
	Password string
	Service []string
}

type Login struct {
	Username    string `form:"username" json:"user" uri:"user" xml:"user" binding:"required"`
	Password string `form:"password" json:"password" uri:"password" xml:"password" binding:"required"`
}