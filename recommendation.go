package main

import (
    "fmt"
    "strconv"
)

type Recommendation struct {
    Id int
    User_id int
    Post_time float64
    Service_name string
    Title string
    Url string
    Thumb_url string
    Thumb_width int
    Thumb_height int
    Audio_html string
    Audio_html_width int
    Audio_html_height int
}

func CreateRecommendation(slack_message SlackMessage) (Recommendation) {
    db, err := GetDatabase()
    defer db.Close()
    checkErr(err)
    var recommendation_id int
    post_time, err := strconv.ParseFloat(slack_message.Ts, 64)
    checkErr(err)
    slack_user, err := GetUserFromSlack(slack_message.User)
    checkErr(err)
    user, _ := GetOrCreateUser(slack_user)
    attachment := slack_message.Attachments[0]
    err = db.QueryRow("INSERT INTO Recommendations (user_id, post_time, service_name, title, url, thumb_url, thumb_width, thumb_height, audio_html, audio_height, audio_width) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) returning id;", user.Id, post_time, attachment.Service_name, attachment.Title, attachment.Title_link, attachment.Thumb_url, attachment.Thumb_width, attachment.Thumb_height, attachment.Audio_html, attachment.Audio_html_width, attachment.Audio_html_height).Scan(&recommendation_id)
    checkErr(err)
    recommendation := Recommendation{recommendation_id, user.Id, post_time, attachment.Service_name, attachment.Title, attachment.Title_link, attachment.Thumb_url, attachment.Thumb_width, attachment.Thumb_height, attachment.Audio_html, attachment.Audio_html_width, attachment.Audio_html_height}
    return recommendation
}

func DownloadRecommendations(params ...float64) (float64) {
    sourceWhitelist := map[string]bool {
        "SoundCloud": true,
        "Spotify": true,
        "YouTube": true,
    }
    messages, hasMore := GetRecommendationMessages(params...)
    var err error
    var newestTime float64
    var postTime float64
    if len(messages) > 0 {
        postTime, err = strconv.ParseFloat(messages[0].Ts, 64)
        checkErr(err)
        newestTime, err = strconv.ParseFloat(messages[len(messages) - 1].Ts, 64)
        checkErr(err)
    }
    for _, message := range messages {
        if message.Attachments != nil && sourceWhitelist[message.Attachments[0].Service_name] {
            GetOrCreateRecommendation(message)
        }
    }
    if hasMore {
        // TODO: recurse over additional pages. the code below does not work
        DownloadRecommendations(postTime, newestTime)
    }
    return postTime
}

func GetRecommendationById(recommendation_id int) (*Recommendation, error) {
    db, err := GetDatabase()
    defer db.Close()
    rows, err := db.Query("SELECT * FROM Recommendations WHERE id = $1;", recommendation_id)
    checkErr(err)
    for rows.Next() {
        var id int
        var user_id int
        var post_time float64
        var service_name string
        var title string
        var url string
        var thumb_url string
        var thumb_width int
        var thumb_height int
        var audio_html string
        var audio_html_width int
        var audio_html_height int
        err = rows.Scan(&id, &user_id, &post_time, &service_name, &title, &url, &thumb_url, &thumb_width, &thumb_height, &audio_html, &audio_html_width, &audio_html_height)
        checkErr(err)
        return &Recommendation{id, user_id, post_time, service_name, title, url, thumb_url, thumb_width, thumb_height, audio_html, audio_html_width, audio_html_height}, err
    }
    return nil, fmt.Errorf("Recommendation not found!")
}

func GetOrCreateRecommendation(slack_message SlackMessage) (Recommendation, bool) {
    db, err := GetDatabase()
    defer db.Close()
    slack_user, err := GetUserFromSlack(slack_message.User)
    user, _ := GetOrCreateUser(slack_user)
    rows, err := db.Query("SELECT * FROM Recommendations WHERE user_id = $1 AND url = $2;", user.Id, slack_message.Attachments[0].Title_link)
    checkErr(err)
    for rows.Next() {
        var id int
        var user_id int
        var post_time float64
        var service_name string
        var title string
        var url string
        var thumb_url string
        var thumb_width int
        var thumb_height int
        var audio_html string
        var audio_html_width int
        var audio_html_height int
        err = rows.Scan(&id, &user_id, &post_time, &service_name, &title, &url, &thumb_url, &thumb_width, &thumb_height, &audio_html, &audio_html_width, &audio_html_height)
        checkErr(err)
        return Recommendation{id, user_id, post_time, service_name, title, url, thumb_url, thumb_width, thumb_height, audio_html, audio_html_width, audio_html_height}, false
    }
    return CreateRecommendation(slack_message), true
}

func GetRecommendationMaxTime() (float64) {
    db, err := GetDatabase()
    defer db.Close()
    rows, err := db.Query("SELECT MAX(post_time) FROM Recommendations;")
    for rows.Next() {
        var max_time float64
        err = rows.Scan(&max_time)
        checkErr(err)
        return max_time
    }
    panic("Unable to get max time")
}

func GetRecommendationMinTime() (float64) {
    db, err := GetDatabase()
    defer db.Close()
    rows, err := db.Query("SELECT MIN(post_time) FROM Recommendations;")
    for rows.Next() {
        var max_time float64
        err = rows.Scan(&max_time)
        checkErr(err)
        return max_time
    }
    panic("Unable to get min time")
}
