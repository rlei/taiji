package taiji

import (
    "bufio"
    "net"
    "testing"
    "time"
)

func TestRunStop(t *testing.T) {
    server := NewServer(5555, 6666)
    go server.ListenAndServe()
    time.Sleep(500 * time.Millisecond)

    conn, err := net.Dial("tcp", ":5555")
    if err != nil {
        t.Error("Failed to connect to push server client port.")
    }
    defer conn.Close()
    conn1, err := net.Dial("tcp", ":6666")
    if err != nil {
        t.Error("Failed to connect to push server control port.")
    }
    defer conn1.Close()

    server.close()
    time.Sleep(500 * time.Millisecond)

    conn2, err := net.Dial("tcp", ":5555")
    if err == nil {
        t.Error("Expecting closed server but still connected successfully.")
        defer conn2.Close()
    }
    conn3, err := net.Dial("tcp", ":6666")
    if err == nil {
        t.Error("Expecting closed server but still connected to control port successfully.")
        defer conn3.Close()
    }
}

func TestBroadcast(t *testing.T) {
    server := NewServer(5555, 6666)
    go server.ListenAndServe()
    time.Sleep(500 * time.Millisecond)
    defer server.close()

    clientConn, err := net.Dial("tcp", ":5555")
    if err != nil {
        t.Errorf("Failed to connect to server %s", err.Error())
        return
    }
    defer clientConn.Close()

    clientReader := bufio.NewReader(clientConn)
    clientWriter := bufio.NewWriter(clientConn)

    clientWriter.WriteString("HELO abc\n")
    clientWriter.Flush()

    clientConn1, _ := net.Dial("tcp", ":5555")
    defer clientConn1.Close()

    clientReader1 := bufio.NewReader(clientConn1)
    clientWriter1 := bufio.NewWriter(clientConn1)

    clientWriter1.WriteString("HELO 123\n")
    clientWriter1.Flush()

    controlConn, _ := net.Dial("tcp", ":6666")
    defer controlConn.Close()

    time.Sleep(500 * time.Millisecond)
    controlWriter := bufio.NewWriter(controlConn)
    controlWriter.WriteString("b \"Hello\"\n")
    controlWriter.Flush()

    expected := "BROADCAST \"Hello\"\n"
    s, _ := clientReader.ReadString('\n')
    if s != expected {
        t.Errorf("Expected %s but received %s.", expected, s)
    }
    s, _ = clientReader1.ReadString('\n')
    if s != expected {
        t.Errorf("Expected %s but received %s.", expected, s)
    }
}

func TestSendTo(t *testing.T) {
    server := NewServer(5555, 6666)
    go server.ListenAndServe()
    time.Sleep(500 * time.Millisecond)
    defer server.close()

    clientConn, err := net.Dial("tcp", ":5555")
    if err != nil {
        t.Errorf("Failed to connect to server %s", err.Error())
        return
    }
    defer func() {
        clientConn.Close()
        time.Sleep(500 * time.Millisecond)
    } ()

    clientReader := bufio.NewReader(clientConn)
    clientWriter := bufio.NewWriter(clientConn)

    clientWriter.WriteString("HELO abc\n")
    clientWriter.Flush()

    clientConn1, _ := net.Dial("tcp", ":5555")
    defer clientConn1.Close()

    clientWriter1 := bufio.NewWriter(clientConn1)

    clientWriter1.WriteString("HELO 123\n")
    clientWriter1.Flush()

    controlConn, _ := net.Dial("tcp", ":6666")
    defer controlConn.Close()

    time.Sleep(500 * time.Millisecond)
    controlWriter := bufio.NewWriter(controlConn)
    controlWriter.WriteString("s abc \"Hello\"\n")
    controlWriter.Flush()

    expected := "BROADCAST \"Hello\"\n"
    s, _ := clientReader.ReadString('\n')
    if s != expected {
        t.Errorf("Expected %s but received %s.", expected, s)
    }
    clientConn1.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
    var buf [32]byte
    n, err := clientConn1.Read(buf[0:])
    if err == nil {
        t.Errorf("Expected nothing read for client 123 but read %d bytes", n)
    }
}

func TestClientHelloTimeout(t *testing.T) {
}

