package main

import (
	"fmt"
	"log"
	"sync/atomic"
	"net/http"
	"net/url"
	"net/http/httputil"
	"net"
	"time"
	"github.com/0xrishabh/balance/src/config"
	"github.com/0xrishabh/balance/src/balance"
)

func checkerr(err error)  {
	if err != nil{
		log.Fatal(err)
	}
}

type Backend struct {
	Url *url.URL
	Online bool
	ReverseProxy *httputil.ReverseProxy
	Total uint64
}

type ServerPool struct {
	Backends []*Backend
	Working []*Backend
	Current uint64
}

func (backend *Backend) Health() bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", backend.Url.Host, timeout)
	if err != nil {

		log.Println("Site unreachable, error: ", err)
		return false
	}
	
	defer conn.Close()
	return true
}

func (sp *ServerPool) Monitor() {
	var status []bool
  
	for _,backend := range sp.Backends{
		backend.Online = backend.Health()
		status = append(status, backend.Online)
	}

	fmt.Println(status)
}

func (sp *ServerPool) Add(u string){

		Url,err := url.Parse(u)
		checkerr(err)

		ReverseProxy := httputil.NewSingleHostReverseProxy(Url)
		sp.Backends = append(sp.Backends, &Backend{Url, true, ReverseProxy, 0})
}

func (sp *ServerPool) Get() *httputil.ReverseProxy{
	index := int(atomic.AddUint64(&sp.Current, uint64(1)) % uint64(len(sp.Working)))
	sp.Working[index].Total += 1
	return sp.Working[index].ReverseProxy
}	

var peers ServerPool

func LoadBalancing(w http.ResponseWriter, r *http.Request){
	
	peers.Get().ServeHTTP(w,r)

}

func main(){
	
	peers.Add("http://0.0.0.0:50001")
	peers.Add("http://0.0.0.0:50002")
	peers.Add("http://0.0.0.0:50003")
	peers.Add("http://0.0.0.0:50004")
	peers.Add("http://0.0.0.0:50005")
	peers.Add("http://0.0.0.0:5000")
	
	config = configuration.Load("/tmp/config.yml")
	LoadBalancerHandler := balance.Run(config)
	
	go func(){
		for{
			peers.Monitor()
			time.Sleep(5 * time.Second)
		}
	}()
	http.ListenAndServe(":8080", http.HandlerFunc(LoadBalancerHandler))

}

