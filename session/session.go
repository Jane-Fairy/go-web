package session

import (
	"crypto/rand"
	"encoding/base64"
	"go-web/util"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Manager struct {
	cookieName  string
	lock        sync.Mutex
	Provider    Provider
	maxLifeTime int64
}

func (manager *Manager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	time.AfterFunc(time.Duration(manager.maxLifeTime), func() { manager.Provider.SessionGC(manager.maxLifeTime) })
}

var provides = make(map[string]Provider)

func NewManager(provideName, cookieName string, maxLifeTime int64) (*Manager, error) {

	provider, ok := provides[provideName]
	log.Println(ok)
	//if !ok {
	//	return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
	//}
	return &Manager{
		cookieName: cookieName,
		//lock:        sync.Mutex{},
		Provider:    provider,
		maxLifeTime: maxLifeTime,
	}, nil
}

type Session interface {
	Set(key, value interface{}) error // set session value
	Get(key interface{}) interface{}  // get session value
	Delete(key interface{}) error     // delete session value
	SessionID() string                // back current sessionID
}

type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}

func Register(name string, provider Provider) {
	if provider == nil {
		panic("session: Register provider is nil")
	}
	if _, dup := provides[name]; dup {
		panic("session: Register called twice for provider " + name)
	}
	provides[name] = provider
}

func (manager *Manager) sessionId() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {

	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		sessionId := manager.sessionId()
		manager.Provider.SessionInit(sessionId) //session is not use.
		cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sessionId), Path: "/", HttpOnly: true, MaxAge: int(manager.maxLifeTime)}
		http.SetCookie(w, &cookie)
	} else {
		sessionid, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.Provider.SessionRead(sessionid)
	}
	return
}

func Login(w http.ResponseWriter, r *http.Request) {

}

func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.cookieName)
	ok := util.CheckErr(err)
	if ok || cookie.Value == "" {
		return
	} else {
		manager.lock.Lock()
		defer manager.lock.Unlock()
		manager.Provider.SessionDestroy(cookie.Value)
		expiration := time.Now()
		cookie := http.Cookie{Name: manager.cookieName, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
		http.SetCookie(w, &cookie)
	}

}
