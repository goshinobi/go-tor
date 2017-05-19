package tor

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	HOME_DIR      = os.Getenv("HOME")
	TOR_MULTI_DIR = HOME_DIR + "/.tor_multi"
)

func init() {
	if !isExist(TOR_MULTI_DIR) {
		err := os.Mkdir(TOR_MULTI_DIR, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func String(s string) *string {
	return &s
}

type Tor struct {
	ID            *int64  `json:"id"`
	SocksPort     *int    `json:"socks_port"`
	ControlPort   *int    `json:"control_port"`
	DataDirectory *string `json:"data_directory"`
	ConfPath      *string `json:"conf_path"`
	cmd           *exec.Cmd
	work          bool
}

func New(c ...string) *Tor {
	var (
		dataDir  string
		libDir   string
		confPath string
	)
	now := time.Now().UnixNano()
	socketPort, err := getPort()
	if err != nil {
		log.Println(err)
		return nil
	}
	controlPort, err := getPort()
	if err != nil {
		log.Println(err)
		return nil
	}

	if len(c) == 0 {
		libDir = TOR_MULTI_DIR + "/lib"
		confPath = fmt.Sprintf("%s/torrc.%d", TOR_MULTI_DIR, now)
	} else {
		libDir = fmt.Sprintf("%s/%s/lib", TOR_MULTI_DIR, c[0])
		confPath = fmt.Sprintf("%s/%s/torrc.%d", TOR_MULTI_DIR, c[0], now)
		if !isExist(fmt.Sprintf("%s/%s", TOR_MULTI_DIR, c[0])) {
			err := os.Mkdir(fmt.Sprintf("%s/%s", TOR_MULTI_DIR, c[0]), 0755)
			if err != nil {
				log.Println(err)
				return nil
			}
		}
	}
	dataDir = fmt.Sprintf("%s/tor%d", libDir, now)
	if !isExist(libDir) {
		err := os.Mkdir(libDir, 0755)
		if err != nil {
			log.Println(err)
			return nil
		}
	}
	tor := &Tor{
		ID:            &now,
		SocksPort:     socketPort,
		ControlPort:   controlPort,
		DataDirectory: &dataDir,
		ConfPath:      &confPath,
	}

	conf := fmt.Sprintf("SocksPort %d\n", *tor.SocksPort)
	conf += fmt.Sprintf("ControlPort %d\n", *tor.ControlPort)
	conf += fmt.Sprintf("DataDirectory %s\n", *tor.DataDirectory)

	if err := ioutil.WriteFile(fmt.Sprintf("%s", confPath), []byte(conf), 0644); err != nil {
		log.Println(err)
		return nil
	}
	return tor
}

func (t *Tor) Start() error {
	cmd := exec.Command("tor", "-f", *t.ConfPath)
	if err := cmd.Start(); err != nil {
		return err
	}
	t.cmd = cmd
	t.work = true
	return nil
}

func (t *Tor) Stop() error {
	if err := t.cmd.Process.Kill(); err != nil {
		return err
	}
	t.cmd = nil
	t.work = false
	return nil
}

func (t *Tor) Kill() error {
	if err := t.Stop(); err != nil {
		return err
	}

	if err := os.RemoveAll(*t.DataDirectory); err != nil {
		return err
	}
	if err := os.RemoveAll(*t.ConfPath); err != nil {
		return err
	}
	t.ID = nil
	t.SocksPort = nil
	t.ControlPort = nil
	t.DataDirectory = nil
	t.ConfPath = nil
	t.cmd = nil
	return nil
}

func (t *Tor) Reload() error {
	if err := t.Stop(); err != nil {
		return err
	}
	return t.Start()
}

func (t *Tor) Dial(timeout time.Duration, r string) (con net.Conn, err error) {
	addr := fmt.Sprintf("127.0.0.1:%s", *t.SocksPort)
	if oconn, err := net.DialTimeout("tcp", addr, timeout); err == nil {
		oconn.Write([]byte{ // VERSION_AUTH
			5, // PROTO_VER5
			1, //
			0, // NO_AUTH
		})
		buffer := [64]byte{}
		oconn.Read(buffer[:])
		buffer[0] = 5 // VER  5
		buffer[1] = 1 // CMD connect
		buffer[2] = 0 // RSV
		buffer[3] = 3 // DOMAINNAME: X'03'

		host, port := splitHostAndPort(r)

		hostBytes := []byte(host)
		buffer[4] = byte(len(hostBytes))
		copy(buffer[5:], hostBytes)
		binary.BigEndian.PutUint16(buffer[5+len(hostBytes):], uint16(port))
		oconn.Write(buffer[:5+len(hostBytes)+2])

		if n, err := oconn.Read(buffer[:]); n > 1 && err == nil && buffer[1] == 0 {
			return oconn, nil
		} else {
			return nil, fmt.Errorf("connet to socks server %s error: %v", addr, err)
		}
	} else {
		return nil, err
	}
}

func (t Tor) String() string {
	bin, err := json.MarshalIndent(t, "", "    ")
	if err != nil {
		return ""
	}
	return string(bin)
}

func getPort() (*int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return nil, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}
	defer l.Close()
	return &l.Addr().(*net.TCPAddr).Port, nil
}

func isExist(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func splitHostAndPort(host string) (string, uint16) {
	if idx := strings.Index(host, ":"); idx < 0 {
		return host, 80
	} else {
		port, _ := strconv.Atoi(host[idx+1:])
		return host[:idx], uint16(port)
	}
}
