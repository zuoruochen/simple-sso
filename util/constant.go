package util

const TGT_PREFIX = "TGT-"
const ST_PREFIX = "ST-"

const SERVICE_TICKET_TIME_TO_LIVE = 12 * 60 * 60 * 1000
const SESSION_TIME_TO_LIVE = 24 * 60 * 60 * 1000

const CAS_HOST = "localhost:8088"

const HELLO_CLIENT_HOST= "localhost:8081"

// cause browser set the cookie to the hostname, this test client2 will not go to cas but auth succeed directly
const HELLO_CLIENT2_HOST= "localhost:8082"

const HTTP = "http"

const HTTPS = "https"