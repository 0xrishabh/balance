package balance

/*
	I have the config file
		- host {
			google.com {
				8.8.8.8
				4.4.4.4
			}
			yahoo.com {
				8.8.8.8
				4.4.4.4
			}
		}
*/

import (
	"sync/atomic"
	"net/http"
	"net/url"
	"net/http/httputil"
	"net"
)
type Server struct {
	Url *url.URL
	Online bool
	ReverseProxy *httputil.ReverseProxy
}

type ServerPool struct {
	Servers []*Server
	AliveServers []*Server
	Current uint64
}

func (sp *ServerPool) Add(u string){
	Url,err := url.Parse(u)
	checkerr(err)

	ReverseProxy := httputil.NewSingleHostReverseProxy(Url)
	sp.Servers = append(sp.Servers, &Server{Url, true, ReverseProxy, 0})
}

func (sp *ServerPool) Get() *httputil.ReverseProxy{
	index := int(atomic.AddUint64(&sp.Current, uint64(1)) % uint64(len(sp.Working)))
	sp.Working[index].Total += 1
	return sp.Working[index].ReverseProxy
}

var Router = make(map[string]ServerPool)

func Config(config map[string]interface{}) func(w http.ResponseWriter, r *http.Request){
	for k,v := range config{
		var Server ServerPool
		for _,ip := range v{
			Server.Add(ip)
		}
		Router[k] = Server
	}
	return Handler
}

func Handler(w http.ResponseWriter, r *http.Request){
	if Router[r.Host] == ""{
		fmt.Println("host does not exists yet: ", r.Host)
	}
	Router[r.Host].ServeHTTP(w,r)
}