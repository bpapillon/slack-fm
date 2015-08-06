package main

import (
    "github.com/mrvdot/golang-utils"
)

type User struct {
    Id int
    Slack_id string
    Photo_url string
    Name string
    Slug string
}

func CreateUser(slack_user SlackUser) (User) {
    db, err := GetDatabase()
    defer db.Close()
    checkErr(err)
    var user_id int
    slug := utils.GenerateSlug(slack_user.Real_name)
    err = db.QueryRow("INSERT INTO Users (slack_id, photo_url, name, slug) VALUES ($1, $2, $3, $4) returning id;", slack_user.Id, slack_user.Profile.Image_192, slack_user.Real_name, slug).Scan(&user_id)
    checkErr(err)
    user := User{user_id, slack_user.Id, slack_user.Profile.Image_192, slack_user.Real_name, slug}
    return user
}

func GetOrCreateUser(slack_user SlackUser) (User, bool) {
    db, err := GetDatabase()
    defer db.Close()
    rows, err := db.Query("SELECT * FROM Users WHERE slack_id = $1;", slack_user.Id)
    checkErr(err)
    for rows.Next() {
        var id int
        var slack_id string
        var photo_url string
        var name string
        var slug string
        err = rows.Scan(&id, &slack_id, &photo_url, &name, &slug)
        checkErr(err)
        return User{id, slack_id, photo_url, name, slug}, false
    }
    return CreateUser(slack_user), true
}
