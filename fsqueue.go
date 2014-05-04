package fsqueue

import (
    "os"
    "os/user"
    "fmt"
    "time"
    "strconv"
    "io/ioutil"
    "sync"
)

type Channel struct {
    name     string
    fullpath string
    sync.RWMutex
}

func (c *Channel) MakeDirs() error {
    current := c.fullpath + "/current"
    failed  := c.fullpath + "/failed"
    success := c.fullpath + "/success"
    deleted := c.fullpath + "/deleted"

    err := os.MkdirAll(current, 0744)
    if err != nil { return err }

    err = os.MkdirAll(failed, 0744)
    if err != nil { return err }

    err = os.MkdirAll(success, 0744)
    if err != nil { return err }

    err = os.MkdirAll(deleted, 0744)
    return err
}

func (c *Channel) Remove() error {
    return os.RemoveAll(c.fullpath)
}

func (c *Channel) MakePayloadFilename() string {
    return fmt.Sprintf("%s/current/-%s", c.fullpath, strconv.FormatInt(time.Now().UnixNano(), 10))
}

func (c *Channel) Push(payloadBytes []byte) (string, error) {
    filename := c.MakePayloadFilename()

    // c.Lock()
    // defer c.Unlock()

    err := ioutil.WriteFile(filename, payloadBytes, 0744)
    return filename, err
}

func (c *Channel) Pop() ([]byte, error) {
    files, _ := ioutil.ReadDir(c.fullpath + "/current")
    file     := files[0]
    current  := c.fullpath + "/current"
    deleted  := c.fullpath + "/deleted"

    os.Rename(current + "/" + file.Name(), deleted + "/" + file.Name())

    payloadBytes, err := ioutil.ReadFile(deleted + "/" + file.Name())
    return payloadBytes, err
}

func DefaultPath() (string, error) {
    currentUser, err := user.Current()

    return currentUser.HomeDir + "/fsqueue-data", err
}

func Fullpath(chanName string) (string, error) {
    defaultPath, err := DefaultPath()
    return defaultPath + "/" + chanName, err
}

func Get(chanName string) (*Channel, error) {
    fullpath, err := Fullpath(chanName)
    return &Channel{name: chanName, fullpath: fullpath}, err
}

func MakeChannel(chanName string) (*Channel, error) {
    channel, err := Get(chanName)

    if err != nil {
        return channel, err
    }

    err = channel.MakeDirs()
    return channel, err
}

func RemoveChannel(chanName string) error {
    channel, err := Get(chanName)

    if err != nil {
        return err
    }

    return channel.Remove()
}
