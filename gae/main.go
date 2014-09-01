package home

import (
    "net/http"
    "fmt"
    //"encoding/gob"
    //"encoding/json"
    "strconv"
    "github.com/pubnub/go/messaging"
    "github.com/gorilla/mux"
    "github.com/gorilla/sessions"
    //"time"
    "html/template"
    "log"
    "appengine"
    "appengine/channel"
    
    //"appengine/module"
    //"appengine/runtime"
    //"appengine/delay"
)

var mainTemplate = template.Must(template.ParseFiles("main.html"))
var Store = sessions.NewCookieStore([]byte("sess-secret-key"))
var pubnubInstMap map[string]interface{}
var subscribeKey = "pam"
var publishKey = "pam"
var secretKey = "pam"

func GetSessionsOptionsObject(age int) (* sessions.Options){
	return &sessions.Options{
			    Path:     "/",
			    MaxAge:   age,
			    HttpOnly: true,
			}		    
}

func init() {
	router := mux.NewRouter()
	router = router.StrictSlash(true)
	router.HandleFunc("/", Handler)
	router.HandleFunc("/subscribe", Subscribe)
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
    router.HandleFunc(`/{rest:[a-zA-Z0-9=\-\/]+}`, Handler)
    
    http.Handle("/", router)
}

func Signout(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	uuid := q.Get("uuid")
	c := appengine.NewContext(r)
	messaging.SetContext(c)
	
	session, err := Store.Get(r, "user-session")
    if(err == nil && 
    	session !=nil && 
    	session.Values["uuid"] != nil) {
    	//session.Values["pubInst"] != nil) {
    	
    	if val, ok := session.Values["uuid"].(string); ok {
    		c.Infof("Deleteing Session ok1 %s", val)
    	
    		delete(pubnubInstMap, val)
		} else {
			delete(pubnubInstMap, uuid)
		}
		session.Values["uuid"] = ""
		session.Options = GetSessionsOptionsObject(-1)
		session.Save(r, w)
		c.Infof("Deleted Session %s")
	} 
}

func GetUserState(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ch := q.Get("ch")
	uuid := q.Get("uuid")
	pubInstance := initPubnub(uuid, w, r) 
	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	
	go pubInstance.GetUserState(ch, successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "Get User State")
}

func DeleteUserState(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ch := q.Get("ch")
	key := q.Get("k")
	uuid := q.Get("uuid")
	pubInstance := initPubnub(uuid, w, r) 
	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	
	go pubInstance.SetUserStateKeyVal(ch, key, "", successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "Del User State")
}

func SetUserStateJson(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ch := q.Get("ch")
	j := q.Get("j")
	uuid := q.Get("uuid")
	pubInstance := initPubnub(uuid, w, r) 
	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	
	go pubInstance.SetUserStateJSON(ch, j, successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "Set User State JSON")
}

func SetUserState(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ch := q.Get("ch")
	k := q.Get("k")
	v := q.Get("v")
	uuid := q.Get("uuid")
	pubInstance := initPubnub(uuid, w, r) 
	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	
	go pubInstance.SetUserStateKeyVal(ch, k, v, successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "Set User State")
}

func AuditPresence(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ch := q.Get("ch")
	uuid := q.Get("uuid")
	pubInstance := initPubnub(uuid, w, r) 
	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	go pubInstance.AuditPresence(ch, successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "Audit Presence")
}

func RevokePresence(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ch := q.Get("ch")
	uuid := q.Get("uuid")
	pubInstance := initPubnub(uuid, w, r) 
	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	go pubInstance.GrantPresence(ch, false, false, 0, successChannel, errorChannel)
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
	
	pubInstance := initPubnub(uuid, w, r) 
	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	go pubInstance.GrantPresence(ch, bRead, bWrite, iTtl, successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "Revoke Presence")

}

func AuditSubscribe(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ch := q.Get("ch")
	uuid := q.Get("uuid")
	pubInstance := initPubnub(uuid, w, r) 
	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	go pubInstance.AuditSubscribe(ch, successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "Audit Subscribe")
}

func RevokeSubscribe(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ch := q.Get("ch")
	uuid := q.Get("uuid")
	pubInstance := initPubnub(uuid, w, r) 
	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	go pubInstance.GrantSubscribe(ch, false, false, 0, successChannel, errorChannel)
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
	pubInstance := initPubnub(uuid, w, r) 
	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	go pubInstance.GrantSubscribe(ch, bRead, bWrite, iTtl, successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "Revoke Subscribe")
}

func SetAuthKey(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	authKey := q.Get("authkey")
	uuid := q.Get("uuid")
	pubInstance := initPubnub(uuid, w, r) 
	pubInstance.SetAuthenticationKey(authKey)
	SendResponseToChannel("Auth key set", r, uuid);
}

