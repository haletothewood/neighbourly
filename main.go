package main

import (
	"net/http"

	app "github.com/ainsleybc/neighbourly/app"
	r "github.com/dancannon/gorethink"
)

func main() {

	session, _ := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "neighbourly",
	})

	router := app.NewRouter(session)

	router.RegisterHandler("user signup", app.SignUpUser)
	router.RegisterHandler("user login", app.LoginUser)

	router.RegisterHandler("post add", app.AddPost)
	router.RegisterHandler("post subscribe", app.SubscribeFeed)
	router.RegisterHandler("post unsubscribe", app.UnsubscribePosts)

	router.RegisterHandler("feed add", app.AddFeed)
	router.RegisterHandler("feed subscribe", app.SubscribeFeed)
	router.RegisterHandler("feed unsubscribe", app.UnsubscribeFeed)

	http.Handle("/", router)
	http.ListenAndServe(":4000", nil)
}
