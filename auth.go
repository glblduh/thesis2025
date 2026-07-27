package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	"go.etcd.io/bbolt"
	"golang.org/x/crypto/bcrypt"
)

const (
	USERS_BUCKET = "Users"
	USERS_PASSWORDS_BUCKET = "UsersPasswords"
	USERS_API_BUCKET = "UsersAPIKeys"
)

func openAuthDB() (*bbolt.DB, error) {
	return bbolt.Open(".auth.db", 0660, &bbolt.Options{
		Timeout: time.Second,
	})
}

func initializeAuthDB() {
	Info.Println("Initializing Auth DB")

	db, dbErr := openAuthDB()
	if dbErr != nil {
		Error.Fatalln("cannot open auth db: " + dbErr.Error())
	}
	defer db.Close()

	db.Update(func(tx *bbolt.Tx) error {
		_, usersBucketCreateErr := tx.CreateBucketIfNotExists([]byte(USERS_BUCKET))
		if usersBucketCreateErr != nil {
			Error.Fatalln("cannot create users bucket: " + usersBucketCreateErr.Error())
		}

		_, usersHashedPasswordsBucketCreateErr := tx.CreateBucketIfNotExists([]byte(USERS_PASSWORDS_BUCKET))
		if usersHashedPasswordsBucketCreateErr != nil {
			Error.Fatalln("cannot create users hashed passwords bucket: " + usersHashedPasswordsBucketCreateErr.Error())
		}

		_, usersAPIKeysBucketCreateErr := tx.CreateBucketIfNotExists([]byte(USERS_API_BUCKET))
		if usersAPIKeysBucketCreateErr != nil {
			Error.Fatalln("cannot create users api keys bucket: " + usersAPIKeysBucketCreateErr.Error())
		}

		return nil
	})

	Info.Println("Successfully initialized Auth DB")
}

func generateAPIKey() (string, error) {
	key := ""
	bytesBuf := make([]byte, 32)

	_, generateErr := rand.Read(bytesBuf)
	if generateErr != nil {
		return key, generateErr
	}

	key = base64.RawURLEncoding.EncodeToString(bytesBuf)
	return key, nil
}

func authCreateUser(username string, password string, userType UserType) {
	db, dbErr := openAuthDB()
	if dbErr != nil {
		Error.Fatalln("cannot open auth db: " + dbErr.Error())
	}
	defer db.Close()

	db.Update(func(tx *bbolt.Tx) error {
		usersBucket := tx.Bucket([]byte(USERS_BUCKET))
		usersPasswordsBucket := tx.Bucket([]byte(USERS_PASSWORDS_BUCKET))
		usersAPIKeysBucket := tx.Bucket([]byte(USERS_API_BUCKET))

		passwordHash, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if hashErr != nil {
			Error.Fatalln("password hashing error: " + hashErr.Error())
		}

		userPutErr := usersBucket.Put([]byte(username), []byte(userType))
		if userPutErr != nil {
			Error.Fatalln("cannot create user: " + userPutErr.Error())
		}

		userPasswordPutErr := usersPasswordsBucket.Put([]byte(username), passwordHash)
		if userPasswordPutErr != nil {
			Error.Fatalln("cannot create user: " + userPasswordPutErr.Error())
		}

		key, genErr := generateAPIKey()
		if genErr != nil {
			Error.Fatalln("cannot create user: " + genErr.Error())
		}

		userAPIPutErr := usersAPIKeysBucket.Put([]byte(username), []byte(key))
		if userAPIPutErr != nil {
			Error.Fatalln("cannot create user: " + userAPIPutErr.Error())
		}

		return nil
	})
}

func authDeleteUser(username string) {
	db, dbErr := openAuthDB()
	if dbErr != nil {
		Error.Fatalln("cannot open auth db: " + dbErr.Error())
	}
	defer db.Close()

	db.Update(func(tx *bbolt.Tx) error {
		usersBucket := tx.Bucket([]byte(USERS_BUCKET))
		usersPasswordsBucket := tx.Bucket([]byte(USERS_PASSWORDS_BUCKET))
		usersAPIKeysBucket := tx.Bucket([]byte(USERS_API_BUCKET))

		userDeleteErr := usersBucket.Delete([]byte(username))
		if userDeleteErr != nil {
			Error.Fatalln("cannot delete user: " + userDeleteErr.Error())
		}

		userPasswordDeleteErr := usersPasswordsBucket.Delete([]byte(username))
		if userPasswordDeleteErr != nil {
			Error.Fatalln("cannot delete user: " + userPasswordDeleteErr.Error())
		}

		userAPIDeleteErr := usersAPIKeysBucket.Delete([]byte(username))
		if userAPIDeleteErr != nil {
			Error.Fatalln("cannot delete user: " + userAPIDeleteErr.Error())
		}

		return nil
	})
}

func validateAuth(username string, password string) (userAuth, error) {
	userInfo := userAuth{}

	db, dbErr := openAuthDB()
	if dbErr != nil {
		Error.Fatalln("cannot open auth db: " + dbErr.Error())
	}
	defer db.Close()

	return userInfo, db.View(func(tx *bbolt.Tx) error {
		usersBucket := tx.Bucket([]byte(USERS_BUCKET))
		usersPasswordsBucket := tx.Bucket([]byte(USERS_PASSWORDS_BUCKET))
		usersAPIKeysBucket := tx.Bucket([]byte(USERS_API_BUCKET))

		user := usersBucket.Get([]byte(username))
		userPassword := usersPasswordsBucket.Get([]byte(username))
		userKey := usersAPIKeysBucket.Get([]byte(username))

		if user == nil || userPassword == nil || userKey == nil {
			return ErrAuthUserNotFound
		}

		userInfo.Username = username
		userInfo.Type = UserType(user)
		userInfo.Key = string(userKey)

		return bcrypt.CompareHashAndPassword(userPassword, []byte(password))
	})
}

func validateKey(username string, key string) error {
	db, dbErr := openAuthDB()
	if dbErr != nil {
		Error.Fatalln("cannot open auth db: " + dbErr.Error())
	}
	defer db.Close()

	return db.View(func(tx *bbolt.Tx) error {
		usersAPIKeysBucket := tx.Bucket([]byte(USERS_API_BUCKET))

		user := usersAPIKeysBucket.Get([]byte(username))
		if user == nil {
			return ErrAuthUserNotFound
		}

		if key != string(user) {
			return ErrAuthKeyIncorrect
		}

		return nil
	})
}

func apiAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errorRes(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		splittedUserKey := strings.Split(authHeader, ":")
		keyValidateErr := validateKey(splittedUserKey[0], splittedUserKey[1])
		if keyValidateErr != nil {
			status := http.StatusUnauthorized

			if errors.Is(keyValidateErr, ErrAuthUserNotFound) {
				status = http.StatusNotFound
			}

			errorRes(w, keyValidateErr.Error(), status)
			return
		}

		next.ServeHTTP(w, r)
	})
}
