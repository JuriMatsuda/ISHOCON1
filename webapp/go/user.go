package main

import (
	"log"
	"time"

	"github.com/gin-gonic/contrib/sessions"
)

// User model
type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	LastLogin string
}

func authenticate(email string, password string) (uid int, result bool) {
	var dbPass string
	err := db().QueryRow("SELECT id, password FROM users WHERE email = ? LIMIT 1", email).Scan(&uid, &dbPass)
	if err != nil {
		return 0, false
	}
	result = password == dbPass
	return
}

func notAuthenticated(session sessions.Session) bool {
	uid := session.Get("uid")
	return !(uid.(int) > 0)
}

func getUser(uid int) User {
	u := User{}
	r := db().QueryRow("SELECT * FROM users WHERE id = ? LIMIT 1", uid)
	err := r.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.LastLogin)
	if err != nil {
		return u
	}

	return u
}

func currentUser(session sessions.Session) User {
	uid := session.Get("uid")
	u := User{}
	r := db().QueryRow("SELECT * FROM users WHERE id = ? LIMIT 1", uid)
	err := r.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.LastLogin)
	if err != nil {
		return u
	}

	return u
}

// BuyingHistory : products which user had bought
func (u *User) BuyingHistory() (products []Product) {
	rows, err := db().Query(
		"SELECT p.id, p.name, p.description, p.image_path, p.price, h.created_at "+
			"FROM histories as h "+
			"LEFT OUTER JOIN products as p "+
			"ON h.product_id = p.id "+
			"WHERE h.user_id = ? "+
			"ORDER BY h.id DESC", u.ID)
	if err != nil {
		return nil
	}

	defer rows.Close()
	for rows.Next() {
		p := Product{}
		var cAt string
		err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.ImagePath, &p.Price, &cAt)
		if err != nil {
			panic(err.Error())
		}
		log.Print(cAt)
		var tmp time.Time
		tmp, err = time.Parse("2006-01-02 15:04:05", cAt)
		tmp = tmp.Add(9 * time.Hour)
		// TODO: +0900 がつかない
		p.CreatedAt = tmp.Format("2006-01-02 15:04:05 -0700")
		if err != nil {
			panic(err.Error())
		}
		products = append(products, p)
	}

	return
}

// BuyProduct : buy product
func (u *User) BuyProduct(pid string) {
	db().Exec(
		"INSERT INTO histories (product_id, user_id, created_at) VALUES (?, ?, ?)",
		pid, u.ID, time.Now())
}

// CreateComment : create comment to the product
func (u *User) CreateComment(pid string, content string) {
	db().Exec(
		"INSERT INTO comments (product_id, user_id, content, created_at) VALUES (?, ?, ?, ?)",
		pid, u.ID, content, time.Now())
}
