package main

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

type UserCache struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Token    string `json:"token"`
	Id       int    `json:"id"`
}

type allCache struct {
	users *cache.Cache
}

const (
	defaultExpiration = 5 * time.Minute
	purgeTime         = 10 * time.Minute
)

func newCache() *allCache {
	Cache := cache.New(defaultExpiration, purgeTime)
	return &allCache{
		users: Cache,
	}
}

func (c *allCache) update(id string, user UserCache) {
	c.users.Set(id, user, cache.DefaultExpiration)
}

var c = newCache()

func doesUserExist(id string) bool {
	fmt.Print("it does exists")
	user, ok := c.users.Get(id)
	if ok {
		fmt.Print("He exists oh!")
		fmt.Printf("%s", user)
		return true
	} else {
		return false
	}
}
