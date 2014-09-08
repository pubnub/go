package home

import (
    "net/http"
    "fmt"
    "math/big"
    "strconv"
    "strings"
    "github.com/pubnub/go/gae/messaging"
    "github.com/gorilla/mux"
    "github.com/gorilla/sessions"
    "html/template"
    "log"
    "appengine"
    "appengine/channel"
    //"time"
    //"appengine/taskqueue"
)

var mainTemplate = template.Must(template.ParseFiles("main.html"))
var subscribeKey = "demo"
var publishKey = "demo"
var secretKey = "demo"

var Store = sessions.NewCookieStore([]byte(secretKey))

func init() {
    router := mux.NewRouter()
    router = router.StrictSlash(true)
    router.HandleFunc("/", Handler)
    //router.HandleFunc("/subscribe", Subscribe)
    //router.HandleFunc("/worker", worker)
    router.HandleFunc("/publish", Publish)
    router.HandleFunc("/globalHereNow", GlobalHereNow)
    router.HandleFunc("/hereNow", HereNow)
    router.HandleFunc("/whereNow", WhereNow)
    router.HandleFunc("/time", Time)
    router.HandleFunc("/setAuthKey", SetAuthKey)
    router.HandleFunc("/getAuthKey", GetAuthKey)
    router.HandleFunc("/deleteUserState", DeleteUserState)
    router.HandleFunc("/setUserStateJson", SetUserStateJson)
    router.HandleFunc("/setUserState", SetUserState)
    router.HandleFunc("/auditPresence", AuditPresence)
    router.HandleFunc("/revokePresence", RevokePresence)
    router.HandleFunc("/grantPresence", GrantPresence)
    router.HandleFunc("/auditSubscribe", AuditSubscribe)
    router.HandleFunc("/revokeSubscribe", RevokeSubscribe)
    router.HandleFunc("/grantSubscribe", GrantSubscribe)
    router.HandleFunc("/getUserState", GetUserState)	
    router.HandleFunc("/signout", Signout)
    router.HandleFunc("/connect", Connect)
    router.HandleFunc("/keepAlive", KeepAlive)
    router.HandleFunc("/detailedHistory", DetailedHistory)	
    router.HandleFunc(`/{rest:[a-zA-Z0-9=\-\/]+}`, Handler)
    
    http.Handle("/", router)
}

func DetailedHistory(w http.ResponseWriter, r *http.Request){
	q := r.URL.Query()
	ch := q.Get("ch")
	uuid := q.Get("uuid")
	start := q.Get("start")
	var iStart int64 = 0
	if(strings.TrimSpace(start) != ""){
		bi := big.NewInt(0)
		if _, ok := bi.SetString(start, 10); !ok {
			iStart = 0
		} else {
			iStart = bi.Int64()
		}
	}
	
	end := q.Get("end")
	var iEnd int64 = 0
	if(strings.TrimSpace(end) != ""){
		bi := big.NewInt(0)
		if _, ok := bi.SetString(end, 10); !ok {
			iEnd = 0
		} else {
			iEnd = bi.Int64()
		}
	}
	
	limit := q.Get("limit")
	reverse := q.Get("reverse")
	
	iLimit := 100
	if ival, err := strconv.Atoi(limit); err == nil {
		iLimit = ival
	}
	
	bReverse := false
	if(reverse == "1"){
		bReverse = true	
	} 
	
	context := appengine.NewContext(r)
	
	pubInstance := messaging.New(context, uuid, w, r, publishKey, subscribeKey, secretKey, "", false) 
	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	
	go pubInstance.History(context, w, r, ch, iLimit, iStart, iEnd, bReverse, successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "Detailed History")
}

func Connect(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	pubKey := q.Get("pubKey")
	subKey := q.Get("subKey")
	uuid := q.Get("uuid")
	secKey := q.Get("secKey")
	cipher := q.Get("cipher")
	ssl := q.Get("ssl")
	bSsl := false
	if(ssl =="1"){
		bSsl=true 
	}
	
	c := appengine.NewContext(r)
	messaging.SetSessionKeys (c, w, r, pubKey, subKey, secKey, cipher, bSsl, uuid);
	
	session, err := Store.Get(r, "example-session")
	if(err!=nil){
		c.Errorf("Session store error %s", err.Error())
		http.Error(w, "Session store error", http.StatusInternalServerError)
		return		
	}
	
	tok, err := channel.Create(c, uuid)
 	
    if err != nil {
        http.Error(w, "Couldn't create Channel", http.StatusInternalServerError)
        c.Errorf("channel.Create: %v", err)
        return
    } else {
    	session.Values["token"] = tok
		err1 := session.Save(r, w)
		if(err1 != nil){
			c.Errorf("error1, %s", err1.Error())
		}    	
    }	
    fmt.Fprintf(w, tok)	
}


