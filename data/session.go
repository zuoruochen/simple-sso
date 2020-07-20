package data

import (
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
	"github.com/satori/go.uuid"
	"log"
	model2 "simple-sso/model"
	"simple-sso/util"
	"time"
)

func getSession(sessionId string) *model2.Session {
	value, err := RedisCli.Get(sessionId).Result()
	if err != nil {
		log.Println(err)
	}
	session:= new(model2.Session)
	if err = json.Unmarshal([]byte(value), session); err != nil {
		log.Println(err)
		return nil
	}
	return session
}

func getServiceTicket(ticket string) *model2.ServiceTicket {
	log.Println(ticket)
	value, err := RedisCli.Get(ticket).Result()
	if err != nil {
		log.Println(err)
	}
	st := new(model2.ServiceTicket)
	if err = json.Unmarshal([]byte(value), st); err != nil {
		log.Println(err)
		return nil
	}
	return st
}

func ValidateSession(sessionId, serviceDomain string) (*model2.Session, error) {
	session := getSession(sessionId)
	if session == nil {
		return nil, errors.New("no tgt session")
	}
	if validateDomain(session.Username, serviceDomain) {
		return session, nil
	}
	return nil, errors.New("service is illegal")
}

func validateDomain(username, serviceDomain string) bool {
	serviceDomainList := getUserServiceDomainList(username)
	if !stringInSlice(serviceDomain, serviceDomainList) {
		return false
	}
	return true
}

func ValidateTicket(serviceDomain, ticket string) (*model2.Session, error) {
	serviceTicket := getServiceTicket(ticket)
	if serviceTicket != nil && serviceTicket.ServiceDomain == serviceDomain {
		// when app go to cas to validate ticket, the cas session may be expired
		if session, err := ValidateSession(serviceTicket.SessionId, serviceDomain); err != nil {
			return nil, errors.New("user session expired, need to login again")
		} else {
			return session, nil
		}
	} else {
		return nil, errors.New("ticket not match service")
	}
}

// create session for user to signal sign on and ticket for service to validate
func GrantUser(username, serviceDomain string) (*model2.Session, *model2.ServiceTicket, error){
	session := &model2.Session{
		Username: username,
		Id: util.TGT_PREFIX + uuid.NewV4().String(),
		Ticket: make([]string, 0),
	}
	if serviceTicket, err := GenerateServiceTicket(RedisCli, session.Id, serviceDomain); err != nil {
		return nil, nil, err
	} else {
		session.Ticket = append(session.Ticket, serviceTicket.Ticket)
		if _, err := saveSession(RedisCli, session.Id, session); err != nil {
			return nil, nil, err
		}
		return session, serviceTicket, nil
	}
}

func AuthUser(username, password, serviceDomain string) (bool, error) {
	if _, err := authUser(username, password); err != nil {
		log.Print(err)
		return false, err
	}
	if status := validateDomain(username, serviceDomain); !status {
		return false, errors.New("unregistered service")
	}
	return true, nil
}

func GenerateServiceTicket(redisCli *redis.Client, sessionId, serviceDomain string) (*model2.ServiceTicket, error) {
	serviceTicket := &model2.ServiceTicket{
		ServiceDomain: serviceDomain,
		Ticket: util.ST_PREFIX + uuid.NewV4().String(),
		SessionId: sessionId,
	}
	if _, err := saveServiceTicket(redisCli, serviceTicket); err != nil {
		return nil, err
	}
	return serviceTicket, nil
}

func AddTicketToSession(redisCli *redis.Client, sessionId string, ticket *model2.ServiceTicket) error{
	val, err := redisCli.Get(sessionId).Result()
	if err == redis.Nil {
		log.Printf("session is not exist, sessionId: %s", sessionId)
		return errors.New("no session exist")
	} else if err != nil {
		log.Printf("get session from redis error, key: %s", sessionId)
		return err
	}
	session := new(model2.Session)
	err = json.Unmarshal([]byte(val), session)
	if err != nil {
		log.Print(err)
		return err
	}
	session.Ticket = append(session.Ticket, ticket.Ticket)
	if _, err = setSessionExpireAt(redisCli, session.ExpiredAt, sessionId, session); err != nil {
		return err
	}
	return nil
}

func saveServiceTicket(redisCli *redis.Client, ticket *model2.ServiceTicket) (bool, error){
	jsonTicket, err := json.Marshal(ticket)
	if err != nil {
		log.Print(err)
		return false, err
	}
	if err := redisCli.Set(ticket.Ticket, string(jsonTicket), time.Millisecond*util.SERVICE_TICKET_TIME_TO_LIVE).Err(); err != nil {
		return false, err
	}
	return true, nil
}

func saveSession(redisCli *redis.Client, key string, session *model2.Session) (bool, error) {
	log.Println("time now: %s, expireAt : %s\n", time.Now(), time.Now().Add(time.Millisecond*util.SESSION_TIME_TO_LIVE))
	return setSessionExpireAt(redisCli, time.Now().Add(time.Millisecond*util.SESSION_TIME_TO_LIVE), key, session)
}

func setSessionExpireAt(redisCli *redis.Client, expireAt time.Time, key string, session *model2.Session) (bool, error) {
	session.ExpiredAt = expireAt
	jsonSession, err := json.Marshal(session)
	if err != nil {
		log.Print(err)
		return false, err
	}
	return setValueExpireAt(redisCli, session.ExpiredAt, key, string(jsonSession))
}

func setValueExpireAt(redisCli *redis.Client, expireAt time.Time, key, value string) (bool, error){
	log.Printf("key : %s, value: %s, expireAt : %s\n", key, value, expireAt)
	if err := redisCli.Set(key, value, 0).Err(); err != nil {
		log.Print(err)
		return false, err
	}
	if err := redisCli.ExpireAt(key, expireAt).Err(); err != nil {
		log.Print(err)
		return false, err
	}
	return true, nil
}

func authUser(username, password string) (bool, error){
	// todo: load user data from db and compare password
	return true, nil
}

func getUserServiceDomainList(username string) []string {
	// todo: load user related serviceDomain list from db
	return []string{"localhost:8081", "localhost:8082"}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}