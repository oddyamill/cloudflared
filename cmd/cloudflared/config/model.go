package config

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"
)

// Forwarder represents a client side listener to forward traffic to the edge
type Forwarder struct {
	URL      string `json:"url"`
	Listener string `json:"listener"`
}

// Tunnel represents a tunnel that should be started
type Tunnel struct {
	URL          string `json:"url"`
	Origin       string `json:"origin"`
	ProtocolType string `json:"type"`
}

// DNSResolver represents a client side DNS resolver
type DNSResolver struct {
	Enabled    bool     `json:"enabled"`
	Address    string   `json:"address"`
	Port       uint16   `json:"port"`
	Upstreams  []string `json:"upstreams"`
	Bootstraps []string `json:"bootstraps"`
}

// Root is the base options to configure the service
type Root struct {
	OrgKey          string      `json:"org_key"`
	ConfigType      string      `json:"type"`
	CheckinInterval int         `json:"checkin_interval"`
	Forwarders      []Forwarder `json:"forwarders,omitempty"`
	Tunnels         []Tunnel    `json:"tunnels,omitempty"`
	Resolver        DNSResolver `json:"resolver"`
}

// Hash returns the computed values to see if the forwarder values change
func (f *Forwarder) Hash() string {
	h := md5.New()
	io.WriteString(h, f.URL)
	io.WriteString(h, f.Listener)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Hash returns the computed values to see if the forwarder values change
func (r *DNSResolver) Hash() string {
	h := md5.New()
	io.WriteString(h, r.Address)
	io.WriteString(h, strings.Join(r.Bootstraps, ","))
	io.WriteString(h, strings.Join(r.Upstreams, ","))
	io.WriteString(h, fmt.Sprintf("%d", r.Port))
	io.WriteString(h, fmt.Sprintf("%v", r.Enabled))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// EnabledOrDefault returns the enabled property
func (r *DNSResolver) EnabledOrDefault() bool {
	return r.Enabled
}

// AddressOrDefault returns the address or returns the default if empty
func (r *DNSResolver) AddressOrDefault() string {
	if r.Address != "" {
		return r.Address
	}
	return "localhost"
}

// PortOrDefault return the port or returns the default if 0
func (r *DNSResolver) PortOrDefault() uint16 {
	if r.Port > 0 {
		return r.Port
	}
	return 53
}

// UpstreamsOrDefault returns the upstreams or returns the default if empty
func (r *DNSResolver) UpstreamsOrDefault() []string {
	if len(r.Upstreams) > 0 {
		return r.Upstreams
	}
	return []string{"https://1.1.1.1/dns-query", "https://1.0.0.1/dns-query"}
}

// BootstrapsOrDefault returns the bootstraps or returns the default if empty
func (r *DNSResolver) BootstrapsOrDefault() []string {
	if len(r.Bootstraps) > 0 {
		return r.Bootstraps
	}
	return []string{"https://162.159.36.1/dns-query", "https://162.159.46.1/dns-query", "https://[2606:4700:4700::1111]/dns-query", "https://[2606:4700:4700::1001]/dns-query"}
}