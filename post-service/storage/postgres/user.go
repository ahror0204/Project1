package postgres

import (
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	pb "github.com/project1/post-service/genproto"
)

type postRepo struct {
	db *sqlx.DB
}

//NewUserRepo ...
func NewPostRepo(db *sqlx.DB) *postRepo {
	return &postRepo{db: db}
}

func (r *postRepo) CreatePost(post *pb.Post) (*pb.Post, error) {
	var (
		rPost = pb.Post{}
	)
	id, err := uuid.NewV4()
	err = r.db.QueryRow("INSERT INTO posts (id, name, description, user_id) VALUES($1, $2, $3, $4) RETURNING id, name, description, user_id", id, post.Name, post.Description, post.UserId).Scan(
		&rPost.Id,
		&rPost.Name,
		&rPost.Description,
		&rPost.UserId,
	)
	if err != nil {
		return &pb.Post{}, err
	}

	for _, media := range post.Medias {
		id, err := uuid.NewV4()
		if err != nil {
			return nil, err
		}
		_, err = r.db.Exec("INSERT INTO medias (id, type, link, post_id) VALUES($1, $2, $3, $4)", id, media.Type, media.Link, rPost.Id)
		if err != nil {
			return &pb.Post{}, err
		}
	}

	return &rPost, nil
}

func (r *postRepo) GetPostById(ID string) (*pb.Post, error) {
	var (
		rPost = pb.Post{}
	)

	err := r.db.QueryRow("SELECT id, name, description, user_id from posts WHERE id = $1", ID).Scan(
		&rPost.Id,
		&rPost.Name,
		&rPost.Description,
		&rPost.UserId,
	)

	if err != nil {
		return nil, err
	}

	return &rPost, nil
}

func (r *postRepo) GetAllUserPosts(userID string) ([]*pb.Post, error) {
	var (
		posts []*pb.Post
	)

	rows, err := r.db.Query("SELECT id, name, description, user_id from posts WHERE user_id = $1", userID)

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var post pb.Post
		err := rows.Scan(
			&post.Id,
			&post.Name,
			&post.Description,
			&post.UserId,
		)
		if err != nil {
			return nil, err
		}

		var medias []*pb.Media
		rows, err := r.db.Query("SELECT id, type, link from medias WHERE post_id = $1", post.Id)

		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var media pb.Media
			err := rows.Scan(
				&media.Id,
				&media.Type,
				&media.Link,
			)
			if err != nil {
				return nil, err
			}

			post.Medias = append(medias, &media)
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *postRepo) GetUserByPostId(postID string) (*pb.GetUserByPostIdResponse, error) {
	var post pb.GetUserByPostIdResponse

	query := `SELECT id, user_id, name, description FROM posts WHERE id = $1`

	err := r.db.QueryRow(query, postID).Scan(
		&post.PostId,
		&post.UserId,
		&post.Name,
		&post.Description,
	)
	if err != nil {
		return nil, err
	}
	mediaQuery := `SELECT type, link FROM medias WHERE post_id = $1`

	rows, err := r.db.Query(mediaQuery, postID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var media pb.Media
		rows.Scan(
			&media.Type,
			&media.Link,
		)
		post.Medias = append(post.Medias, &media)
	}
	return &post, nil
}

func (r *postRepo) CreatePostUser(user *pb.User) error {
	_, err := r.db.Query("INSERT INTO post_users (id, first_name, last_name) VALUES ($1,$2,$3)", user.Id, user.FirstName, user.LastName)

	if err != nil {
		return err
	}

	return nil
}
