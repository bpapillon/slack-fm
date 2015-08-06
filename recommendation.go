package main

import (
    "fmt"
    "strconv"
    "time"
)

type Recommendation struct {
    Id int
    User_id int
    Post_time time.Time
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

// func DownloadRecommendations() (error) {
//     messages := GetRecommendationMessages()
//     for _, message := range messages {
//         for _, attachment := range message.Attachments {
//
//         }
//     }
// }

func UnixTimeStringToTime(unix_time string) (time.Time, error) {
    i, err := strconv.ParseInt(unix_time, 10, 64)
    checkErr(err)
    return time.Unix(i, 0), err
}

func CreateRecommendation(slack_message SlackMessage) (Recommendation) {
    db, err := GetDatabase()
    defer db.Close()
    checkErr(err)
    var recommendation_id int
    post_time, err := UnixTimeStringToTime(slack_message.Ts)
    slack_user, err := GetUserFromSlack(slack_message.User)
    user, _ := GetOrCreateUser(slack_user)
    attachment := slack_message.Attachments[0]
    err = db.QueryRow("INSERT INTO Recommendations (user_id, post_time, service_name, title, url, thumb_url, thumb_width, thumb_height, audio_html, audio_height, audio_width) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) returning id;", user.Id, post_time, attachment.Service_name, attachment.Title, attachment.Title_link, attachment.Thumb_url, attachment.Thumb_width, attachment.Thumb_height, attachment.Audio_html, attachment.Audio_html_width, attachment.Audio_html_height).Scan(&recommendation_id)
    checkErr(err)
    recommendation := Recommendation{recommendation_id, user.Id, post_time, attachment.Service_name, attachment.Title, attachment.Title_link, attachment.Thumb_url, attachment.Thumb_width, attachment.Thumb_height, attachment.Audio_html, attachment.Audio_html_width, attachment.Audio_html_height}
    return recommendation
}

func GetRecommendationById(recommendation_id int) (*Recommendation, error) {
    db, err := GetDatabase()
    defer db.Close()
    rows, err := db.Query("SELECT * FROM Recommendations WHERE id = $1;", recommendation_id)
    checkErr(err)
    for rows.Next() {
        var id int
        var user_id int
        var post_time time.Time
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
        var post_time time.Time
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
