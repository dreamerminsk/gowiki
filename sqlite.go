package main

import (
	"database/sql"

	"github.com/dreamerminsk/gowiki/model"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteStorage struct {
	db *sql.DB
}

const (
	topicCreateSQL    = `create table if not exists topics (id integer primary key, title text, author text, published datetime, magnet text, likes integer)`
	topicSelectAllSQL = `select * from topics`
	topicSelectOneSQL = `select * from topics where id=?`
	topicInsertSQL    = `insert into topics (id,title,author,published,magnet,likes) values (?,?,?,?,?,?)`
	topicUpdateSQL    = `update topics set title=?,author=?,published=?,magnet=?,likes=? where id=?`
	topicDeleteSQL    = `"delete from topics where id=?"`
)

func NewStorage() (*SqliteStorage, error) {
	db, err := sql.Open("sqlite3", "nnmclub.sqlite3.db")
	if err != nil {
		return nil, err
	}
	db.Exec(topicCreateSQL)
	s := &SqliteStorage{db: db}
	return s, nil
}

func (s *SqliteStorage) getTopics() ([]*model.Topic, error) {
	rows, err := s.db.Query(topicSelectAllSQL)
	if err != nil {
		return []*model.Topic{}, err
	}
	topics := []*model.Topic{}
	for rows.Next() {
		var t model.Topic
		err = rows.Scan(&t.ID, &t.Title, &t.Author, &t.Published, &t.Magnet, &t.Likes)
		if err != nil {
			return topics, err
		}
		topics = append(topics, &t)
	}
	return topics, nil
}

func (s *SqliteStorage) getTopic(id int) (*model.Topic, error) {
	t := &model.Topic{}
	err := s.db.QueryRow(topicSelectOneSQL, id).Scan(&t.ID, &t.Title, &t.Author, &t.Published, &t.Magnet, &t.Likes)
	switch {
	case err == sql.ErrNoRows:
		return t, nil
	case err != nil:
		return nil, err
	default:
		return t, nil
	}
}

func (s *SqliteStorage) addTopic(t *model.Topic) error {
	tx, _ := s.db.Begin()
	stmt, _ := tx.Prepare(topicInsertSQL)
	_, err := stmt.Exec(t.ID, t.Title, t.Author, t.Published, t.Magnet, t.Likes)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *SqliteStorage) updateTopic(t *model.Topic) error {
	tx, _ := s.db.Begin()
	stmt, _ := tx.Prepare(topicUpdateSQL)
	_, err := stmt.Exec(t.Title, t.Author, t.Published, t.Magnet, t.Likes, t.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *SqliteStorage) deleteTopic(topicId int) error {
	tx, _ := s.db.Begin()
	stmt, _ := tx.Prepare(topicDeleteSQL)
	_, err := stmt.Exec(topicId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *SqliteStorage) Close() {
	s.db.Close()
}
