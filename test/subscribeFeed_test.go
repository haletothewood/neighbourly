package app

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	. "github.com/ainsleybc/neighbourly/app"
	r "github.com/dancannon/gorethink"
	"github.com/posener/wstest"
)

// "encoding/json"
// "fmt"
// "testing"
// "time"

// r "github.com/dancannon/gorethink"
// "github.com/posener/wstest"

func TestSubscribeFeed(t *testing.T) {

	// // connect to rethinkDB
	session, _ := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "test",
	})

	// close session on end test
	defer session.Close()

	// create the tables for test
	r.TableCreate("feeds").RunWrite(session)

	// new router
	testRouter := NewRouter(session)

	// mock server thingy
	d := wstest.NewDialer(testRouter, nil)

	// open websocket connection (skip error)
	conn, resp, _ := d.Dial("ws://localhost:4000", nil)

	got, want := resp.StatusCode, http.StatusSwitchingProtocols
	if got != want {
		t.Errorf("resp.StatusCode: %q, want: %q", got, want)
	}

	// register handler for addFeed message
	testRouter.RegisterHandler("feed subscribe", SubscribeFeed)

	// creating test message and passing it through websocket
	rawMessage := json.RawMessage(`{"name":"feed subscribe"}`)
	conn.WriteJSON(rawMessage)

	// create feed & add to database
	feed := &Feed{
		Address: "Makers Academy",
	}
	r.Table("feeds").Insert(feed).RunWrite(session)

	// simple timeout to allow to database writes
	time.Sleep(time.Second * 1)

	// readJSON from socket
	var output Message
	conn.ReadJSON(&output)

	// write assertion
	got2, want2 := output.Name, "feed add"
	if got2 != want2 {
		t.Errorf("got: %v, want: %v", got2, want2)
	}

	r.TableDrop("feeds").Wait().Exec(session)
}
