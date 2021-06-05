package main

import (
	"fmt"
	"log"
	"sync/atomic"
	"net/http"
	"net/url"
	"net/http/httputil"
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
	Current uint64
}

func (sp *ServerPool) Add(u string){

		var backend Backend
		Url,err := url.Parse(u)
		checkerr(err)

		ReverseProxy := httputil.NewSingleHostReverseProxy(Url)
		
		backend.Url = Url
		backend.Online = true
		backend.ReverseProxy = ReverseProxy
		sp.Backends = append(sp.Backends, &backend)
}
func (sp *ServerPool) Get() *httputil.ReverseProxy{
	index := int(atomic.AddUint64(&sp.Current, uint64(1)) % uint64(len(sp.Backends)))
	sp.Backends[index].Total += 1
	fmt.Println(sp.Backends[0].Total,sp.Backends[1].Total,sp.Backends[2].Total,sp.Backends[3].Total,sp.Backends[4].Total)
	return sp.Backends[index].ReverseProxy
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
	
	http.ListenAndServe(":8080", http.HandlerFunc(LoadBalancing))

}