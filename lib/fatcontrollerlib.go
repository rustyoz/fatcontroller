package fatcontroller

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/spf13/viper"
)

type VirtualHost struct {
	Name string
	Url  string
	Host string
	Path string
	Port int
}

var Proxies []*httputil.ReverseProxy

var Port int
var Host string
var VirtualHosts []VirtualHost
var LogRequests bool

func Run() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("$HOME/.fatcontroller")
	viper.AddConfigPath(".")
	configname := flag.String("config", "config", "path/configfilename")
	flag.BoolVar(&LogRequests, "logrequests", false, "log each http request")
	flag.Parse()

	viper.SetConfigName(*configname)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error in config file %s", err))
	}
	Port = viper.GetInt("port")
	Host = viper.GetString("host")
	vhs := viper.Get("virtualhost")

	VirtualHosts = make([]VirtualHost, len(vhs.([]map[string]interface{})))
	for n, vh := range vhs.([]map[string]interface{}) {
		VirtualHosts[n].Name = vh["name"].(string)
		VirtualHosts[n].Url = vh["url"].(string)
		VirtualHosts[n].Port = int(vh["port"].(int64))
		VirtualHosts[n].Host = vh["host"].(string)
		VirtualHosts[n].Path = vh["path"].(string)
	}

	mux, err := ShuntingYard(VirtualHosts)

	if err != nil {
		fmt.Println(err)
	}

	addr := net.JoinHostPort(Host, strconv.Itoa(Port))

	if LogRequests {
		log.Fatal(http.ListenAndServe(addr, Log(mux)))
	} else {
		log.Fatal(http.ListenAndServe(addr, mux))
	}

}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func ShuntingYard(vhs []VirtualHost) (mux *http.ServeMux, err error) {
	mux = http.NewServeMux()
	for _, vh := range vhs {

		target, e := url.Parse(vh.Url)
		if e != nil {
			return nil, e
		}
		target.Scheme = "http"
		target.Path = vh.Path
		target.Host = net.JoinHostPort(vh.Host, strconv.Itoa(vh.Port))

		Proxies = append(Proxies, httputil.NewSingleHostReverseProxy(target))
		fmt.Println(vh.Name)
		fmt.Printf("From %s:%d%s\n", Host, Port, vh.Url)
		fmt.Printf("To   %s:%d%s\n\n", vh.Host, vh.Port, vh.Path)

	}

	for n, p := range Proxies {
		mux.Handle(vhs[n].Url, http.StripPrefix(vhs[n].Url, p))
	}

	return mux, err
}
