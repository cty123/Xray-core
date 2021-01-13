package outbound

import (
	"context"
	"github.com/xtaci/smux"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/transport/internet"
	"sync"
)

type SmuxManager struct {
	access           sync.RWMutex
	muxConnectionMap map[net.Destination]net.Conn
}

func NewSmuxManager() SmuxManager {
	return SmuxManager{
		muxConnectionMap: make(map[net.Destination]net.Conn),
	}
}

func (sm *SmuxManager) addConnection(ctx context.Context, dest net.Destination, dialer internet.Dialer) (net.Conn, error) {
	conn, err := dialer.Dial(ctx, dest)
	if err != nil {
		return nil, err
	}

	session, err := smux.Client(conn, nil)
	if err != nil {
		return nil, err
	}

	muxConn, err := session.OpenStream()
	if err != nil {
		return nil, err
	}

	sm.muxConnectionMap[dest] = muxConn
	return muxConn, nil
}

func (sm *SmuxManager) getConnection(ctx context.Context, dest net.Destination, dialer internet.Dialer) (net.Conn, error) {
	sm.access.Lock()
	defer sm.access.Unlock()

	if conn, ok := sm.muxConnectionMap[dest]; ok {
		return conn, nil
	}

	conn, err := sm.addConnection(ctx, dest, dialer)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (sm *SmuxManager) removeConnection(dest net.Destination) {

}
