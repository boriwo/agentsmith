package main

type SimpleSessionManager struct {
	sessions map[string]*UserSession
}

type SessionManager interface {
	GetSession(user *User) *UserSession
	DeleteSession(id string)
}

func NewSimpleSessionManager() SessionManager {
	mgr := new(SimpleSessionManager)
	mgr.sessions = make(map[string]*UserSession, 0)
	return mgr
}

func (mgr *SimpleSessionManager) GetSession(user *User) *UserSession {
	if user == nil || user.Id == "" {
		return nil
	}
	session, ok := mgr.sessions[user.Id]
	if !ok {
		session = new(UserSession)
		session.User = user
		session.State = STATE_QA
		mgr.sessions[user.Id] = session
	}
	return session
}

func (mgr *SimpleSessionManager) DeleteSession(id string) {
	delete(mgr.sessions, id)
}