func GetAuthKey(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	uuid := q.Get("uuid")
	pubInstance := initPubnub(uuid, w, r) 
	SendResponseToChannel("Auth key: "+pubInstance.GetAuthenticationKey(), r, uuid);
}

func Publish(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	message := q.Get("m")
	uuid := q.Get("uuid")
	
	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	
	pubInstance := initPubnub(uuid, w, r) 
	go pubInstance.Publish("test", message, successChannel, errorChannel)
	
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
	
	pubInstance := initPubnub(uuid, w, r) 
	go pubInstance.GlobalHereNow(disableUUID, includeUserState, successChannel, errorChannel)
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
	
	pubInstance := initPubnub(uuid, w, r) 
	go pubInstance.HereNow(channel, disableUUID, includeUserState, successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "HereNow")
}

func WhereNow(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	whereNowUUID := q.Get("whereNowUUID")
	uuid := q.Get("uuid")

	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	
	pubInstance := initPubnub(uuid, w, r) 
	go pubInstance.WhereNow(whereNowUUID, successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, errorChannel, messaging.GetNonSubscribeTimeout(), "WhereNow")
}

func Time(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	uuid := q.Get("uuid")

	errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	
	pubInstance := initPubnub(uuid, w, r) 
	go pubInstance.GetTime(successChannel, errorChannel)
	handleResult(w, r, uuid, successChannel, successChannel, messaging.GetNonSubscribeTimeout(), "Time")
}

func initPubnub(uuid string, w http.ResponseWriter, r *http.Request) *messaging.Pubnub{
	cipher:= ""
	ssl:= false

	c := appengine.NewContext(r)
	messaging.SetContext(c)
	
	session, err := Store.Get(r, "user-session")
    if(err == nil && 
    	session !=nil && 
    	session.Values["uuid"] != nil) {
    	//session.Values["pubInst"] != nil) {
    	
    	if val, ok := session.Values["uuid"].(string); ok {
    		c.Infof("Session ok1 %s", val)
    		var pubInst2 = pubnubInstMap[val]
    		if(pubInst2 != nil){
    			return pubInst2.(*messaging.Pubnub)
    		} else {
	    		c.Errorf("pubInst2 val nil")
    		}
	/*var pubInst3 *messaging.Pubnub
	pubInst3 = session.Values["pubInst"].(*messaging.Pubnub)
	
	c.Infof("pubInst4:%s",pubInst3.GetUUID())
    		
		if val, ok := session.Values["pubInst"].(messaging.PubnubRaw); ok {
    		c.Infof("Session ok1 %s", val)
	var pubInst3 *messaging.Pubnub
	pubInst3 = &messaging.Pubnub{ PN : val }
	//intd := session.Values["pubInst"].(map[string]interface{})
	//pubInst3 = intd[uuidn].(*messaging.Pubnub)
	
	c.Infof("pubInst3:%s",pubInst3.GetUUID())
	c.Infof("pubInst3 Auth:%s",pubInst3.GetAuthenticationKey())
	c.Infof("pubInst3 newPubnub.PN.Origin:%s",pubInst3.PN.Origin)
	c.Infof("pubInst3 newPubnub.PN.Auth:%s",pubInst3.PN.AuthenticationKey)
	
    		
    		
    		if(pubInst3 != nil){
    			return pubInst3
    		} else {
	    		c.Errorf("pubInst2 val nil")
    		}*/    		
    	} else {
    		c.Errorf("Session val nil")
    	}
	} else {
		if(err != nil){
			c.Errorf("Session error:", err.Error())
		}
		if(session == nil){
			c.Errorf("Session nil")
		}
		/*if(session.Values["pubInst"] == nil){
			c.Errorf("pubInst nil")
		}*/
		if(session.Values["uuid"] == nil){
			c.Errorf("uuid nil")
		}
	}
	c.Infof("Creating instance")
	var pubInstance = messaging.NewPubnub(publishKey, subscribeKey, secretKey, cipher, ssl, uuid)
	//
	if(pubnubInstMap == nil){
		pubnubInstMap = make (map[string]interface{})
	}
	uuidn := pubInstance.GetUUID()
	
	/*jsonSerialized, err := json.Marshal(pubInstance)
	if err != nil {
		c.Errorf("err, %s", err.Error())
	} else {
		c.Infof("jsonSerialized:%v",jsonSerialized)
	}*/
	
	//var pubInterface interface{}
	//pubInterface = pubInstance
	session.Values["uuid"] = uuidn
	pubnubInstMap[uuidn] = interface{}(pubInstance)
	//pubnubInstMap["0"] = pubInterface
	//gob.Register(pubnubInstMap)
	//gob.Register(pubInstance.PN)
	//session.Values["pubInst"] = pubInstance.PN//pubInterface
	
	//c.Infof("uuid: %s %s",uuidn, pubInstance.GetUUID())
	session.Options = GetSessionsOptionsObject(60*20)
	err1 := session.Save(r, w)
	if(err1!=nil){
		c.Errorf("error1, %s", err1.Error())
	}
	
	//var pubInst3 *messaging.Pubnub
	//pubInst3 = &messaging.Pubnub{ PN : session.Values["pubInst"].(messaging.PubnubRaw) }
	//intd := session.Values["pubInst"].(map[string]interface{})
	//pubInst3 = intd[uuidn].(*messaging.Pubnub)
	
	//c.Infof("pubInst3:%s",pubInst3.GetUUID())
	/*errUnmarshalMessages := json.Unmarshal([]byte(session.Values["pubInst"].(string)), &pubInst3)

	if errUnmarshalMessages != nil {
		c.Errorf("errUnmarshalMessages, %s", errUnmarshalMessages.Error())
	} else {
		c.Infof("pubInst3:%s",pubInst3.GetUUID())
	}*/
	var pubInst2 = pubnubInstMap[session.Values["uuid"].(string)]
	if val, ok := session.Values["uuid"].(string); ok {
    		c.Infof("Session ok %s", val)
    		pubInst2 = pubnubInstMap[val]
    		if(pubInst2!=nil){
    			return pubInst2.(*messaging.Pubnub)
    		} else {
    			return nil
    		}
    	} else {
    		c.Errorf("Session val nil")
    	}
	return nil
}

