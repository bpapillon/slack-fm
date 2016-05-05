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
    Tags []string
}

type RecommendationTag struct {
    Id int
    Recommendation_id int
    Tag_id int
    User_id int
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
    var tags []string
    recommendation := Recommendation{recommendation_id, user.Id, post_time, attachment.Service_name, attachment.Title, attachment.Title_link, attachment.Thumb_url, attachment.Thumb_width, attachment.Thumb_height, attachment.Audio_html, attachment.Audio_html_width, attachment.Audio_html_height, tags}
    return recommendation
}

func DownloadRecommendations(params ...string) (string) {
    sourceWhitelist := map[string]bool {
        "SoundCloud": true,
        "Spotify": true,
        "YouTube": true,
    }
    messages, hasMore := GetRecommendationMessages(params...)
    var newestTime string
    var postTime string
    if len(messages) > 0 {
        postTime = messages[0].Ts
        newestTime = messages[len(messages) - 1].Ts
    }
    for _, message := range messages {
        if message.Attachments != nil && sourceWhitelist[message.Attachments[0].Service_name] {
            recommendation, _ := GetOrCreateRecommendation(message)
            tags := ParseTags(message.Text)
            for i := 0; i < len(tags); i++ {
                GetOrCreateRecommendationTag(recommendation.Id, tags[i], recommendation.User_id)
            }
            for i := 0; i < len(message.Reactions); i++ {
                for j := 0; j < len(message.Reactions[i].Users); j++ {
                    slack_user, err := GetUserFromSlack(message.Reactions[i].Users[j])
                    checkErr(err)
                    user, _ := GetOrCreateUser(slack_user)
                    GetOrCreateRecommendationTag(recommendation.Id, message.Reactions[i].Name, user.Id)
                }
            }
        }
    }
    if hasMore {
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
        tags := GetRecommendationTags(id)
        return &Recommendation{id, user_id, post_time, service_name, title, url, thumb_url, thumb_width, thumb_height, audio_html, audio_html_width, audio_html_height, tags}, err
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
        tags := GetRecommendationTags(id)
        return Recommendation{id, user_id, post_time, service_name, title, url, thumb_url, thumb_width, thumb_height, audio_html, audio_html_width, audio_html_height, tags}, false
    }
    return CreateRecommendation(slack_message), true
}

func CreateRecommendationTag(recommendation_id int, tag_id int, user_id int) (RecommendationTag) {
    db, err := GetDatabase()
    defer db.Close()
    checkErr(err)
    var recommendation_tag_id int
    err = db.QueryRow("INSERT INTO Recommendation_tags (tag_id, recommendation_id, user_id) VALUES ($1, $2, $3) returning id;", tag_id, recommendation_id, user_id).Scan(&recommendation_tag_id)
    checkErr(err)
    recommendation_tag := RecommendationTag{recommendation_tag_id, recommendation_id, tag_id, user_id}
    return recommendation_tag
}

func GetOrCreateRecommendationTag(recommendation_id int, tag_string string, user_id int) (RecommendationTag, bool) {
    db, err := GetDatabase()
    defer db.Close()
    tag, _ := GetOrCreateTag(tag_string)
    rows, err := db.Query("SELECT * FROM Recommendation_tags WHERE tag_id = $1 and recommendation_id = $2 and user_id = $3;", tag.Id, recommendation_id, user_id)
    checkErr(err)
    for rows.Next() {
        var id int
        var tag_id int
        err = rows.Scan(&id, &tag_id, &recommendation_id, &user_id)
        checkErr(err)
        return RecommendationTag{id, tag_id, recommendation_id, user_id}, false
    }
    return CreateRecommendationTag(recommendation_id, tag.Id, user_id), true
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

func GetRecommendationTags(recommendation_id int) ([]string) {
    db, err := GetDatabase()
    defer db.Close()
    rows, err := db.Query("SELECT DISTINCT tags.tag FROM recommendation_tags INNER JOIN tags ON tags.id = recommendation_tags.tag_id WHERE recommendation_tags.recommendation_id = $1;", recommendation_id)
    checkErr(err)
    var tags []string
    for rows.Next() {
        var tag string
        err = rows.Scan(&tag)
        checkErr(err)
        tags = append(tags, tag)
    }
    return tags
}

func GetRecommendationsByTag(tag string) ([]Recommendation) {
    db, err := GetDatabase()
    defer db.Close()
    rows, err := db.Query("SELECT DISTINCT recommendations.* FROM recommendations INNER JOIN recommendation_tags ON recommendation_tags.recommendation_id = recommendations.id INNER JOIN tags ON tags.id = recommendation_tags.tag_id WHERE tags.tag = $1;", tag)
    checkErr(err)
    var recommendations []Recommendation
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
        tags := GetRecommendationTags(id)
        recommendations = append(recommendations, Recommendation{id, user_id, post_time, service_name, title, url, thumb_url, thumb_width, thumb_height, audio_html, audio_html_width, audio_html_height, tags})
    }
    return recommendations
}

func GetRecommendationsByUser(user_id int) ([]Recommendation) {
    db, err := GetDatabase()
    defer db.Close()
    rows, err := db.Query("SELECT * FROM recommendations WHERE user_id = $1;", user_id)
    checkErr(err)
    var recommendations []Recommendation
    for rows.Next() {
        var id int
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
        tags := GetRecommendationTags(id)
        recommendations = append(recommendations, Recommendation{id, user_id, post_time, service_name, title, url, thumb_url, thumb_width, thumb_height, audio_html, audio_html_width, audio_html_height, tags})
    }
    return recommendations
}
