package vpn

import (
	"astro/persistence"
	"astro/vpn/types"
	"github.com/dgraph-io/badger/v3"
	"github.com/goccy/go-json"
)

func SaveUser(db *badger.DB, user types.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return persistence.Save(db, user.Username, string(data))
}

func RemoveUser(db *badger.DB, username string) error {
	return persistence.Delete(db, username)
}

func loadUser(db *badger.DB, username string) (types.User, error) {
	bytes, err := persistence.Get(db, username)
	if err != nil {
		return types.User{}, err
	}

	var user types.User
	err = json.Unmarshal(bytes, &user)
	return user, err
}

func GetAllUsers(db *badger.DB) ([]types.User, error) {
	mapUser, err := persistence.GetAll(db)
	if err != nil {
		return nil, err
	}

	var users []types.User

	for _, v := range mapUser {
		var user types.User
		err = json.Unmarshal([]byte(v), &user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}