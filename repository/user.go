package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	uuid "github.com/satori/go.uuid"
	"golang-echo-elasticsearch-rest-api-example/exception"
	"golang-echo-elasticsearch-rest-api-example/model"
	"golang-echo-elasticsearch-rest-api-example/util"
	"reflect"
)

var ctx context.Context = context.Background()

const (
	UserIndexName = "users_index"
)

type UserRepository interface {
	Count() int64
	GetAllUser(page int64, limit int64) (*util.PagedModel, error)
	SaveUser(user *model.User) (*model.User, error)
	GetUser(id string) (*model.User, error)
	UpdateUser(id string, user *model.User) (*model.User, error)
	DeleteUser(id string) error
}

type userRepositoryImpl struct {
	client *elastic.Client
}

func NewUserRepository(client *elastic.Client) UserRepository {
	return &userRepositoryImpl{client: client}
}

func (userRepository *userRepositoryImpl) Count() int64 {
	count, _ := userRepository.client.Count(UserIndexName).Do(ctx)
	return count
}

func (userRepository *userRepositoryImpl) GetAllUser(page int64, limit int64) (*util.PagedModel, error) {
	cc := make(chan int64, 0)

	go countRecords(userRepository.client, cc)
	count := <-cc

	paginator := util.Paging(page, limit, count)

	result, err := userRepository.client.Search(UserIndexName).Size(int(paginator.Limit)).From(int(paginator.Offset)).Do(ctx)
	if err != nil {
		return nil, err
	}

	var existingUser model.User

	users := result.Each(reflect.TypeOf(&existingUser))

	return paginator.PagedData(users), nil
}

func (userRepository *userRepositoryImpl) SaveUser(user *model.User) (*model.User, error) {
	id := uuid.NewV4().String()
	user.ID = id
	_, err := userRepository.client.Index().Index(UserIndexName).Id(user.ID).BodyJson(&user).Do(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userRepository *userRepositoryImpl) GetUser(id string) (*model.User, error) {
	var existingUser model.User

	result, err := userRepository.client.Get().Index(UserIndexName).Id(id).Do(ctx)
	if err != nil {
		switch {
		case elastic.IsNotFound(err):
			return nil, exception.ResourceNotFoundException("User", "id", id)
		default:
			return nil, err
		}
	}

	err = json.Unmarshal(result.Source, &existingUser)
	if err != nil {
		fmt.Println(err)
	}

	return &existingUser, nil
}

func (userRepository *userRepositoryImpl) UpdateUser(id string, user *model.User) (*model.User, error) {
	_, err := userRepository.client.Update().Index(UserIndexName).Id(id).Doc(&user.UserInput).Do(ctx)
	if err != nil {
		switch {
		case elastic.IsNotFound(err):
			return nil, exception.ResourceNotFoundException("User", "id", id)
		default:
			return nil, err
		}
	}

	user.ID = id
	return user, nil
}

func (userRepository *userRepositoryImpl) DeleteUser(id string) error {

	_, err := userRepository.client.Delete().Index(UserIndexName).Id(id).Do(ctx)
	if err != nil {
		switch {
		case elastic.IsNotFound(err):
			return exception.ResourceNotFoundException("User", "id", id)
		default:
			return err
		}
	}

	return nil
}

func countRecords(client *elastic.Client, cc chan int64) {
	count, _ := client.Count(UserIndexName).Do(ctx)
	cc <- count
}
