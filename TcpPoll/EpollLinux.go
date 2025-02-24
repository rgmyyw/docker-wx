// +build linux

package TcpPoll

import (
	"log"
	"net"
	"reflect"
	"sync"
	"syscall"

	"golang.org/x/sys/unix"
)

type epoll struct {
	fd          int					// poll的文件描述
	connections map[int]net.Conn	// 长连接列表, key: fd, value: conn
	lock        *sync.RWMutex		// 读写锁
}

// 获取连接池
func MkEpoll() (*epoll, error) {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &epoll{
		fd:          fd,
		lock:        &sync.RWMutex{},
		connections: make(map[int]net.Conn),
	}, nil
}

// 连接池增加长连接
func (e *epoll) Add(conn net.Conn) (int, error) {
	// Extract file descriptor associated with the connection
	fd := socketFD(conn)
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: unix.POLLIN | unix.POLLHUP | unix.EPOLLPRI, Fd: int32(fd)})
	if err != nil {
		return 0, err
	}
	e.lock.Lock()
	defer e.lock.Unlock()
	e.connections[fd] = conn
	if len(e.connections)%100 == 0 {
		log.Printf("total number of connections: %v", len(e.connections))
	}
	return fd, nil
}

// 连接池移除长连接
func (e *epoll) Remove(conn net.Conn) error {
	fd := socketFD(conn)
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		return err
	}
	e.lock.Lock()
	defer e.lock.Unlock()
	delete(e.connections, fd)
	if len(e.connections)%100 == 0 {
		log.Printf("total number of connections: %v", len(e.connections))
	}
	return nil
}

// 等待
func (e *epoll) Wait() ([]int, error) {
	events := make([]unix.EpollEvent, 100)
retry:
	n, err := unix.EpollWait(e.fd, events, 100)
	if err != nil {
		if err == unix.EINTR {
			goto retry
		}
		return nil, err
	}
	e.lock.RLock()
	defer e.lock.RUnlock()
	var fds []int
	for i := 0; i < n; i++ {
		fds = append(fds, int(events[i].Fd))
	}
	return fds, nil
}

// 获取长连接的文件描述fd
func socketFD(conn net.Conn) int {
	//tls := reflect.TypeOf(conn.UnderlyingConn()) == reflect.TypeOf(&tls.Conn{})
	// Extract the file descriptor associated with the connection
	//connVal := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn").Elem()
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	//if tls {
	//	tcpConn = reflect.Indirect(tcpConn.Elem())
	//}
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")

	return int(pfdVal.FieldByName("Sysfd").Int())
}