func KeepAlive(w http.ResponseWriter, r *http.Request) {

}

func Signout(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	messaging.DeleteSession(c, w, r, secretKey)
}

func GetUserState(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ch := q.Get("ch")
	uuid := q.Get("uuid")
	
	c := appengine.NewContext(r)
	pubInstance := messaging.New(c, uuid, w, r, publishKey, subscribeKey, secretKey, "", false) 
	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	
	go pubInstance.GetUserState(c, w, r, ch, successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "Get User State")
}

func DeleteUserState(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ch := q.Get("ch")
	key := q.Get("k")
	uuid := q.Get("uuid")
	c := appengine.NewContext(r)
	pubInstance := messaging.New(c, uuid, w, r, publishKey, subscribeKey, secretKey, "", false) 
	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	
	go pubInstance.SetUserStateKeyVal(c, w, r, ch, key, "", successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "Del User State")
}

func SetUserStateJson(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ch := q.Get("ch")
	j := q.Get("j")
	uuid := q.Get("uuid")
	c := appengine.NewContext(r)
	pubInstance := messaging.New(c, uuid, w, r, publishKey, subscribeKey, secretKey, "", false) 
	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	
	go pubInstance.SetUserStateJSON(c, w, r, ch, j, successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "Set User State JSON")
}

func SetUserState(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ch := q.Get("ch")
	k := q.Get("k")
	v := q.Get("v")
	uuid := q.Get("uuid")
	c := appengine.NewContext(r)
	pubInstance := messaging.New(c, uuid, w, r, publishKey, subscribeKey, secretKey, "", false) 
	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	
	//setUserState
	
	go pubInstance.SetUserStateKeyVal(c, w, r, ch, k, v, successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "Set User State")
}

func AuditPresence(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ch := q.Get("ch")
	uuid := q.Get("uuid")
	c := appengine.NewContext(r)
	pubInstance := messaging.New(c, uuid, w, r, publishKey, subscribeKey, secretKey, "", false) 
	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	go pubInstance.AuditPresence(c, w, r, ch, successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "Audit Presence")
}

func RevokePresence(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ch := q.Get("ch")
	uuid := q.Get("uuid")
	c := appengine.NewContext(r)
	pubInstance := messaging.New(c, uuid, w, r, publishKey, subscribeKey, secretKey, "", false) 
	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	go pubInstance.GrantPresence(c, w, r, ch, false, false, 0, successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "Revoke Presence")
}

func GrantPresence(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ch := q.Get("ch")
	uuid := q.Get("uuid")
	read := q.Get("r")
	write := q.Get("w")
	ttl := q.Get("ttl")
	bRead := false
	if(read =="1"){
		bRead=true 
	}
	bWrite := false
	if(write =="1"){
		bWrite=true 
	}
	iTtl := 1440
	if ival, err := strconv.Atoi(ttl); err == nil {
		iTtl = ival
	} 
	
	c := appengine.NewContext(r)
	pubInstance := messaging.New(c, uuid, w, r, publishKey, subscribeKey, secretKey, "", false) 
	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	go pubInstance.GrantPresence(c, w, r, ch, bRead, bWrite, iTtl, successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "Revoke Presence")

}

func AuditSubscribe(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ch := q.Get("ch")
	uuid := q.Get("uuid")
	c := appengine.NewContext(r)
	pubInstance := messaging.New(c, uuid, w, r, publishKey, subscribeKey, secretKey, "", false) 
	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	go pubInstance.AuditSubscribe(c, w, r, ch, successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "Audit Subscribe")
}

