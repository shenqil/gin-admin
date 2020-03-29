package model

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/model"
	"github.com/LyricTian/gin-admin/internal/app/model/impl/gorm/entity"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/errors"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

var _ model.IUser = new(User)

// UserSet 注入User
var UserSet = wire.NewSet(NewUser, wire.Bind(new(model.IUser), new(*User)))

// NewUser 创建用户存储实例
func NewUser(db *gorm.DB) *User {
	return &User{db}
}

// User 用户存储
type User struct {
	db *gorm.DB
}

func (a *User) getQueryOption(opts ...schema.UserQueryOptions) schema.UserQueryOptions {
	var opt schema.UserQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *User) Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error) {
	db := entity.GetUserDB(ctx, a.db)
	if v := params.UserName; v != "" {
		db = db.Where("user_name=?", v)
	}
	if v := params.LikeUserName; v != "" {
		db = db.Where("user_name LIKE ?", "%"+v+"%")
	}
	if v := params.LikeRealName; v != "" {
		db = db.Where("real_name LIKE ?", "%"+v+"%")
	}
	if v := params.Status; v > 0 {
		db = db.Where("status=?", v)
	}
	if v := params.RoleIDs; len(v) > 0 {
		subQuery := entity.GetUserRoleDB(ctx, a.db).Select("user_id").Where("role_id IN(?)", v).SubQuery()
		db = db.Where("record_id IN ?", subQuery)
	}
	db = db.Order("id DESC")

	opt := a.getQueryOption(opts...)
	var list entity.Users
	pr, err := WrapPageQuery(ctx, db, opt.PageParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &schema.UserQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaUsers(),
	}
	return qr, nil
}

// Get 查询指定数据
func (a *User) Get(ctx context.Context, recordID string, opts ...schema.UserQueryOptions) (*schema.User, error) {
	var item entity.User
	ok, err := FindOne(ctx, entity.GetUserDB(ctx, a.db).Where("record_id=?", recordID), &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaUser(), nil
}

// Create 创建数据
func (a *User) Create(ctx context.Context, item schema.User) error {
	sitem := entity.SchemaUser(item)
	result := entity.GetUserDB(ctx, a.db).Create(sitem.ToUser())
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update 更新数据
func (a *User) Update(ctx context.Context, recordID string, item schema.User) error {
	eitem := entity.SchemaUser(item).ToUser()
	result := entity.GetUserDB(ctx, a.db).Where("record_id=?", recordID).Omit("record_id", "creator").Updates(eitem)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete 删除数据
func (a *User) Delete(ctx context.Context, recordID string) error {
	result := entity.GetUserDB(ctx, a.db).Where("record_id=?", recordID).Delete(entity.User{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdateStatus 更新状态
func (a *User) UpdateStatus(ctx context.Context, recordID string, status int) error {
	result := entity.GetUserDB(ctx, a.db).Where("record_id=?", recordID).Update("status", status)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdatePassword 更新密码
func (a *User) UpdatePassword(ctx context.Context, recordID, password string) error {
	result := entity.GetUserDB(ctx, a.db).Where("record_id=?", recordID).Update("password", password)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}