package fsqueue

import (
    "os"
    "os/user"
    "time"
    "strconv"
    "io/ioutil"
    "reflect"
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

    channel.name         = chanName
    channel.fullPath     = fullPath
    channel.currentPath  = fullPath + "/current"
    channel.failedPath   = fullPath + "/failed"
    channel.successPath  = fullPath + "/success"
    channel.deletedPath  = fullPath + "/deleted"
    channel.donePushChan = make(chan bool)

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
    name         string
    fullPath     string
    currentPath  string
    failedPath   string
    successPath  string
    deletedPath  string
    donePushChan chan bool
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

func (c *Channel) Payloads(bucket string) ([]os.FileInfo, error) {
    bucketPath := reflect.Indirect(reflect.ValueOf(c)).FieldByName(bucket + "Path")

    return ioutil.ReadDir(bucketPath.String())
}

func (c *Channel) Count(bucket string) (int, error) {
    files, err := c.Payloads(bucket)

    if err != nil { return 0, err }

    return len(files), err
}

func (c *Channel) Oldest(bucket string) ([]byte, error) {
    files, err := c.Payloads(bucket)

    if err != nil { return nil, err }

    if len(files) <= 0 { return nil, nil }

    bucketPath := reflect.Indirect(reflect.ValueOf(c)).FieldByName(bucket + "Path")

    return ioutil.ReadFile(bucketPath.String() + "/" + files[0].Name())
}

func (c *Channel) Newest(bucket string) ([]byte, error) {
    files, err := c.Payloads(bucket)
    numFiles   := len(files)

    if err != nil { return nil, err }

    if numFiles <= 0 { return nil, nil }

    bucketPath := reflect.Indirect(reflect.ValueOf(c)).FieldByName(bucket + "Path")

    return ioutil.ReadFile(bucketPath.String() + "/" + files[numFiles - 1].Name())
}

func (c *Channel) PushBytes(payloadBytes []byte) (string, string) {
    id, fullFilePath := c.PayloadId()

    go func() {
        ioutil.WriteFile(fullFilePath, payloadBytes, 0744)
        c.donePushChan <- true
    }()

    return id, fullFilePath
}

func (c *Channel) PopBytes() ([]byte, error) {
    files, _ := c.Payloads("current")

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