func RevokeSubscribe(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ch := q.Get("ch")
	uuid := q.Get("uuid")
	c := appengine.NewContext(r)
	pubInstance := messaging.New(c, uuid, w, r, publishKey, subscribeKey, secretKey, "", false) 
	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	go pubInstance.GrantSubscribe(c, w, r, ch, false, false, 0, successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "Revoke Subscribe")
}

func GrantSubscribe(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ch := q.Get("ch")
	read := q.Get("r")
	write := q.Get("w")
	ttl := q.Get("ttl")
	bRead := false
	if(read =="1"){
		bRead=true 
	}
	bWrite := false
	if(write =="1"){
		bWrite=true 
	}
	iTtl := 1440
	if ival, err := strconv.Atoi(ttl); err == nil {
		iTtl = ival
	} 
		
	uuid := q.Get("uuid")
	c := appengine.NewContext(r)
	pubInstance := messaging.New(c, uuid, w, r, publishKey, subscribeKey, secretKey, "", false) 
	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	go pubInstance.GrantSubscribe(c, w, r, ch, bRead, bWrite, iTtl, successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "Revoke Subscribe")
}

func SetAuthKey(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	authKey := q.Get("authkey")
	uuid := q.Get("uuid")
	
	c := appengine.NewContext(r)
	pubInstance := messaging.New(c, uuid, w, r, publishKey, subscribeKey, secretKey, "", false) 
	pubInstance.SetAuthenticationKey(c, w, r, authKey)
	SendResponseToChannel(w, "Auth key set", r, uuid);
}

func GetAuthKey(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	uuid := q.Get("uuid")
	c := appengine.NewContext(r)
	pubInstance := messaging.New(c, uuid, w, r, publishKey, subscribeKey, secretKey, "", false) 
	SendResponseToChannel(w, "Auth key: "+pubInstance.GetAuthenticationKey(), r, uuid);
}

func Publish(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	message := q.Get("m")
	uuid := q.Get("uuid")
	ch := q.Get("ch")
	
	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	
	c := appengine.NewContext(r)
	pubInstance := messaging.New(c, uuid, w, r, publishKey, subscribeKey, secretKey, "", false) 
	go pubInstance.Publish(c, w, r, ch, message, successChannel, errorChannel)
	
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "Publish")
}

func GlobalHereNow(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	uuid := q.Get("uuid")
	globalHereNowShowUUID := q.Get("showUUID")
	globalHereNowIncludeUserState := q.Get("includeUserState")	
	disableUUID := false;
	includeUserState := false;
	if(globalHereNowShowUUID == "1"){
		disableUUID = true;
	}
	if(globalHereNowIncludeUserState == "1"){
		includeUserState = true;
	}
	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	
	c := appengine.NewContext(r)
	pubInstance := messaging.New(c, uuid, w, r, publishKey, subscribeKey, secretKey, "", false) 
	go pubInstance.GlobalHereNow(c, w, r, disableUUID, includeUserState, successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "Global Here Now")
}

func HereNow(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	channel := q.Get("ch")
	uuid := q.Get("uuid")
	hereNowShowUUID := q.Get("showUUID")
	hereNowIncludeUserState := q.Get("includeUserState")
	
	disableUUID := false;
	includeUserState := false;
	if(hereNowShowUUID == "1"){
		disableUUID = true;
	}
	if(hereNowIncludeUserState == "1"){
		includeUserState = true;
	}

	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	
	c := appengine.NewContext(r)
	pubInstance := messaging.New(c, uuid, w, r, publishKey, subscribeKey, secretKey, "", false) 
	go pubInstance.HereNow(c, w, r, channel, disableUUID, includeUserState, successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "HereNow")
}

func WhereNow(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	whereNowUUID := q.Get("whereNowUUID")
	uuid := q.Get("uuid")

	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	
	c := appengine.NewContext(r)
	pubInstance := messaging.New(c, uuid, w, r, publishKey, subscribeKey, secretKey, "", false) 
	go pubInstance.WhereNow(c, w, r, whereNowUUID, successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "WhereNow")
}

func Time(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	uuid := q.Get("uuid")

	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	
	c := appengine.NewContext(r)
	pubInstance := messaging.New(c, uuid, w, r, publishKey, subscribeKey, secretKey, "", false) 
	go pubInstance.GetTime(c, w, r, successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "Time")
}

