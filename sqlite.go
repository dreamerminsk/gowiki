package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func NewStorage() (*Storage, error) {
	db, err := sql.Open("sqlite3", "nnmclub.sqlite3.db")
	if err != nil {
		return nil, err
	}
	db.Exec("create table if not exists topics (id integer primary key, title text, author text, published text, magnet text, likes integer)")
	s := &Storage{db: db}
	return s, nil
}

func (s *Storage) addTopic(t *Topic) error {
	tx, _ := s.db.Begin()
	stmt, _ := tx.Prepare("insert into topics (id,title,author,published,magnet,likes) values (?,?,?,?,?,?)")
	_, err := stmt.Exec(t.ID, t.Title, t.Author, t.Published, t.Magnet, t.Likes)
	if err != nil {
		tx.Commit()
		return err
	}
	return tx.Commit()
}

func (s *Storage) getTopics() ([]*Topic, error) {
	rows, err := s.db.Query("select * from topics")
	if err != nil {
		return []*Topic{}, err
	}
	topics := []*Topic{}
	for rows.Next() {
		var t Topic
		err = rows.Scan(&t.ID, &t.Title, &t.Author, &t.Published, &t.Magnet, &t.Likes)
		if err != nil {
			return topics, err
		}
		topics = append(topics, &t)
	}
	return topics, nil
}

func (s *Storage) updateTopic(t *Topic) error {
	tx, _ := s.db.Begin()
	stmt, _ := tx.Prepare("update topics set title=?,author=?,published=?,magnet=?,likes=? where id=?")
	_, err := stmt.Exec(t.ID, t.Title, t.Author, t.Published, t.Magnet, t.Likes)
	if err != nil {
		tx.Commit()
		return err
	}
	return tx.Commit()
}

func (s *Storage) deleteTopic(topicId int) error {
	tx, _ := s.db.Begin()
	stmt, _ := tx.Prepare("delete from topics where id=?")
	_, err := stmt.Exec(topicId)
	if err != nil {
		tx.Commit()
		return err
	}
	return tx.Commit()
}

func (s *Storage) Close() {
	s.db.Close()
}
