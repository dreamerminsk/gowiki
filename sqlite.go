package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

const (
	topicCreateSQL = `create table if not exists topics (id integer primary key, title text, author text, published text, magnet text, likes integer)`
	topicSelectSQL = `select * from topics`
	topicInsertSQL = `insert into topics (id,title,author,published,magnet,likes) values (?,?,?,?,?,?)`
	topicUpdateSQL = `update topics set title=?,author=?,published=?,magnet=?,likes=? where id=?`
	topicDeleteSQL = `"delete from topics where id=?"`
)

func NewStorage() (*Storage, error) {
	db, err := sql.Open("sqlite3", "nnmclub.sqlite3.db")
	if err != nil {
		return nil, err
	}
	db.Exec(topicCreateSQL)
	s := &Storage{db: db}
	return s, nil
}

func (s *Storage) addTopic(t *Topic) error {
	tx, _ := s.db.Begin()
	stmt, _ := tx.Prepare(topicInsertSQL)
	_, err := stmt.Exec(t.ID, t.Title, t.Author, t.Published, t.Magnet, t.Likes)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *Storage) getTopics() ([]*Topic, error) {
	rows, err := s.db.Query(topicSelectSQL)
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
	stmt, _ := tx.Prepare(topicUpdateSQL)
	_, err := stmt.Exec(t.ID, t.Title, t.Author, t.Published, t.Magnet, t.Likes)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *Storage) deleteTopic(topicId int) error {
	tx, _ := s.db.Begin()
	stmt, _ := tx.Prepare(topicDeleteSQL)
	_, err := stmt.Exec(topicId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *Storage) Close() {
	s.db.Close()
}
