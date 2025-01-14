package pluto

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

/*
TODO
  Support channel expiration

const (
	MaxChannelLife = time.Hour * 24
)
*/

var (
	Channels      = map[uuid.UUID]Channel{}
	ChannelsMutex = new(sync.RWMutex)
)

type Joinable interface {
	Identifier
	Processor
}

type Channel struct {
	ID uuid.UUID `json:"id"`
	//OwnerID  uuid.NullUUID `json:"owner_id"`
	Name                        string     `json:"name"`
	Members                     []Joinable `json:"members"`
	Capacity                    int        `json:"capacity"`
	Length                      int        `json:"length"`
	DoNotAcceptDuplicateMembers bool       `json:"do_not_accept_duplicate_members"`

	OnJoin        Processor `json:"-"`
	OnLeave       Processor `json:"-"`
	OnMaxCapacity Processor `json:"-"`
	OnExpire      Processor `json:"-"`
}

func NewChannel(name string, length int) Channel {
	return Channel{
		ID:            uuid.New(),
		Name:          name,
		Members:       []Joinable{},
		Capacity:      length,
		Length:        length,
		OnJoin:        EmptyProcessor{},
		OnLeave:       EmptyProcessor{},
		OnMaxCapacity: EmptyProcessor{},
		OnExpire:      EmptyProcessor{},
	}
}

func (c *Channel) Publish(processable Processable) {
	for _, member := range c.Members {
		go member.Process(processable)
	}
}

func (c *Channel) Join(r Joinable) {
	if c.IsMember(r) && c.DoNotAcceptDuplicateMembers {
		return
	}

	if c.Capacity <= 0 {
		return
	}

	c.Members = append(c.Members, r)
	c.Capacity -= 1

	go c.OnJoin.Process(&InternalProcessable{
		ID: uuid.New(),
		Body: map[string]any{
			"channel":    c,
			"identifier": r,
		},
		CreatedAt: time.Now(),
	})

	if c.Capacity <= 0 {
		go c.OnMaxCapacity.Process(&InternalProcessable{
			ID: uuid.New(),
			Body: map[string]any{
				"channel": c,
			},
			CreatedAt: time.Now(),
		})

		// Do not call OnMaxCapacity several time
		c.OnMaxCapacity = EmptyProcessor{}
	}
}

func (c *Channel) Leave(identifier Identifier) {
	_, index := c.GetMember(identifier)
	if index == -1 {
		return
	}

	c.Members = append(c.Members[:index], c.Members[index+1:]...)
	c.Capacity += 1

	go c.OnLeave.Process(&InternalProcessable{
		ID: uuid.New(),
		Body: map[string]any{
			"channel":    c,
			"identifier": identifier,
		},
		CreatedAt: time.Now(),
	})
}

func (c *Channel) IsMember(identifier Identifier) bool {
	for _, member := range c.Members {
		if CompareIdentifiers(identifier, member) {
			return true
		}
	}
	return false
}

func (c *Channel) GetMember(identifier Identifier) (Joinable, int) {
	for i, member := range c.Members {
		if CompareIdentifiers(identifier, member) {
			return c.Members[i], i
		}
	}
	return nil, -1
}

func (c *Channel) UniqueProperty() string {
	return c.ID.String()
}

func (c *Channel) PredefinedKind() string {
	return KindChannel
}

type LockChannels struct {
	// TODO: Add support for read locks
	Write bool
}

func (p LockChannels) Process(processable Processable) (Processable, bool) {
	ChannelsMutex.Lock()
	return processable, true
}

type UnLockChannels struct {
}

func (p UnLockChannels) Process(processable Processable) (Processable, bool) {
	ChannelsMutex.Unlock()
	return processable, true
}

func getChannel(id uuid.UUID) (Channel, bool) {
	for _, channel := range Channels {
		if channel.ID == id {
			return channel, true
		}
	}
	return Channel{}, false
}
