/**
 * Copyright 2024 Boris Wolf
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