/*func worker(w http.ResponseWriter, r *http.Request) {
	uuid := r.FormValue("uuid")
	
    c := appengine.NewContext(r)
    for i:=0; i<10; i++{
    	time.Sleep(time.Duration(1) * time.Second)
    	c.Infof("looping:" + uuid)
    }
    c.Infof("message: %s", uuid)
    SendResponseToChannel(w, "message:" + uuid, r, uuid);
}

func Subscribe(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	uuid := q.Get("uuid")

	c := appengine.NewContext(r)
    t := taskqueue.NewPOSTTask("/worker", map[string][]string{"uuid": {uuid}})
    if _, err := taskqueue.Add(c, t, "subscribe"); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}*/

func Handler(w http.ResponseWriter, r *http.Request) {
	uuid := ""
	c := appengine.NewContext(r)
	pubInstance :=  messaging.New(c, uuid, w, r, publishKey, subscribeKey, secretKey, "", false)
	if(pubInstance == nil){
		c.Errorf("Couldn't create pubnub instance")
		http.Error(w, "Couldn't create pubnub instance", http.StatusInternalServerError)
		return
	}
	nuuid := pubInstance.GetUUID()
	
	session, err := Store.Get(r, "example-session")
	if(err!=nil){
		c.Errorf("Session store error %s", err.Error())
		http.Error(w, "Session store error", http.StatusInternalServerError)
		return		
	}
		
	//Enhancement: can be kept in memcache	
 	tok, err := channel.Create(c, nuuid)
    if err != nil {
        http.Error(w, "Couldn't create Channel", http.StatusInternalServerError)
        c.Errorf("channel.Create: %v", err)
        return
    } else {
    	session.Values["token"] = tok
    }	
	
	err1 := mainTemplate.Execute(w, map[string]string{
        "token":    tok,
        "uuid":       nuuid,
        "subscribeKey": subscribeKey,
        "publishKey": publishKey,
        "secretKey": secretKey,
    })
    if err1 != nil {
        c.Errorf("mainTemplate: %v", err1)
    }
}

func flush(w http.ResponseWriter) {
	f, ok := w.(http.Flusher)
	if ok && f != nil {
		log.Println("flush")
		f.Flush()
	} else {
   		// Response writer does not support flush.
   		fmt.Fprintf(w, fmt.Sprintf(" Response writer does not support flush.:"))
	}
	
}

/*func SendResponseToChannel(w http.ResponseWriter, message string, r *http.Request, uuid string){
	c := appengine.NewContext(r)
	c.Infof(message);
	fmt.Fprintf(w, fmt.Sprintf(message))
}*/

func SendResponseToChannel(w http.ResponseWriter, message string, r *http.Request, uuid string){
	c := appengine.NewContext(r)
	err := channel.SendJSON(c, uuid, message)
	c.Infof("json");
	if err != nil {
    	c.Errorf("sending Game: %v", err)
    }
}

func handleResultSubscribe(w http.ResponseWriter, r *http.Request, uuid string, successChannel, errorChannel chan []byte, timeoutVal uint16, action string) {
	for {
		select {
		
		case success, ok := <-successChannel:
			if !ok {
				log.Println(fmt.Sprintf("success!OK"))
				
			}
			if string(success) != "[]" {
				SendResponseToChannel(w, string(success), r, uuid);
			}
			flush(w)
		case failure, ok := <-errorChannel:
			if !ok {
				log.Println(fmt.Sprintf("failure!OK"))
			}
			if string(failure) != "[]" {
				SendResponseToChannel(w, string(failure), r, uuid);
			}
		}
	}
}

func handleResult(w http.ResponseWriter, r *http.Request, uuid string, successChannel, errorChannel chan []byte, timeoutVal uint16, action string) {
	c := appengine.NewContext(r)
	for {
		select {
		
		case success, ok := <-successChannel:
			if !ok {
				c.Infof("success!OK");				
				break
			}
			if string(success) != "[]" {
				c.Infof("success:",string(success));
				SendResponseToChannel(w, string(success), r, uuid);
			}
			
			return
		case failure, ok := <-errorChannel:
			if !ok {
				c.Infof("fail1:",string("failure"));
				break
			}
			if string(failure) != "[]" {
				c.Infof("fail:",string(failure));
				SendResponseToChannel(w, string(failure), r, uuid);
			}
			return
		}
	}
}