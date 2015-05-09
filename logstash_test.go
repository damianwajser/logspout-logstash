package logstash

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/gliderlabs/logspout/router"
    "net"
    "github.com/fsouza/go-dockerclient"
    "time"
    "encoding/json"
)

var res string

type MockConn struct {
}

func (m MockConn) Close() error {
    return nil
}

func (m MockConn) Read(b []byte) (n int, err error) {
    return 0, nil
}

func (m MockConn) Write(b []byte) (n int, err error) {
    res = string(b)
    return 0, nil
}

func (m MockConn) LocalAddr() net.Addr {
    return nil
}

func (m MockConn) RemoteAddr() net.Addr {
    return nil
}

func (m MockConn) SetDeadline(t time.Time) error {
    return nil
}

func (m MockConn) SetReadDeadline(t time.Time) error {
    return nil
}

func (m MockConn) SetWriteDeadline(t time.Time) error {
    return nil
}


func TestStreamNotJson(t *testing.T) {
    assert := assert.New(t)

    conn := MockConn{}

    adapter := LogstashAdapter{
        route: new(router.Route),
        conn:  conn,
    }

    assert.NotNil(adapter)

    logstream := make(chan *router.Message)

    containerConfig := docker.Config{}
    containerConfig.Image = "image"
    containerConfig.Hostname = "hostname"

    container := docker.Container{}
    container.Name = "name"
    container.ID = "ID"
    container.Config = &containerConfig

    str := `foo bananas`

    message := router.Message{
        Container: &container,
        Source: "FOOOOO",
        Data: str,
        Time: time.Now(),
    }

    go func() {
        logstream <- &message
        close(logstream)
    }()

    adapter.Stream(logstream)

    var data map[string]interface{}
    err := json.Unmarshal([]byte(res), &data)
    assert.Nil(err)

    assert.Equal("name", data["docker.name"])
    assert.Equal("ID", data["docker.id"])
    assert.Equal("image", data["docker.image"])
    assert.Equal("hostname", data["docker.hostname"])
    assert.Equal("foo bananas", data["message"])
}

func TestStreamJson(t *testing.T) {
    assert := assert.New(t)

    conn := MockConn{}

    adapter := LogstashAdapter{
        route: new(router.Route),
        conn:  conn,
    }

    assert.NotNil(adapter)

    logstream := make(chan *router.Message)

    containerConfig := docker.Config{}
    containerConfig.Image = "image"
    containerConfig.Hostname = "hostname"

    container := docker.Container{}
    container.Name = "name"
    container.ID = "ID"
    container.Config = &containerConfig

    str := `{ "remote_user": "-", "body_bytes_sent": "25", "request_time": "0.821", "status": "200", "request_method": "POST", "http_referrer": "-", "http_user_agent": "-" }`

    message := router.Message{
        Container: &container,
        Source: "FOOOOO",
        Data: str,
        Time: time.Now(),
    }

    go func() {
        logstream <- &message
        close(logstream)
    }()

    adapter.Stream(logstream)

    var data map[string]interface{}
    err := json.Unmarshal([]byte(res), &data)
    assert.Nil(err)

    assert.Equal("name", data["docker.name"])
    assert.Equal("ID", data["docker.id"])
    assert.Equal("image", data["docker.image"])
    assert.Equal("hostname", data["docker.hostname"])
    assert.Equal("-", data["remote_user"])
    assert.Equal("25", data["body_bytes_sent"])
    assert.Equal("0.821", data["request_time"])
    assert.Equal("200", data["status"])
    assert.Equal("POST", data["request_method"])
    assert.Equal("-", data["http_referrer"])
    assert.Equal("-", data["http_user_agent"])
}

