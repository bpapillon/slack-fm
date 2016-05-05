package main

type Tag struct {
    Id int
    Tag string
}

type TagCount struct {
    Tag_id int
    Tag string
    Count int
}

func CreateTag(tag string) (Tag) {
    db, err := GetDatabase()
    defer db.Close()
    checkErr(err)
    var tag_id int
    err = db.QueryRow("INSERT INTO Tags (tag) VALUES ($1) returning id;", tag).Scan(&tag_id)
    checkErr(err)
    recommendation := Tag{tag_id, tag}
    return recommendation
}

func GetOrCreateTag(tag string) (Tag, bool) {
    db, err := GetDatabase()
    defer db.Close()
    rows, err := db.Query("SELECT * FROM Tags WHERE tag = $1;", tag)
    checkErr(err)
    for rows.Next() {
        var id int
        err = rows.Scan(&id, &tag)
        checkErr(err)
        return Tag{id, tag}, false
    }
    return CreateTag(tag), true
}

func GetTagCounts() ([]TagCount) {
    db, err := GetDatabase()
    defer db.Close()
    rows, err := db.Query("SELECT recommendation_tags.tag_id, tags.tag,  COUNT(*) AS n FROM recommendation_tags INNER JOIN tags ON (tags.id = recommendation_tags.tag_id) GROUP BY recommendation_tags.tag_id, tags.tag ORDER BY n DESC;")
    checkErr(err)
    var counts []TagCount
    for rows.Next() {
        var tag_id int
        var tag string
        var count int
        err = rows.Scan(&tag_id, &tag, &count)
        checkErr(err)
        counts = append(counts, TagCount{tag_id, tag, count})
    }
    return counts
}
