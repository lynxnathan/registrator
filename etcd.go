package main

import (
	"net"
	"net/url"
	"strconv"
	"strings"
	"log"
	"github.com/coreos/go-etcd/etcd"
)

type EtcdRegistry struct {
	client *etcd.Client
	path   string
}

func NewEtcdRegistry(uri *url.URL) ServiceRegistry {
	urls := make([]string, 0)
	if uri.Host != "" {
		urls = append(urls, "http://"+uri.Host)
	}
	return &EtcdRegistry{client: etcd.NewClient(urls), path: uri.Path}
}

func (r *EtcdRegistry) registerdomains(service *Service) {
	log.Println("Called registerdomains")
	path := r.path + "/" + service.Name + "/" + service.ID + "/domains"
	if val, prs := service.Attrs["domains"]; prs {
		domains := strings.Split(val, ",")
		if len(domains) > 0 {
			log.Println("Domains:", domains)
			r.client.SetDir(path, uint64(0))
			for _, domain := range domains {
				_, err := r.client.Set(path + "/" + domain, domain, uint64(0))
				if err != nil {
					log.Println("Unable to create domain:", err)
				}
			}
		}
	}
}

func (r *EtcdRegistry) Register(service *Service) error {
	log.Println("register called for:", service.Name)
	path := r.path + "/" + service.Name + "/" + service.ID + "/endpoint"
	r.registerdomains(service)
	port := strconv.Itoa(service.Port)
	addr := net.JoinHostPort(service.IP, port)
	_, err := r.client.Set(path, addr, uint64(0))
	return err
}

func (r *EtcdRegistry) Deregister(service *Service) error {
	path := r.path + "/" + service.Name + "/" + service.ID
	_, err := r.client.Delete(path, false)
	return err
}
