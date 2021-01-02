package repo

import (
	"errors"
	"fmt"
	"github.com/aaronchen2k/tester/internal/server/biz/domain"
	"github.com/aaronchen2k/tester/internal/server/model"
	"github.com/fatih/color"
	"gorm.io/gorm"
)

type UserRepo struct {
	CommonRepo
	DB *gorm.DB `inject:""`
}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (r *UserRepo) NewUser() *model.User {
	return &model.User{}
}

// GetUser get user
func (r *UserRepo) GetUser(search *domain.Search) (*model.User, error) {
	t := r.NewUser()
	err := r.Found(search).First(t).Error
	if !r.IsNotFound(err) {
		return t, err
	}
	return t, nil
}

// DeleteUser del user . if user's username is username ,can't del it.
func (r *UserRepo) DeleteUser(id uint) error {
	s := &domain.Search{
		Fields: []*domain.Filed{
			{
				Key:       "id",
				Condition: "=",
				Value:     id,
			},
		},
	}
	u, err := r.GetUser(s)
	if err != nil {
		return err
	}
	if u.Username == "username" {
		return errors.New(fmt.Sprintf("不能删除管理员 : %s \n ", u.Username))
	}

	if err := r.DB.Delete(u, id).Error; err != nil {
		color.Red(fmt.Sprintf("DeleteUserByIdErr:%s \n ", err))
		return err
	}
	return nil
}

// GetAllUsers get all users
func (r *UserRepo) GetAllUsers(s *domain.Search) ([]*model.User, int64, error) {
	var users []*model.User
	var count int64
	q := r.GetAll(&model.User{}, s)
	if err := q.Count(&count).Error; err != nil {
		return nil, count, err
	}
	q = q.Scopes(r.Paginate(s.Offset, s.Limit), r.Relation(s.Relations))
	if err := q.Find(&users).Error; err != nil {
		color.Red(fmt.Sprintf("GetAllUserErr:%s \n ", err))
		return nil, count, err
	}
	return users, count, nil
}
