package gen_lobby

import (
	"errors"
	"strings"

	"github.com/itfantasy/gonode-toolkit/toolkit"
	"github.com/itfantasy/gonode-toolkit/toolkit/gamedb"
	"github.com/itfantasy/gonode/behaviors/gen_server"
)

// --------------------- super admin

func InitGameDB(mongoConf string) error {
	if err := gamedb.InitMongo(mongoConf); err != nil {
		return err
	}
	oldlobby, _ := lobbyManager().FindLobby(toolkit.DEFAULT_LOBBY)
	if oldlobby == nil {
		if _, err := CreateLobby(toolkit.DEFAULT_LOBBY); err != nil {
			return err
		}
	}
	return nil
}

func CreateLobby(lobbyId string) (*LobbyEntity, error) {
	oldlobby, _ := lobbyManager().FindLobby(lobbyId)
	if oldlobby != nil {
		return nil, errors.New("there has been a lobby has the same lobbyid!" + lobbyId)
	}
	lobby, err := lobbyManager().CreateLobby(lobbyId)
	if err != nil {
		return nil, err
	}
	return lobby, nil
}

func DisposeLobby(lobbyId string) (*LobbyEntity, error) {
	lobby, err := lobbyManager().FindLobby(lobbyId)
	if err != nil {
		return nil, err
	}
	if err := lobbyManager().DisposeLobby(lobbyId); err != nil {
		return nil, err
	}
	return lobby, nil
}

// --------------------- guest usr for all

func LobbyStats(lobbyId string) (*LobbyEntity, error) {
	lobby, err := lobbyManager().FindLobby(lobbyId)
	if err != nil {
		return nil, err
	}
	return lobby, nil
}

func RoomList(lobbyId string, startIndex int, endIndex int) ([]*LiteRoomEntity, error) {
	lobby, err := lobbyManager().FindLobby(lobbyId)
	if err != nil {
		return nil, err
	}
	return lobby.Rooms(startIndex, endIndex)
}

// --------------------- guest usr

func peerCannotFind(peerId string) error {
	return errors.New("cannot find the peer:" + peerId)
}

func JoinLobby(peerId string, lobbyId string) error {
	l, exist := getLobbyPeer(peerId)
	if !exist {
		return peerCannotFind(peerId)
	}
	l.SetLobbyId(lobbyId)
	return nil
}

func LeaveLobby(peerId string, lobbyId string) error {
	l, exist := getLobbyPeer(peerId)
	if !exist {
		return peerCannotFind(peerId)
	}
	l.SetDefaultLobby()
	return nil
}

func CreateRoom(peerId string, roomId string) (*LiteRoomEntity, error) {
	l, exist := getLobbyPeer(peerId)
	if !exist {
		return nil, peerCannotFind(peerId)
	}
	if !strings.HasSuffix(roomId, l.LobbyId()) {
		roomId = roomId + "@" + l.LobbyId()
	}
	lobby, err := lobbyManager().FindLobby(l.LobbyId())
	if err != nil {
		return nil, err
	}
	room, err := lobby.CreateRoom(roomId)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func JoinRoom(peerId string, roomId string) (*LiteRoomEntity, error) {
	l, exist := getLobbyPeer(peerId)
	if !exist {
		return nil, peerCannotFind(peerId)
	}
	lobby, err := lobbyManager().FindLobby(l.LobbyId())
	if err != nil {
		return nil, err
	}
	room, err := lobby.FindRoom(roomId)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func JoinRandomRoom(peerId string) (*LiteRoomEntity, error) {
	l, exist := getLobbyPeer(peerId)
	if !exist {
		return nil, peerCannotFind(peerId)
	}
	lobby, err := lobbyManager().FindLobby(l.LobbyId())
	if err != nil {
		return nil, err
	}
	room, err := lobby.RandomRoom()
	if err != nil {
		return nil, err
	}
	return room, nil
}

func AddPeer(peer *LobbyPeer) error {
	return peerManager().AddPeer(peer)
}

func RemovePeer(peerId string) error {
	return peerManager().RemovePeer(peerId)
}

func GetPeer(peerId string) (*LobbyPeer, bool) {
	return getLobbyPeer(peerId)
}

func RoomPubDomain(info *gen_server.NodeInfo) (string, bool) {
	pub, exist := info.UsrDatas[toolkit.USRDATA_PUBDOMAIN]
	return pub, exist
}
