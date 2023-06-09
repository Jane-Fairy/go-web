package session

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"net/url"
	"sync"
)

type Manager struct {
	cookieName  string
	lock        sync.Mutex
	provider    Provider
	maxLifeTime int64
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
		provider:    provider,
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
		manager.provider.SessionInit(sessionId) //session is not use.
		cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sessionId), Path: "/", HttpOnly: true, MaxAge: int(manager.maxLifeTime)}
		http.SetCookie(w, &cookie)
	} else {
		sessionid, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.provider.SessionRead(sessionid)
	}
	return
}

func Login(w http.ResponseWriter, r *http.Request) {

}
