package main

import (
	"fmt"

	"github.com/Philipp15b/go-steam"
	"github.com/Philipp15b/go-steam/socialcache"
	"github.com/Philipp15b/go-steam/steamid"
)

var userID steamid.SteamId
var user steam.GroupStateEvent

var group = socialcache.Group{
	SteamId: 3412828,
	Name:    "CodeGifts Sorteio",
}

func main() {
	userID = 76561198057002771
	user.SteamId = userID
	fmt.Println(user.IsMember())

	gl := socialcache.NewGroupsList()
	gl.Add(group)

	fmt.Println(gl.Count())
}
