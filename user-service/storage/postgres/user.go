package postgres

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	pb "github.com/project1/user-service/genproto"
	
)

type userRepo struct {
	db *sqlx.DB
}

//NewUserRepo ...
func NewUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) LogIn(login *pb.LogInRequest) (*pb.LogInResponse, error) {
	var ruser pb.LogInResponse
	loginQuery := `SELECT id, first_name, last_name, email, bio, status, phone_number, user_name, password FROM users WHERE email = $1`
	err := r.db.QueryRow(loginQuery, login.Email).Scan(
		&ruser.Id,
		&ruser.FirstName,
		&ruser.LastName,
		&ruser.Email,
		&ruser.Bio,
		&ruser.Status,
		pq.Array(&ruser.PhoneNumbers),
		&ruser.UserName,
		&ruser.Password,
	)
	if err != nil {
		return &pb.LogInResponse{}, err
	}
	getByIdAdressQuery := `SELECT city, country, district, postal_code FROM adress WHERE user_id = $1`
	rows, err := r.db.Query(getByIdAdressQuery, ruser.Id)
	if err != nil {
		return nil, err
	}
	var tempUser pb.User
	for rows.Next() {
		var adressById pb.Address
		err = rows.Scan(
			&adressById.City,
			&adressById.Country,
			&adressById.District,
			&adressById.PostalCode,
		)
		if err != nil {
			return nil, err
		}
		tempUser.Address = append(tempUser.Address, &adressById)
	}
	ruser.Address = tempUser.Address

	

	return &ruser, nil
}

func (r *userRepo) CreateUser(user *pb.User) (*pb.User, error) {
	var (
		time_at = time.Now()
		ruser   = pb.User{}
	)

	id, err := uuid.NewV4()
	if err != nil {
		return &pb.User{}, err
	}
	insertUserQuery := `INSERT INTO users (id, first_name, last_name, email, bio, status, created_at, phone_number, user_name, password, refresh_token, access_token) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	RETURNING id, first_name, last_name, email, bio, phone_number, status, user_name, refresh_token, access_token`
	err = r.db.QueryRow(insertUserQuery, id, user.FirstName, user.LastName, user.Email, user.Bio, user.Status, time_at, pq.Array(user.PhoneNumbers), user.UserName, user.Password, user.RefreshToken, user.AccessToken).Scan(
		&ruser.Id,
		&ruser.FirstName,
		&ruser.LastName,
		&ruser.Email,
		&ruser.Bio,
		pq.Array(&ruser.PhoneNumbers),
		&ruser.Status,
		// &ruser.CreatedAt,
		&ruser.UserName,
		&ruser.RefreshToken,
		&ruser.AccessToken,
	)
	if err != nil {
		return &pb.User{}, err
	}
	for _, value := range user.Address {
		var adrsId string
		insertAdressQuery := `INSERT INTO adress (user_id, city, country, district, postal_code)
			VALUES ($1, $2, $3, $4, $5) RETURNING user_id`
		err = r.db.QueryRow(insertAdressQuery, ruser.Id, value.City, value.Country, value.District, value.PostalCode).Scan(&adrsId)
		if err != nil {
			return nil, err
		}
	}
	return &ruser, nil
}

func (r *userRepo) UpdateUser(user *pb.User) (string, error) {
	var (
		timeAt = time.Now()
		// ruser = pb.User{}
	)

	insertUserQuery := `UPDATE users SET first_name = $1, last_name = $2, email = $3, bio = $4, status =$5,
	updated_at = $6, phone_number = $7 where id = $8`
	_, err := r.db.Query(insertUserQuery, user.FirstName, user.LastName, user.Email, user.Bio, user.Status, timeAt, pq.Array(user.PhoneNumbers), user.Id)
	if err != nil {
		return "", err
	}

	for _, value := range user.Address {
		insertAdressQuery := `UPDATE adress SET city = $1, country = $2, district = $3, postal_code = $4 
			WHERE user_id = $5`
		_, err = r.db.Exec(insertAdressQuery, value.City, value.Country, value.District, value.PostalCode, user.Id)
		if err != nil {
			return "", err
		}
	}
	return "", err
}

