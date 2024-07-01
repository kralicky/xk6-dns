package dns

import (
	"context"
	"fmt"
	"net"
	"time"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/dns", &DNS{})
}

type DNS struct{}

type CommonFields struct {
	Err      string        `js:"err"`
	Duration time.Duration `js:"duration"`
}

func (c *CommonFields) SetCommonFields(duration time.Duration, err error) {
	c.Duration = duration
	if err != nil {
		c.Err = err.Error()
	}
}

type LookupAddrResponse struct {
	CommonFields
	Names []string `js:"names"`
}

type LookupCNAMEResponse struct {
	CommonFields
	Cname string `js:"cname"`
}

type LookupHostResponse struct {
	CommonFields
	Addrs []string `js:"addrs"`
}

type LookupIPResponse struct {
	CommonFields
	Ips []string `js:"ips"`
}

type LookupIPAddrResponse struct {
	CommonFields
	Ips []string `js:"ips"`
}

type LookupMXResponse struct {
	CommonFields
	Records []*net.MX `js:"records"`
}

type LookupNSResponse struct {
	CommonFields
	Records []*net.NS `js:"records"`
}

type LookupNetIPResponse struct {
	CommonFields
	Ips []string `js:"ips"`
}

type LookupPortResponse struct {
	CommonFields
	Port int `js:"port"`
}

type LookupSRVResponse struct {
	CommonFields
	Cname   string     `js:"cname"`
	Records []*net.SRV `js:"records"`
}

type LookupTXTResponse struct {
	CommonFields
	Records []string `js:"records"`
}

func measure[T interface {
	SetCommonFields(duration time.Duration, err error)
}](fn func() (T, error)) T {
	now := time.Now()
	t, err := fn()
	t.SetCommonFields(time.Since(now), err)
	return t
}

func toStrings[T fmt.Stringer, S ~[]T](s S) []string {
	strs := make([]string, len(s))
	for i, v := range s {
		strs[i] = v.String()
	}
	return strs
}

func toPtrs[T any, S ~[]T](s S) []*T {
	ptrs := make([]*T, len(s))
	for i, v := range s {
		ptrs[i] = &v
	}
	return ptrs
}

func (*DNS) LookupHost(host string) *LookupHostResponse {
	return measure(func() (*LookupHostResponse, error) {
		addrs, err := net.DefaultResolver.LookupHost(context.Background(), host)
		return &LookupHostResponse{Addrs: addrs}, err
	})
}

func (r *DNS) LookupAddr(addr string) *LookupAddrResponse {
	return measure(func() (*LookupAddrResponse, error) {
		names, err := net.DefaultResolver.LookupAddr(context.Background(), addr)
		return &LookupAddrResponse{Names: names}, err
	})
}

func (r *DNS) LookupCNAME(host string) *LookupCNAMEResponse {
	return measure(func() (*LookupCNAMEResponse, error) {
		cname, err := net.DefaultResolver.LookupCNAME(context.Background(), host)
		return &LookupCNAMEResponse{Cname: cname}, err
	})
}

func (r *DNS) LookupIP(network string, host string) *LookupIPResponse {
	return measure(func() (*LookupIPResponse, error) {
		ips, err := net.DefaultResolver.LookupIP(context.Background(), network, host)
		return &LookupIPResponse{Ips: toStrings(ips)}, err
	})
}

func (r *DNS) LookupIPAddr(host string) *LookupIPAddrResponse {
	return measure(func() (*LookupIPAddrResponse, error) {
		ips, err := net.DefaultResolver.LookupIPAddr(context.Background(), host)
		return &LookupIPAddrResponse{Ips: toStrings(toPtrs(ips))}, err
	})
}

func (r *DNS) LookupMX(name string) *LookupMXResponse {
	return measure(func() (*LookupMXResponse, error) {
		records, err := net.DefaultResolver.LookupMX(context.Background(), name)
		return &LookupMXResponse{Records: records}, err
	})
}

func (r *DNS) LookupNS(name string) *LookupNSResponse {
	return measure(func() (*LookupNSResponse, error) {
		records, err := net.DefaultResolver.LookupNS(context.Background(), name)
		return &LookupNSResponse{Records: records}, err
	})
}

func (r *DNS) LookupNetIP(network string, host string) *LookupNetIPResponse {
	return measure(func() (*LookupNetIPResponse, error) {
		ips, err := net.DefaultResolver.LookupNetIP(context.Background(), network, host)
		return &LookupNetIPResponse{Ips: toStrings(ips)}, err
	})
}

func (r *DNS) LookupPort(network string, service string) *LookupPortResponse {
	return measure(func() (*LookupPortResponse, error) {
		port, err := net.DefaultResolver.LookupPort(context.Background(), network, service)
		return &LookupPortResponse{Port: port}, err
	})
}

func (r *DNS) LookupSRV(service string, proto string, name string) *LookupSRVResponse {
	return measure(func() (*LookupSRVResponse, error) {
		cname, records, err := net.DefaultResolver.LookupSRV(context.Background(), service, proto, name)
		return &LookupSRVResponse{Cname: cname, Records: records}, err
	})
}

func (r *DNS) LookupTXT(name string) *LookupTXTResponse {
	return measure(func() (*LookupTXTResponse, error) {
		res, err := net.DefaultResolver.LookupTXT(context.Background(), name)
		return &LookupTXTResponse{Records: res}, err
	})
}
