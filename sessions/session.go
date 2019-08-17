package session

// https://astaxie.gitbooks.io/build-web-application-with-golang/en/06.2.html

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Manager struct {
	cookieName  string
	lock        sync.Mutex
	provider    Provider
	maxlifetime int64
}

func NewManager(provideName, cookieName string, maxlifetime int64) (*Manager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (check imports)", provideName)
	}
	return &Manager{provider: provider, cookieName: cookieName, maxlifetime: maxlifetime}, nil
}

func (manager *Manager) sessionID() string {
	b := make([]byte, 32)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		// no session, create new one
		sid := manager.sessionID()
		session, _ := manager.provider.SessionInit(sid)
		cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxlifetime)}
		http.SetCookie(w, &cookie)
		return session
	} else {
		// has session, return it
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ := manager.provider.SessionRead(sid)
		return session
	}
}

func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		manager.lock.Lock()
		defer manager.lock.Unlock()
		manager.provider.SessionDestroy(cookie.Value)
		expiration := time.Now()
		cookie := http.Cookie{Name: manager.cookieName, Expires: expiration, Path: "/", HttpOnly: true, MaxAge: -1}
		http.SetCookie(w, &cookie)
	}
}

func (manager *Manager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxlifetime)
	time.AfterFunc(time.Duration(manager.maxlifetime), func() { manager.GC() })
}

type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionRead(sid string) (Session, error)
	SessionGC(maxLifeTime int64)
}

var provides = make(map[string]Provider)

func Register(name string, provider Provider) {
	if provider == nil {
		panic("session: Register provider is nil")
	}
	if _, dup := provides[name]; dup {
		panic("session: Register called twice for provider " + name)
	}

	provides[name] = provider
}

type Session interface {
	Set(key, data interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionID() string
}
