package fsqueue

import (
    "os"
    "os/user"
    "time"
    "strconv"
    "io/ioutil"
)

func DefaultPath() (string, error) {
    currentUser, err := user.Current()

    return currentUser.HomeDir + "/fsqueue-data", err
}

func FullPath(chanName string) (string, error) {
    defaultPath, err := DefaultPath()

    return defaultPath + "/" + chanName, err
}

func NewChannel(chanName string) (*Channel, error) {
    fullPath, err := FullPath(chanName)
    channel       := &Channel{}

    channel.name        = chanName
    channel.fullPath    = fullPath
    channel.currentPath = fullPath + "/current"
    channel.failedPath  = fullPath + "/failed"
    channel.successPath = fullPath + "/success"
    channel.deletedPath = fullPath + "/deleted"

    return channel, err
}

func MakeChannel(chanName string) (*Channel, error) {
    channel, err := NewChannel(chanName)

    if err == nil {
        err = channel.MakeDirs()
    }

    return channel, err
}

func RemoveChannel(chanName string) error {
    channel, err := NewChannel(chanName)

    if err != nil { return err }

    return channel.Remove()
}

type Channel struct {
    name        string
    fullPath    string
    currentPath string
    failedPath  string
    successPath string
    deletedPath string
}

func (c *Channel) MakeDirs() error {
    err := os.MkdirAll(c.currentPath, 0744)
    if err != nil { return err }

    err = os.MkdirAll(c.failedPath, 0744)
    if err != nil { return err }

    err = os.MkdirAll(c.successPath, 0744)
    if err != nil { return err }

    err = os.MkdirAll(c.deletedPath, 0744)
    return err
}

func (c *Channel) Remove() error {
    return os.RemoveAll(c.fullPath)
}

func (c *Channel) PayloadId() (string, string) {
    id           := strconv.FormatInt(time.Now().UnixNano(), 10)
    fullFilePath := c.currentPath + "/" + id
    return id, fullFilePath
}

func (c *Channel) Push(payloadBytes []byte) (string, string, error) {
    id, fullFilePath := c.PayloadId()
    err := ioutil.WriteFile(fullFilePath, payloadBytes, 0744)

    return id, fullFilePath, err
}

func (c *Channel) Pop() ([]byte, error) {
    files, _ := ioutil.ReadDir(c.currentPath)

    if len(files) <= 0 {
        return nil, nil
    }

    file    := files[0]
    current := c.currentPath + "/" + file.Name()
    deleted := c.deletedPath + "/" + file.Name()

    os.Rename(current, deleted)

    payloadBytes, err := ioutil.ReadFile(deleted)

    go os.Remove(deleted)

    return payloadBytes, err
}

