package persistence

import (
	os2 "astro/os"
	"astro/utils"
	"github.com/dgraph-io/badger/v3"
)

type (
	Persists struct {
		DBUsers   *badger.DB
		DBRoles   *badger.DB
		DBModules *badger.DB
		DBActions *badger.DB
		DBConfig  *badger.DB
		DbVPN     *badger.DB
	}
)

func NewPersists() (Persists, error) {
	dbUsers, err := NewPersistence("users")
	if err != nil {
		return Persists{}, err
	}

	dbRoles, err := NewPersistence("roles")
	if err != nil {
		return Persists{}, err
	}

	dbModules, err := NewPersistence("modules")
	if err != nil {
		return Persists{}, err
	}

	dbActions, err := NewPersistence("actions")
	if err != nil {
		return Persists{}, err
	}

	dbConfig, err := NewPersistence("config")
	if err != nil {
		return Persists{}, err
	}

	dbVPN, err := NewPersistence("vpn")
	if err != nil {
		return Persists{}, err
	}

	return Persists{
		DBUsers:   dbUsers,
		DBRoles:   dbRoles,
		DBModules: dbModules,
		DBActions: dbActions,
		DBConfig:  dbConfig,
		DbVPN:     dbVPN,
	}, nil
}

func NewPersistence(dbname string) (*badger.DB, error) {

	existsKey := utils.ExistsEncryptionKey(os2.ASTRO_KEY_PATH)

	var key []byte
	var err error

	if !existsKey {
		key, err = utils.GenerateAndSaveEncryptionKey(os2.ASTRO_KEY_PATH)
		if err != nil {
			return nil, err
		}
	} else {
		key, err = utils.LoadEncryptionKey(os2.ASTRO_KEY_PATH)
		if err != nil {
			return nil, err
		}
	}

	opt := badger.DefaultOptions(os2.ASTRO_ETC_PATH + dbname + ".astro").WithEncryptionKey(key)
	opt.IndexCacheSize = 200 << 20

	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return db, nil
}

func Save(db *badger.DB, key, value string) error {
	return db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), []byte(value))
		return err
	})
}

func Get(db *badger.DB, key string) ([]byte, error) {
	var value []byte
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		valBytes, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		value = valBytes
		return nil
	})

	if err != nil {
		return nil, err
	}
	return value, nil
}

func Delete(db *badger.DB, key string) error {
	return db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))
		return err
	})
}

func GetAll(db *badger.DB) (map[string]string, error) {

	values := map[string]string{}

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := string(item.Key())
			value, _ := item.ValueCopy(nil)

			values[key] = string(value)
		}

		return nil
	})
	return values, err
}
