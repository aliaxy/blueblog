package mysql

import "blueblog/models"

func CreatePost(post *models.Post) (err error) {
	sqlStr := `insert into
		post (post_id, title, content, author_id, community_id)
		values (?, ?, ?, ?, ?)`

	_, err = db.Exec(sqlStr, post.ID, post.Title, post.Content, post.AuthorID, post.CommunityID)
	return
}

func GetPostByID(pid int64) (post *models.Post, err error) {
	sqlStr := `select
		post_id, title, content, author_id, community_id, create_time
		from post
		where post_id = ?`
	post = new(models.Post)
	err = db.Get(post, sqlStr, pid)
	return
}
