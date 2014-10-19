package main

import (
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
// TODO: Fix this mess when conscious
func (r *EtcdRegistry) registerattributes(service *Service) {
	log.Println("Called registerattributes with", service.Attrs)
	base_path := r.path + "/" + service.Name + "/" + service.ID
	if mode, prs := service.Attrs["mode"]; prs {
		mode = strings.ToLower(mode)
		if mode != "http" && mode != "tcp" {
			log.Println("Unknown mode requested, defaulting to TCP")
			mode = "tcp"
		}
		log.Println("Setting network mode", mode, "for service", service.Name)
		_, err := r.client.Set(base_path+"/mode", mode, uint64(0))
		if err != nil {
			log.Println("Unable to set network mode:", err, "ignoring service.")
			return
		}
		if val, prs := service.Attrs["external_port"]; prs && mode == "tcp" {
			_, err := r.client.Set(base_path+"/external_port", val, uint64(0))
			if err != nil {
				log.Println("Unable to set external port:", err)
			}
		} else if val, prs := service.Attrs["internal_port"]; prs && mode == "tcp" {
			_, err := r.client.Set(base_path+"/internal_port", val, uint64(0))
			if err != nil {
				log.Println("Unable to set internal port:", err)
			}
		}
		if val, prs := service.Attrs["domains"]; prs && mode == "http" {
			domains := strings.Split(val, ",")
			if len(domains) > 0 {
				log.Println("Domains:", domains)
				r.client.SetDir(base_path+"/domains", uint64(0))
				for _, domain := range domains {
					_, err := r.client.Set(base_path+"/domains/"+domain, domain, uint64(0))
					if err != nil {
						log.Println("Unable to create domain:", err)
					}
				}
			}
		}
	} else {
		log.Println("Service mode not specified, default to TCP for service", service.Name)
		r.client.Set(base_path+"/mode", "tcp", uint64(0))
	}
}

func (r *EtcdRegistry) Register(service *Service) error {
	log.Println("register called for:", service.Name)
	base_path := r.path + "/" + service.Name + "/" + service.ID
	r.registerattributes(service)
	port := strconv.Itoa(service.Port)
	r.client.Set(base_path+"/port", port, uint64(0))
	r.client.Set(base_path+"/address", service.IP, uint64(0))
	return nil
}

func (r *EtcdRegistry) Deregister(service *Service) error {
	path := r.path + "/" + service.Name + "/" + service.ID
	_, err := r.client.Delete(path, false)
	return err
}
