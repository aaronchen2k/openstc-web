package _utils

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	_fileUtils "github.com/aaronchen2k/tester/internal/pkg/libs/file"
	"gopkg.in/ini.v1"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// 本机 ip
func LocalIP() string {
	ip := ""
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && !ipnet.IP.IsMulticast() && !ipnet.IP.IsLinkLocalUnicast() && !ipnet.IP.IsLinkLocalMulticast() && ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}
		}
	}
	return ip
}

// md5
func MD5(str string) string {
	encoder := md5.New()
	encoder.Write([]byte(str))
	return hex.EncodeToString(encoder.Sum(nil))
}

// 当前目录
func CWD() string {
	// 兼容 travis 集成测试
	if os.Getenv("TRAVIS_BUILD_DIR") != "" {
		return os.Getenv("TRAVIS_BUILD_DIR")
	}

	path, err := os.Executable()
	if err != nil {
		return ""
	}

	dir := filepath.Dir(path)
	return dir
}

func WwwPath() string {
	return "./www/dist"
}

// 工作目录
var workInDirLock sync.Mutex

func WorkInDir(f func(), dir string) {
	wd, _ := os.Getwd()
	workInDirLock.Lock()
	defer workInDirLock.Unlock()
	os.Chdir(dir)
	defer os.Chdir(wd)
	f()
}

func EXEName() string {
	path, err := os.Executable()
	if err != nil {
		return ""
	}
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

func LogDir() string {
	dir := filepath.Join(CWD(), "logs")
	_fileUtils.MkDirIfNeeded(dir)
	return dir
}

func ErrorLogFilename() string {
	return filepath.Join(LogDir(), fmt.Sprintf("%s-error.log", strings.ToLower(EXEName())))
}

func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

func IsPortInUse(port int) bool {
	if conn, err := net.DialTimeout("tcp", net.JoinHostPort("", fmt.Sprintf("%d", port)), 3*time.Second); err == nil {
		conn.Close()
		return true
	}
	return false
}

func init() {
	gob.Register(map[string]interface{}{})
	ini.PrettyFormat = false
}