func Subscribe(w http.ResponseWriter, r *http.Request) {
	
	/*errorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	
	c := appengine.NewContext(r)*/
	/*t := taskqueue.NewPOSTTask("/worker", url.Values{
            // ...
        })
        // Use the transaction's context when invoking taskqueue.Add.
        _, err := taskqueue.Add(c, t, "")
t := &taskqueue.Task{Name: "foo"}
err := taskqueue.Delete(c, t, "queue1")	
	
	//err := runtime.RunInBackground(c, func(c appengine.Context) {
	err := runtime.RunInBackground(c, func(c appengine.Context) {*/
/*var laterFunc = delay.Func("key", func(c appengine.Context, x string) {
    // ...
    pubInstance := initPubnub("customuuid") 
   	pubInstance.Subscribe("test", "", successChannel, false, errorChannel)
	handleResultSubscribe(w, r, successChannel, errorChannel, 310, "Subscribe")
    
})	
laterFunc.Call(c, "")*/
	  //go 
	//})		
	/*if err != nil {
        c.Errorf("Background process: %v", err)
        return
    }*/	
	    
}

func Handler(w http.ResponseWriter, r *http.Request) {
	uuid := ""
	c := appengine.NewContext(r)
	messaging.SetContext(c)
	
	pubInstance :=  initPubnub(uuid, w, r)
	if(pubInstance == nil){
		c.Errorf("Couldn't create pubnub instance")
		http.Error(w, "Couldn't create pubnub instance", http.StatusInternalServerError)
		return
	}
	nuuid := pubInstance.GetUUID()
	
 	tok, err := channel.Create(c, nuuid)
    if err != nil {
        http.Error(w, "Couldn't create Channel", http.StatusInternalServerError)
        c.Errorf("channel.Create: %v", err)
        return
    }	
	
	err1 := mainTemplate.Execute(w, map[string]string{
        "token":    tok,
        "uuid":       nuuid,
        "subscribeKey": subscribeKey,
        "publishKey": publishKey,
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

func SendResponseToChannel(message string, r *http.Request, uuid string){
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
				SendResponseToChannel(string(success), r, uuid);
			}
			flush(w)
		case failure, ok := <-errorChannel:
			if !ok {
				log.Println(fmt.Sprintf("failure!OK"))
			}
			if string(failure) != "[]" {
				SendResponseToChannel(string(failure), r, uuid);
			}
		}
	}
}

func handleResult(w http.ResponseWriter, r *http.Request, uuid string, successChannel, errorChannel chan []byte, timeoutVal uint16, action string) {
	for {
		select {
		
		case success, ok := <-successChannel:
			if !ok {
				log.Println(fmt.Sprintf("success!OK"))
				break
			}
			if string(success) != "[]" {
				SendResponseToChannel(string(success), r, uuid);
			}
			if f, ok := w.(http.Flusher); ok {
			    f.Flush()
			} else {
   				// Response writer does not support flush.
   				//fmt.Fprintf(w, fmt.Sprintf("<p> Response writer does not support flush.:</p>"))
			}
			
			return
		case failure, ok := <-errorChannel:
			if !ok {
				log.Println(fmt.Sprintf("failure!OK"))
				break
			}
			if string(failure) != "[]" {
				SendResponseToChannel(string(failure), r, uuid);
			}
			return
		}
	}
}