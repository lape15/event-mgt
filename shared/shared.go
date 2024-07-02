package shared

import (
	"fmt"
	"strconv"
	"time"

	"event-mgt/database"

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

type Credential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
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

func (c *allCache) Update(id string, user UserCache) {
	c.users.Set(id, user, cache.DefaultExpiration)
}

var C = newCache()

func DoesUserExist(id string) bool {
	fmt.Print("it does exists")
	user, ok := C.users.Get(id)
	if ok {
		fmt.Print("He exists oh!")
		fmt.Printf("%s", user)
		return true
	} else {
		return false
	}
}

func IsUserHostOfEvent(userID string, eventID int) bool {

	var hostID int
	row, err := database.Db.QueryRow("SELECT user_id FROM events WHERE id = ?", eventID)
	row.Scan(&hostID)
	if err != nil {
		fmt.Println("Error checking host:", err)
		return false
	}
	userIDInt, _ := strconv.Atoi(userID)
	return userIDInt == hostID
}

func GetEventRSVPCount(eventID string) (int, error) {
	var count int
	row, err := database.Db.QueryRow("SELECT COUNT(*) FROM event_rsvps WHERE event_id = ?", eventID)
	row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
