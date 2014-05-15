package fsqueue

import (
    "os"
    "bytes"
    "testing"
)

func TestMakeThenRemove(t *testing.T) {
    chanName := "fsqueue-test"

    if _, err := MakeChannel(chanName); err != nil {
        t.Fatal(err)
    }

    if err := RemoveChannel(chanName); err != nil {
        t.Fatal(err)
    }
}

func TestPush(t *testing.T) {
    chanName := "fsqueue-test"
    data     := []byte("{\"method\": \"Run\"}")

    channel, _  := MakeChannel(chanName)
    _, filename := channel.Push(data)
    <- channel.donePushChan

    if _, err := os.Stat(filename); os.IsNotExist(err) {
        t.Fatal(err)
    }

    RemoveChannel(chanName)
}

func TestPop(t *testing.T) {
    chanName := "fsqueue-test"
    data     := []byte("{\"method\": \"Run\"}")

    channel, _  := MakeChannel(chanName)
    channel.Push(data)
    <- channel.donePushChan

    payload, err := channel.Pop()

    if err != nil { t.Fatal(err) }

    if bytes.Compare(payload, data) != 0 {
        t.Fatal("Payload doesn't match data")
    }

    RemoveChannel(chanName)
}

func TestCount(t *testing.T) {
    chanName := "fsqueue-test"
    data     := []byte("{\"method\": \"Run\"}")

    channel, _  := MakeChannel(chanName)
    channel.Push(data)
    <- channel.donePushChan

    if count, _ := channel.Count("current"); count != 1 {
        t.Fatal("Channel current bucket should contain 1 item")
    }

    channel.Pop()

    if count, _ := channel.Count("deleted"); count != 1 {
        t.Fatal("Channel deleted bucket should contain 1 item")
    }

    RemoveChannel(chanName)
}

func TestOldestNewest(t *testing.T) {
    chanName   := "fsqueue-test"
    channel, _ := MakeChannel(chanName)

    channel.Push([]byte("{\"method\": \"Run\"}"))
    <- channel.donePushChan

    if count, _ := channel.Count("current"); count != 1 {
        t.Fatal("Channel current bucket should contain 1 items")
    }

    channel.Push([]byte("{\"method\": \"Run2\"}"))
    <- channel.donePushChan

    if count, _ := channel.Count("current"); count != 2 {
        t.Fatal("Channel current bucket should contain 2 items")
    }

    data, _ := channel.Oldest("current")

    if bytes.Compare([]byte("{\"method\": \"Run\"}"), data) != 0 {
        t.Fatal("Payload doesn't match data: ", data)
    }

    data, _ = channel.Newest("current")

    if bytes.Compare([]byte("{\"method\": \"Run2\"}"), data) != 0 {
        t.Fatal("Payload doesn't match data: ", data)
    }

    RemoveChannel(chanName)
}


func BenchmarkPushPop(b *testing.B) {
    b.StopTimer()

    chanName := "fsqueue-test"
    data     := []byte("{\"method\": \"Run\"}")

    channel, _ := MakeChannel(chanName)

    b.StartTimer()

    for n := 0; n < b.N; n++ {
        go channel.Push(data)
        go channel.Pop()
    }

    b.StopTimer()

    RemoveChannel(chanName)
}