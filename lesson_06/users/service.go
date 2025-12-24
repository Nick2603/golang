package users

import (
	"errors"

	"github.com/Nick2603/golang/lesson_06/documentstore"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

//go:generate mockery --name=CollectionStore --output=mocks --outpkg=mocks

type CollectionStore interface {
	Put(doc documentstore.Document) error
	Get(id string) (*documentstore.Document, error)
	List() []documentstore.Document
	Delete(id string) error
}

type Service struct {
	coll CollectionStore
}

func NewService(coll CollectionStore) *Service {
	return &Service{coll: coll}
}

func (s *Service) CreateUser(id, name string) (*User, error) {
	user := &User{ID: id, Name: name}

	doc, err := documentstore.MarshalDocument(user)
	if err != nil {
		return nil, err
	}

	if err := s.coll.Put(*doc); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) ListUsers() ([]User, error) {
	docs := s.coll.List()
	users := make([]User, 0, len(docs))

	for _, doc := range docs {
		var u User
		if err := documentstore.UnmarshalDocument(&doc, &u); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (s *Service) GetUser(userID string) (*User, error) {
	doc, err := s.coll.Get(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	var user User
	if err := documentstore.UnmarshalDocument(doc, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) DeleteUser(userID string) error {
	if err := s.coll.Delete(userID); err != nil {
		return ErrUserNotFound
	}
	return nil
}
