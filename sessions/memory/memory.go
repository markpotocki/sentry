package memory

// https://astaxie.gitbooks.io/build-web-application-with-golang/en/06.3.html

import (
	"container/list"
	session "idendity-provider/sessions"
	"sync"
	"time"
)

var pder = &MemoryProvider{list: list.New()}

type SessionStore struct {
	sid          string
	timeAccessed time.Time
	value        map[interface{}]interface{}
}

func (s *SessionStore) Set(key, data interface{}) error {
	s.value[key] = data
	pder.SessionUpdate(s.sid)
	return nil
}

func (s *SessionStore) Get(key interface{}) interface{} {
	pder.SessionUpdate(s.sid)
	if v, ok := s.value[key]; ok {
		return v
	} else {
		return nil
	}
	return nil
}

func (s *SessionStore) Delete(key interface{}) error {
	delete(s.value, key)
	pder.SessionUpdate(s.sid)
	return nil
}

func (s *SessionStore) SessionID() string {
	return s.sid
}

type MemoryProvider struct {
	lock     sync.Mutex
	sessions map[string]*list.Element
	list     *list.List
}

func (pder *MemoryProvider) SessionInit(sid string) (session.Session, error) {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	newsess := &SessionStore{sid: sid, timeAccessed: time.Now(), value: v}
	element := pder.list.PushBack(newsess)
	pder.sessions[sid] = element
	return newsess, nil
}

func (pder *MemoryProvider) SessionRead(sid string) (session.Session, error) {
	if element, ok := pder.sessions[sid]; ok {
		return element.Value.(*SessionStore), nil
	} else {
		sess, err := pder.SessionInit(sid)
		return sess, err
	}
	return nil, nil
}

func (pder *MemoryProvider) SessionDestroy(sid string) error {
	if element, ok := pder.sessions[sid]; ok {
		delete(pder.sessions, sid)
		pder.list.Remove(element)
		return nil
	}
	return nil
}

func (pder *MemoryProvider) SessionGC(maxlifetime int64) {
	pder.lock.Lock()
	defer pder.lock.Unlock()

	for {
		element := pder.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*SessionStore).timeAccessed.Unix() + maxlifetime) < time.Now().Unix() {
			pder.list.Remove(element)
			delete(pder.sessions, element.Value.(*SessionStore).sid)
		} else {
			break
		}
	}
}

func (pder *MemoryProvider) SessionUpdate(sid string) error {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	if element, ok := pder.sessions[sid]; ok {
		element.Value.(*SessionStore).timeAccessed = time.Now()
		pder.list.MoveToFront(element)
		return nil
	}
	return nil
}

func init() {
	pder.sessions = make(map[string]*list.Element, 0)
	session.Register("memory", pder)
}
