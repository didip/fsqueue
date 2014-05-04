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

    channel, _       := MakeChannel(chanName)
    _, filename, err := channel.Push(data)

    if err != nil {
        t.Fatal(err)
    }

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

    payload, err := channel.Pop()

    if err != nil { t.Fatal(err) }

    if bytes.Compare(payload, data) != 0 {
        t.Fatal("Payload doesn't match data")
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