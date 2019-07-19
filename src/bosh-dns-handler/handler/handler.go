package handler

import (
	"github.com/aemengo/bosh-dns-handler/config"
	"github.com/miekg/dns"
	"log"
	"net"
)

type Handler struct{
	cfg config.Config
}

func New(cfg config.Config) *Handler {
	return &Handler{
		cfg: cfg,
	}
}

func (h *Handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	log.Println("Received DNS request...")

	msg := dns.Msg{}
	msg.SetReply(r)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name
		address, ok := h.fetchAddress(domain)

		if ok {
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{ Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60 },
				A: net.ParseIP(address),
			})

			log.Printf("answering with %s, for DNS domain %s", address, domain)
		} else {
			log.Printf("unable to find appropriate address for DNS domain %s", domain)
		}
	}

	w.WriteMsg(&msg)
}

func (h *Handler) fetchAddress(domain string) (string, bool) {
	for _, alias := range h.cfg.Aliases {
		if alias.Domain == domain || alias.Domain+"." == domain {
			return alias.Address, true
		}
	}

	return "", false
}