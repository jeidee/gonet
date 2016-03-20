package gonet

import (
	"errors"
	"net"
)

// SessionGroupManager has a tree structure.
// It means that SessionGroupManager has a SessionGroupManager as child,
// and also has a map of sessions to manage its own sessions.
type SessionGroupManager struct {
	sessionGroups map[string]*SessionGroupManager
	sessions      map[net.Conn]*Session
	parent        *SessionGroupManager
	name          string
}

func NewSessionGroupManager(groupName string, parent *SessionGroupManager) *SessionGroupManager {
	sgm := new(SessionGroupManager)

	sgm.sessions = make(map[net.Conn]*Session)
	sgm.sessionGroups = make(map[string]*SessionGroupManager)
	sgm.parent = parent
	sgm.name = groupName

	return sgm
}

func (self *SessionGroupManager) NewGroup(groupNames ...string) (*SessionGroupManager, error) {
	if self.FindGroup(self, groupNames...) != nil {
		return nil, errors.New("Already exists.")
	}

	groups := self.sessionGroups
	parent := self

	for _, groupName := range groupNames {

		if groupName == "root" {
			return nil, errors.New("<root> is reserved.")
		}

		if groups[groupName] == nil {
			newGroup := NewSessionGroupManager(groupName, parent)
			groups[groupName] = newGroup
			parent = groups[groupName]
			groups = parent.sessionGroups
		}
	}

	return parent, nil
}

func (self *SessionGroupManager) RemoveGroup(groupNames ...string) error {
	group := self.FindGroup(self, groupNames...)
	if group == nil {
		return errors.New("Not exists group.")
	}

	if group.parent == nil {
		return errors.New("Invalid a parent of this group.")
	}

	delete(group.parent.sessionGroups, group.name)

	return nil
}

func (self *SessionGroupManager) FindGroup(sessionGroupManager *SessionGroupManager, groupNames ...string) *SessionGroupManager {
	if sessionGroupManager == nil {
		return nil
	}

	if len(groupNames) == 1 {
		return sessionGroupManager.sessionGroups[groupNames[0]]
	}

	nextFindSgm := sessionGroupManager
	for _, groupName := range groupNames {
		nextFindSgm := nextFindSgm.FindGroup(nextFindSgm, groupName)
		if nextFindSgm == nil {
			return nil
		}
	}

	return nextFindSgm
}

func (self *SessionGroupManager) JoinGroup(groupName string, session *Session) error {
	if self.sessionGroups[groupName] == nil {
		return errors.New("Not exists group.")
	}

	return self.sessionGroups[groupName].Join(session)
}

func (self *SessionGroupManager) LeaveGroup(groupName string, session *Session) error {
	if self.sessionGroups[groupName] == nil {
		return errors.New("Not exists group.")
	}

	return self.sessionGroups[groupName].Leave(session)
}

func (self *SessionGroupManager) Join(session *Session) error {
	if self.sessions[session.conn] != nil {
		return errors.New("Already joined.")
	}
	self.sessions[session.conn] = session
	return nil
}

func (self *SessionGroupManager) Leave(session *Session) error {
	if self.sessions[session.conn] == nil {
		return errors.New("Not exists session.")
	}
	delete(self.sessions, session.conn)
	return nil
}

func (self *SessionGroupManager) Broadcast(data interface{}, exceptSessions ...*Session) {
	for _, session := range self.sessions {
		isExcepted := false
		for _, exceptSession := range exceptSessions {
			if exceptSession == session {
				isExcepted = true
				break
			}
		}
		if isExcepted {
			continue
		}

		session.Send(data)
	}

	for _, sg := range self.sessionGroups {
		sg.Broadcast(data)
	}
}
