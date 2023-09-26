package prober

import (
	"dtool/utils"
)

type SVCBResult struct {
	Ip       string             `json:"ip"`
	Response utils.SVCBResponse `json:"response"`
}

func SVCBProbeOnce(ip string, domain string) (SVCBResult, error) {
	result := SVCBResult{Ip: ip}
	res, err := utils.SendSVCBQuery(ip, domain)
	if err != nil {
		return result, err
	}
	resp, err := utils.ParseSVCBResponse(res)
	if err != nil {
		return result, err
	}
	result.Response = resp
	return result, nil
}

func SVCBProbe(ip string, domain string) SVCBResult {
	result := SVCBResult{Ip: ip}
	res, err := utils.SendSVCBQuery(ip, domain)
	if err != nil {
		return result
	}
	resp, err := utils.ParseSVCBResponse(res)
	if err != nil {
		return result
	}
	result.Response = resp
	return result
}
