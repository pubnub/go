package e2e

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"

	mux "github.com/gorilla/mux"
	pubnub "github.com/pubnub/go"
)

const (
	SPECIAL_CHARACTERS = "-.,_~:/?#[]@!$&'()*+;=`|"
	SPECIAL_CHANNEL    = "-._~:/?#[]@!$&'()*+;=`|"
)

var pamConfig *pubnub.Config
var config *pubnub.Config

var (
	serverErrorTemplate     = "pubnub/server: Server respond with error code %d"
	validationErrorTemplate = "pubnub/validation: %s"
	connectionErrorTemplate = "pubnub/connection: %s"
)

func init() {
	rand.Seed(time.Now().UnixNano())

	config = pubnub.NewConfig()
	config.PublishKey = "pub-c-071e1a3f-607f-4351-bdd1-73a8eb21ba7c"
	config.SubscribeKey = "sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f"

	pamConfig = pubnub.NewConfig()
	pamConfig.PublishKey = "pub-c-1bd448ed-05ba-4dbc-81a5-7d6ff5c6e2bb"
	pamConfig.SubscribeKey = "sub-c-90c51098-c040-11e5-a316-0619f8945a4f"
	pamConfig.SecretKey = "sec-c-ZDA1ZTdlNzAtYzU4Zi00MmEwLTljZmItM2ZhMDExZTE2ZmQ5"
}

func configCopy() *pubnub.Config {
	cfg := new(pubnub.Config)
	*cfg = *config
	return cfg
}

func pamConfigCopy() *pubnub.Config {
	config := new(pubnub.Config)
	*config = *pamConfig
	return config
}

func randomized(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, rand.Intn(10000000))
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
		time.Sleep(2000 * time.Millisecond)
		closed <- true
	}()

	s.Serve(l)
}
