package pntests

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	mux "github.com/gorilla/mux"
	pubnub "github.com/pubnub/go"
)

var pamConfig *pubnub.Config

var (
	serverErrorTemplate     = "pubnub/server: Server respond with error code %d"
	validationErrorTemplate = "pubnub/validation: %s"
	connectionErrorTemplate = "pubnub/connection: %s"
)

func init() {
	pamConfig = pubnub.NewConfig()
	pamConfig.PublishKey = "pub-c-1bd448ed-05ba-4dbc-81a5-7d6ff5c6e2bb"
	pamConfig.SubscribeKey = "sub-c-90c51098-c040-11e5-a316-0619f8945a4f"
}

func pamConfigCopy() *pubnub.Config {
	config := new(pubnub.Config)
	*config = *pamConfig
	return config
}

func makeResponseRoot(hangSeconds int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		log.Printf("Sleeping %d seconds\n", hangSeconds)
		time.Sleep(time.Duration(hangSeconds) * time.Second)

		if vars["pubKey"] == "my_pub_key" {
			fmt.Fprint(w, "[1, \"Sent\", 123]")
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "[{\"eror\": true}]")
		}
	}
}

func servePublish(hangSeconds int, close, closed chan bool) {
	r := mux.NewRouter()
	r.HandleFunc("/publish/{pubKey}/{subKey}/0/{channel}/0/{msg}",
		makeResponseRoot(hangSeconds))

	s := &http.Server{
		Handler: r,
	}

	l, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}

	go func() {
		<-close
		fmt.Println(">>> closing listener")
		l.Close()
		fmt.Println("<<< listener closed")

		// HACK: let the server release resources before a next test started
		time.Sleep(2000 * time.Millisecond)
		closed <- true
	}()

	s.Serve(l)
}