func (r *userRepo) GetUserById(ID string) (*pb.User, error) {
	var ruser pb.User
	getByIdQuery := `SELECT id, first_name, last_name, email, bio, status, created_at, phone_number FROM users WHERE id = $1`
	err := r.db.QueryRow(getByIdQuery, ID).Scan(
		&ruser.Id,
		&ruser.FirstName,
		&ruser.LastName,
		&ruser.Email,
		&ruser.Bio,
		&ruser.Status,
		&ruser.CreatedAt,
		// &ruser.DeletedAt,
		pq.Array(&ruser.PhoneNumbers),
	)
	if err != nil {
		return &pb.User{}, err
	}
	getByIdAdressQuery := `SELECT city, country, district, postal_code FROM adress WHERE user_id = $1`
	rows, err := r.db.Query(getByIdAdressQuery, ID)
	if err != nil {
		return nil, err
	}
	var tempUser pb.User
	for rows.Next() {
		var adressById pb.Address
		err = rows.Scan(
			&adressById.City,
			&adressById.Country,
			&adressById.District,
			&adressById.PostalCode,
		)
		if err != nil {
			return nil, err
		}
		tempUser.Address = append(tempUser.Address, &adressById)
	}
	ruser.Address = tempUser.Address
	return &ruser, nil
}

func (r *userRepo) GetAllUser() ([]*pb.User, error) {
	var ruser1 []*pb.User

	getByIdQuery := `SELECT id, first_name, last_name, email, bio, status, created_at, phone_number FROM users`
	rowss, err := r.db.Query(getByIdQuery)

	if err != nil {
		return nil, err
	}

	for rowss.Next() {
		var ruser pb.User
		err = rowss.Scan(
			&ruser.Id,
			&ruser.FirstName,
			&ruser.LastName,
			&ruser.Email,
			&ruser.Bio,
			&ruser.Status,
			&ruser.CreatedAt,
			pq.Array(&ruser.PhoneNumbers),
		)
		if err != nil {
			return nil, err
		}

		getByIdAdressQuery := `SELECT city, country, district, postal_code FROM adress`
		rows, err := r.db.Query(getByIdAdressQuery)

		if err != nil {
			return nil, err
		}

		var tempUser pb.User
		for rows.Next() {
			var adressById pb.Address
			err = rows.Scan(
				&adressById.City,
				&adressById.Country,
				&adressById.District,
				&adressById.PostalCode,
			)

			if err != nil {
				return nil, err
			}

			tempUser.Address = append(tempUser.Address, &adressById)
		}
		ruser.Address = tempUser.Address
		ruser1 = append(ruser1, &ruser)
	}

	return ruser1, nil
}

func (r *userRepo) GetUserFromPost(userID string) (*pb.GetUserFromPostResponse, error) {
	var user *pb.GetUserFromPostResponse

	query := `SELECT first_name, last_name FROM users WHERE id = $1`

	err := r.db.QueryRow(query, userID).Scan(
		&user.FirstName,
		&user.LastName,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepo) UserList(limit, page int64) ([]*pb.User, int64, error) {

	var users []*pb.User

	offset := (page - 1) * limit

	listQuery := `SELECT id, first_name, last_name, bio, email, status, created_at, phone_number FROM users LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(listQuery, limit, offset)

	if err != nil {
		return nil, 0, err
	}

	for rows.Next() {
		var user pb.User
		err = rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.Bio,
			&user.Email,
			&user.Status,
			&user.CreatedAt,
			pq.Array(&user.PhoneNumbers),
		)

		if err != nil {
			return nil, 0, err
		}

		var address pb.Address

		addressQuery := `SELECT city, country, district, postal_code FROM adress WHERE user_id = $1`

		rows1, err := r.db.Query(addressQuery, user.Id)

		if err != nil {
			return nil, 0, err
		}

		for rows1.Next() {
			err := rows1.Scan(
				&address.City,
				&address.Country,
				&address.District,
				&address.PostalCode,
			)

			if err != nil {
				return nil, 0, err
			}
		}
		user.Address = append(user.Address, &address)
		users = append(users, &user)
	}

	var count int64
	err = r.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)

	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (r *userRepo) CheckFeild(field, value string) (bool, error) {
	var cle int
	if field == "username" {
		err := r.db.QueryRow("SELECT COUNT(1) FROM users WHERE user_name = $1 AND deleted_at = NULL", value).Scan(&cle)
		if err != nil {
			return false, err
		}
	} else if field == "email" {
		err := r.db.QueryRow("SELECT COUNT(1) FROM users WHERE user_name = $1 AND deleted_at = NULL", value).Scan(&cle)
		if err != nil {
			return false, err
		}
	} else {
		err := errors.New("ERROR IN CheckField")
		return false, err
	}

	if cle == 0 {
		return false, nil
	}

	return true, nil
}
