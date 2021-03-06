package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "regexp"
)

type SlackMessageAttachment struct {
    Service_name string
    Service_url string
    Title string
    Title_link string
    Thumb_url string
    Thumb_width int
    Thumb_height int
    Fallback string
    Audio_html string
    Audio_html_width int
    Audio_html_height int
    From_url string
    Id int
}

type SlackMessageReaction struct {
    Name string
    Users []string
    Count int
}

type SlackMessage struct {
    User string
    Subtype string
    Text string
    Ts string
    Attachments []SlackMessageAttachment
    Reactions []SlackMessageReaction
}

type SlackRecommendationResponse struct {
    Ok bool
    Messages []SlackMessage
    Has_more bool
}

type SlackProfile struct {
    Image_24 string
    Image_32 string
    Image_48 string
    Image_72 string
    Image_192 string
    First_name string
    Last_name string
    Title string
    Skype string
    Phone string
    Real_name string
    Real_name_normalized string
    Email string
}

type SlackUser struct {
    Id string
    Name string
    Deleted bool
    Status string
    Color string
    Real_name string
    Profile SlackProfile
}

type SlackUserResponse struct {
    Ok bool
    User SlackUser
}

func GetRecommendationMessages(params ...string) ([]SlackMessage, bool) {
    var url string
    if len(params) > 1 {
        url = fmt.Sprintf("https://slack.com/api/channels.history?token=%s&channel=%s&count=1000&oldest=%s&newest=%s", SLACK_API_KEY, SLACK_CHANNEL_ID, params[0], params[1])
    } else {
        url = fmt.Sprintf("https://slack.com/api/channels.history?token=%s&channel=%s&count=1000", SLACK_API_KEY, SLACK_CHANNEL_ID)
    }
    log.Println(url)
    response, err := http.Get(url)
    defer response.Body.Close()
    checkErr(err)
    contents, err := ioutil.ReadAll(response.Body)
    checkErr(err)
    responseData := SlackRecommendationResponse{}
    err = json.Unmarshal([]byte(contents), &responseData)
    checkErr(err)
    return responseData.Messages, responseData.Has_more
}

func GetUserFromSlack(userId string) (SlackUser, error) {
    response, err := http.Get(fmt.Sprintf("https://slack.com/api/users.info?token=%s&user=%s", SLACK_API_KEY, userId))
    if err != nil {
        return SlackUser{}, err
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            return SlackUser{}, err
        }
        responseData := SlackUserResponse{}
        err = json.Unmarshal([]byte(contents), &responseData)
        if err != nil {
            return SlackUser{}, err
        }
        return responseData.User, nil
    }
}

func ParseTags(messageText string) ([]string) {
    re, _ := regexp.Compile(":([a-z]+):")
    result := re.FindAllStringSubmatch(messageText, -1)
    matches := []string{}
    for _, v := range result {
        matches = append(matches, v[1])
    }
    return matches
}
