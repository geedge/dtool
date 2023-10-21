// dns_utils.go contains the necessary functions for building
// and parsing a DNS packet
package utils

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

type WrongAnswerError struct {
	Message string
}

type QueryStruct struct {
	Id     uint16
	RD     bool
	Qname  string
	Qclass uint16
	Qtype  uint16
}

type DomainInfo struct {
	IPList []string
	NSList map[string][]string
}

type SVCBRecord struct {
	Target string            `json:"target"`
	Data   map[string]string `json:"value"`
}

type SVCBResponse struct {
	Records []SVCBRecord `json:"records"`
}

func (e *WrongAnswerError) Error() string {
	return fmt.Sprintf("Wrong Answer: %s", e.Message)
}

// build the question section of a dns packet
func QuestionMaker(domain string, qclass uint16, qtype uint16) *dns.Question {
	return &dns.Question{Name: domain, Qtype: qtype, Qclass: qclass}
}

// build a dns query message
func QueryMaker(query QueryStruct) *dns.Msg {
	msg := new(dns.Msg)
	msg.Id = query.Id
	msg.RecursionDesired = query.RD

	var query_name string
	var query_class, query_type uint16
	msg.Question = make([]dns.Question, 1)
	query_name = dns.Fqdn(query.Qname)

	// default class INET
	if query.Qclass == 0 {
		query_class = 1
	} else {
		query_class = query.Qclass
	}

	// default type A
	if query.Qtype == 0 {
		query_type = 1
	} else {
		query_type = query.Qtype
	}

	question := QuestionMaker(query_name, query_class, query_type)
	msg.Question[0] = *question
	return msg
}

func ParseRcode(msg *dns.Msg) int {
	return msg.MsgHdr.Rcode
}

func ParseAResponse(msg *dns.Msg) (string, error) {
	var ip_str string
	if len(msg.Answer) == 1 {
		if a, ok := msg.Answer[0].(*dns.A); ok {
			ip_str = a.A.String()
		} else {
			err := &WrongAnswerError{
				Message: "Wrong record type",
			}
			return ip_str, err
		}
	} else {
		err := &WrongAnswerError{
			Message: "Wrong answer section",
		}
		return ip_str, err
	}
	return ip_str, nil
}

func ParseCNAMEChain(msg *dns.Msg) (int, string, error) {
	var cache_id int
	var rdns_ip string
	if len(msg.Answer) == 3 {
		if cname, ok := msg.Answer[0].(*dns.CNAME); ok {
			rdns_ip = strings.Join(strings.Split(strings.Split(cname.Target, ".")[0], "-")[1:], ".")
			if a, ok := msg.Answer[2].(*dns.A); ok {
				cache_id = int(binary.BigEndian.Uint32(a.A))
			}
		}
	} else {
		err := &WrongAnswerError{
			Message: "Wrong record number",
		}
		return cache_id, rdns_ip, err
	}
	return cache_id, rdns_ip, nil
}

func ParseTXTResponse(msg *dns.Msg) (string, error) {
	var txt_string string
	if len(msg.Answer) == 1 {
		if txt, ok := msg.Answer[0].(*dns.TXT); ok {
			txt_string = txt.String()
		} else {
			err := &WrongAnswerError{
				Message: "Wrong record type",
			}
			return txt_string, err
		}
	} else {
		err := &WrongAnswerError{
			Message: "wrong record number",
		}
		return txt_string, err
	}
	return txt_string, nil
}

func SendQuery(ip, domain, port string, qtype, qclass uint16) (*dns.Msg, error) {
	addr := ip + ":" + port
	query := new(QueryStruct)
	query.Qname = domain
	query.RD = true
	query.Qtype = qtype
	query.Qclass = uint16(qclass)
	query.Id = dns.Id()
	m := QueryMaker(*query)

	res, err := dns.Exchange(m, addr)

	return res, err
}

func SendAQuery(ip string, domain string) (*dns.Msg, error) {
	addr := ip + ":53"
	query := new(QueryStruct)
	query.Qname = domain
	query.RD = true
	query.Id = dns.Id()
	m := QueryMaker(*query)

	res, err := dns.Exchange(m, addr)

	return res, err
}

func SendVersionQuery(ip string) (*dns.Msg, error) {
	addr := ip + "53"
	query := new(QueryStruct)
	query.Id = dns.Id()
	query.Qname = "version.bind"
	query.Qclass = dns.ClassCHAOS
	query.Qtype = dns.TypeTXT
	m := QueryMaker(*query)

	res, err := dns.Exchange(m, addr)

	return res, err
}

func SendSVCBQuery(ip string, domain string) (*dns.Msg, error) {
	addr := ip + ":53"
	query := QueryStruct{Id: dns.Id(), Qname: dns.Fqdn(domain), RD: true, Qtype: dns.TypeSVCB}
	m := QueryMaker(query)

	res, err := dns.Exchange(m, addr)

	return res, err
}

func toSVCBKEY(key dns.SVCBKey) string {
	switch key {
	case dns.SVCB_MANDATORY:
		return "mandatory"
	case dns.SVCB_ALPN:
		return "alpn"
	case dns.SVCB_NO_DEFAULT_ALPN:
		return "no_default_alpn"
	case dns.SVCB_PORT:
		return "port"
	case dns.SVCB_IPV4HINT:
		return "ipv4_hint"
	case dns.SVCB_ECHCONFIG:
		return "ech_config"
	case dns.SVCB_IPV6HINT:
		return "ipv6_hint"
	case dns.SVCB_DOHPATH:
		return "doh_path"
	default:
		return "unknown"
	}
}

func ParseSVCBResponse(msg *dns.Msg) (SVCBResponse, error) {
	response := SVCBResponse{Records: make([]SVCBRecord, 0)}
	if len(msg.Answer) > 0 {
		for _, rr := range msg.Answer {
			if svcb, ok := rr.(*dns.SVCB); ok {
				record := SVCBRecord{Data: make(map[string]string)}
				record.Target = svcb.Target
				for _, kv := range svcb.Value {
					key := toSVCBKEY(kv.Key())
					value := kv.String()
					record.Data[key] = value
				}
				response.Records = append(response.Records, record)
			}
		}
	}
	if len(response.Records) > 0 {
		return response, nil
	}
	return response, &WrongAnswerError{Message: "no valid SVCB records"}
}
