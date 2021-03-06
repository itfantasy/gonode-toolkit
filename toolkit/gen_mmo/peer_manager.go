package gen_mmo

import (
	"github.com/itfantasy/gonode/core/actors"
)

type MmoPeer struct {
	peerId            string
	requestExecutor   *actors.Executor
	actor             *MmoActor
	radarSubscription IDisposer
	eventer           IMmoEventer
}

func NewMmoPeer(peerId string, eventer IMmoEventer) *MmoPeer {
	m := new(MmoPeer)
	m.peerId = peerId
	m.requestExecutor = actors.Spawn(1024)
	m.eventer = eventer
	return m
}

func (m *MmoPeer) PeerId() string {
	return m.peerId
}

func (m *MmoPeer) AvatarId() string {
	if m.actor != nil && m.actor.Avatar() != nil {
		return m.actor.Avatar().Id()
	}
	return ""
}

func (m *MmoPeer) RequestExecutor() *actors.Executor {
	return m.requestExecutor
}

func (m *MmoPeer) MmoActor() *MmoActor {
	return m.actor
}

func (m *MmoPeer) MmoEventer() IMmoEventer {
	return m.eventer
}

func (m *MmoPeer) SetActorHandler(actor *MmoActor) {
	m.actor = actor
}

func (m *MmoPeer) DisposeActor() {
	m.actor.Dispose()
	m.actor = nil
}

func (m *MmoPeer) RadarSubscription() IDisposer {
	return m.radarSubscription
}

func (m *MmoPeer) SetRadarSubscription(r IDisposer) {
	m.radarSubscription = r
}

func (m *MmoPeer) DisposeRadarSubscription() {
	m.radarSubscription = nil
}

func getMmoPeer(peerId string) (*MmoPeer, bool) {
	peer, ok := peerManager().GetPeer(peerId)
	if !ok {
		return nil, false
	}
	cntpeer, ok := peer.(*MmoPeer)
	if !ok {
		return nil, false
	}
	return cntpeer, true
}
