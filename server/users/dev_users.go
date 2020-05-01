package users

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
)

type (
	User struct {
		Password string
		Name string
		Id uuid.UUID
		Token string
	}

	UserFinder interface {
		FindByName(name string) User
		FindByToken(token string) User
		FindById(uuid uuid.UUID) User
	}

	DevUsers struct {
		byName map[string] *User
		byToken map[string] *User
		byId map[uuid.UUID] *User
	}
)

func NewDevUsers(filename string) *DevUsers {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Print(err)
		return nil
	}
	users := DevUsers{
		byName : make(map[string] *User),
		byToken : make(map[string] *User),
	}
	var usersOnDisk = make([]User, 0)
	err = json.Unmarshal(data, &usersOnDisk)
	if err != nil {
		return nil
	}
	for _, onDisk := range usersOnDisk {
		user := onDisk
		user.GenerateInsecureDevToken()
		users.byToken[user.Token] = &user
		users.byName[user.Name] = &user
	}
	return &users
}

func (user *User) GenerateInsecureDevToken() {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprint(user.Name, ":", user.Password)))
	user.Token = base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

func (users *DevUsers) FindByName(name string) User {
	return *users.byName[name]
}

func (users *DevUsers) FindByToken(token string) User {
	return *users.byToken[token]
}

func (users *DevUsers) FindById(uuid uuid.UUID) User {
	return *users.byId[uuid]
}

func (users *DevUsers) ValidCredentials(username, password string) bool {
	if user, ok := users.byName[username]; ok {
		return user.Password == password
	}
	return false
}

func (users *DevUsers) ValidateToken(token string) bool {
	_, ok := users.byToken[token]
	return ok
}