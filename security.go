package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// all data available in a single row in the users table
type Row struct {
	ID   int    `db:"id"`
	Hash string `db:"hash"`
	User string `db:"user"`
}

// check a provided (hashed) password against a user's stored password. returns a bool and the hash on success
func userIsAuthenticated(username string, testHash []byte) (bool, []byte) {
	// without rate limiting, it's vulnerable to brute force attacks, but for the sake of time i'll forgo things like rate limiting for now...
	// if i did implement it, i would track login attempts against an account over a specific span of time in memory, maybe using something like redis for convenience. maybe ip based rate limiting too

	user := new(Row)

	err := dbcon.SQL().SelectFrom("users").Where("user=?", username).One(user)
	if err != nil { // exit on fail
		fmt.Println(err)
		return false, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), testHash)

	if err != nil {
		return false, nil
	}

	return true, []byte(user.Hash)
}

// creates the default account + password if it doesn't exist. TODO: if it does, just update the default password. don't run more than once...
func setDefaultPassword(defaultPassword string) {

	_, err := os.Stat(".newuserlockout")
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("Default user account already generated, skipping..")
			return //we've already created the default user. skip this
		}
	}

	fmt.Printf("Setting new password for account `default` to `%+v`\n", defaultPassword)

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	/*
		row := new(Row)
		defaultUser := dbcon.Collection("users")

		res := defaultUser.Find(db.Cond{"user": "default"})
		err = res.One(row)

		if err != nil && row.User != "default" { // default user doesn't exist, create

		}
	*/

	_, err = dbcon.SQL().InsertInto("users").Columns(
		"hash",
		"user",
	).Values(
		string(pwdHash),
		"default",
	).Exec()

	if err != nil {
		fmt.Println("Error setting password!")
		fmt.Println(err)
	} else {
		fmt.Println("Successfully set new password.")
		os.OpenFile(".newuserlockout", os.O_RDONLY|os.O_CREATE, 0644)
	}

}
