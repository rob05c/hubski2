package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

func dbGet(connectionUri string) (*sql.DB, error) {
	return sql.Open("postgres", connectionUri)
}

// dbGetTokenUser returns the user for the given token, whether any user was found, or any database errors
func dbGetTokenUser(token string, db *sql.DB) (string, bool, error) {
	user := ""
	err := db.QueryRow("select username from tokens where token = $1", token).Scan(&user)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return user, true, nil
}

func dbInsertToken(token, user string, db *sql.DB) error {
	_, err := db.Exec("insert into tokens (token, username) values ($1, $2);", token, user)
	return err
}

// dbGetDonations returns the donation amount in cents
// TODO(move to its own package/file)
// TODO(use prepared statement)
func dbGetDonations(user string, db *sql.DB) (int, error) {
	donationCents := 0
	err := db.QueryRow("select sum(donation_cents) from donations where username = $1", user).Scan(&donationCents)
	return donationCents, err
}

type DbUserPageUserData struct {
	Name                 string    `json:"name"`
	Badges               int       `json:"badges"`
	Followers            int       `json:"followers"`
	Bio                  string    `json:"bio"`
	Following            int       `json:"following"`
	FollowedTags         int       `json:"followed_tags"`
	FollowedDomains      int       `json:"followed_domains"`
	BadgesGiven          int       `json:"badges_given"`
	BadgesTotal          int       `json:"badges_total"`
	Joined               time.Time `json:"joined"`
	Style                string    `json:"style"`
	DonatedCents         int       `json:"donated_cents"`
	DonatedCentsThisYear int       `json:"donated_cents_this_year"`
	TagsUsed             []string  `json:"tags_used"`
}

func dbGetUserLastFiveTags(user string, db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT tag FROM (SELECT DISTINCT(tag) as tag, MAX(time) as time FROM publications as pubs WHERE username = $1 AND tag IS NOT NULL GROUP BY tag ORDER BY MAX(time) DESC LIMIT 5) as tags", user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []string{}
	for rows.Next() {
		tag := ""
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tags, nil
}

// dbGetUserPageUserData returns the user page data for the given user, whether the user exists, or any error.
// TODO(add tags used)
func dbGetUserPageUserData(user string, db *sql.DB) (DbUserPageUserData, bool, error) {
	q := "select followers, badges, bio, following, followed_tags, followed_domains, badges_given, badges_total, joined, style, donated_cents, donated_cents_this_year from user_page_user_data where id = $1"

	d := DbUserPageUserData{Name: user}
	err := db.QueryRow(q, user).Scan(&d.Followers, &d.Badges, &d.Bio, &d.Following, &d.FollowedTags, &d.FollowedDomains, &d.BadgesGiven, &d.BadgesTotal, &d.Joined, &d.Style, &d.DonatedCents, &d.DonatedCentsThisYear)
	if err == sql.ErrNoRows {
		return DbUserPageUserData{}, false, nil
	}

	tags, err := dbGetUserLastFiveTags(user, db)
	if err == sql.ErrNoRows {
		return DbUserPageUserData{}, false, nil
	}
	d.TagsUsed = tags

	return d, true, err
}

// TODO(rename)
type DbComment struct {
	Id        int       `json:"id"`
	User      string    `json:"user"`
	Time      time.Time `json:"time"`
	PostTitle string    `json:"post_title"`
	PostId    string    `json:"post_id"`
	Text      string    `json:"md"`
	Score     int       `json:"score"`
	Deleted   bool      `json:"deleted"`
	Draft     bool      `json:"draft"`
	ParentId  int       `json:"parent_id"`
}

func dbGetUserPageComments(user string, numComments int, db *sql.DB) ([]DbComment, error) {
	q := "select id, time, text, score, deleted, draft, parent_id, post_id, post_title from user_page_comments where username = $1 order by time desc fetch first 12 rows only;"
	comments := []DbComment(nil)
	rows, err := db.Query(q, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		c := DbComment{User: user}
		if err := rows.Scan(&c.Id, &c.Time, &c.Text, &c.Score, &c.Deleted, &c.Draft, &c.ParentId, &c.PostId, &c.PostTitle); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

type DbUserPagePostSharePublication struct {
	Id           int    `json:"id"`
	Title        string `json:"title"`
	User         string `json:"user"`
	Score        int    `json:"score"`
	Tag          string `json:"tag"`
	Tag2         string `json:"tag2"`
	CommunityTag string `json:"community_tag"`
	Domain       string `json:"domain"`
	CommentCount int    `json:"comment_count"`
	BadgeCount   int    `json:"badge_count"`
}

// TODO make this a view
func dbGetPublicationChildren(id int, db *sql.DB) ([]int, error) {
	q := `WITH RECURSIVE publication_children(id, parent_id)
AS (
  -- get parent
  SELECT id, parent_id
  FROM publications where id = $1
  UNION ALL
  -- get all children
  SELECT t.id, t.parent_id
  FROM publication_children c
  JOIN publications t
  ON t.parent_id = c.id
)
SELECT id, parent_id from publication_children where parent_id is not null;`

	publications := []int{}
	rows, err := db.Query(q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var parent_id int
		err := rows.Scan(&id, &parent_id)
		if err != nil {
			return nil, err
		}
		publications = append(publications, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return publications, err
}

// TODO(CommunityTag, Badges, CommentCount)
func dbGetUserPagePostsAndShares(user string, numPostsAndShares int, db *sql.DB) ([]DbUserPagePostSharePublication, error) {
	// select shares UNION posts
	// TODO make view?
	q := `(
SELECT p.id, p.time, p.title, p.username, p.score, p.tag, COALESCE(p.tag2, ''), COALESCE(p.web_domain, ''), COALESCE(c.tag, '') AS community_tag, COUNT(b.id) AS badge_count
FROM publications AS p
INNER JOIN users_shareds AS s ON p.id = s.publication_id
INNER JOIN publication_community_tag AS c on c.id = p.id
LEFT OUTER JOIN publication_badged_by AS b ON p.id = b.id
WHERE s.id = $1 AND p.type_id = (SELECT id FROM publication_types WHERE publication_type = 'story')
GROUP BY p.id, p.time, p.title, p.username, p.score, p.tag, p.tag2, p.web_domain, c.tag
UNION ALL
SELECT p.id, p.time, p.title, p.username, p.score, p.tag, COALESCE(p.tag2, ''), COALESCE(p.web_domain, ''), COALESCE(c.tag, '') AS community_tag, COUNT(b.id) AS badge_count 
FROM publications AS p
INNER JOIN publication_community_tag AS c ON p.id = c.id
LEFT OUTER JOIN publication_badged_by AS b ON p.id = b.id
WHERE p.username = $1 AND p.type_id = (SELECT id FROM publication_types WHERE publication_type = 'story')
GROUP BY p.id, p.time, p.title, p.username, p.score, p.tag, p.tag2, p.web_domain, c.tag
) ORDER BY time DESC LIMIT $2;`

	publications := []DbUserPagePostSharePublication{}
	rows, err := db.Query(q, user, numPostsAndShares)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		d := DbUserPagePostSharePublication{}
		var pubTime int
		err := rows.Scan(&d.Id, &pubTime, &d.Title, &d.User, &d.Score, &d.Tag, &d.Tag2, &d.Domain, &d.CommunityTag, &d.BadgeCount)
		if err != nil {
			return nil, err
		}
		publications = append(publications, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i, publication := range publications {
		children, err := dbGetPublicationChildren(publication.Id, db)
		if err != nil {
			return nil, err
		}
		publications[i].CommentCount = len(children)
	}
	return publications, nil
}
