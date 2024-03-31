package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"

	"xhyovo.cn/community/pkg/mysql"

	"xhyovo.cn/community/pkg/cache"
	"xhyovo.cn/community/pkg/constant"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/model"
)

type UserService struct {
}

func (s UserService) IsAdmin(userId int) (bool, error) {

	query := mysql.GetInstance().Table("users as u").Select("m.name").
		Joins("join invite_codes as inv on u.invite_code = inv.code").
		Joins("join member_infos as m on m.id = inv.member_id").
		Where("u.id = ? and u.deleted_at is NULL", userId)
	rows, err := query.Rows()
	if err != nil {
		return false, err
	}
	defer rows.Close()
	var name string
	if rows.Next() {
		rows.Scan(&name)
		if name == "admin" {
			return true, nil
		}
	}
	defer rows.Close()
	return false, nil
}

func (s UserService) ListUsers(name string) (users []model.Users) {
	model.User().Where("name like ?", "%"+name+"%").Select("name", "id").Limit(10).Find(&users)
	return
}

// get user information
func (*UserService) GetUserById(id int) *model.Users {

	user := userDao.QueryUser(&model.Users{ID: id})
	user.InviteCode = 0
	return user
}

// get user information
func (*UserService) GetUserSimpleById(id int) *model.UserSimple {
	user, _ := userDao.QueryUserSimple(&model.Users{ID: id})

	return &user
}

// update user information
func (*UserService) UpdateUser(user *model.Users) {

	userDao.UpdateUser(user)

}

func (*UserService) ListByIdsSelectEmail(id ...int) []string {
	return userDao.ListByIds(id...)
}

func (s *UserService) ListByIdsToMap(ids []int) map[int]model.Users {

	m := make(map[int]model.Users)
	users := userDao.ListByIdsSelectIdName(ids)
	for i := range users {
		user := users[i]
		user.Password = ""
		m[user.ID] = user
	}
	return m
}

func Register(account, pswd, name string, inviteCode int) (int, error) {

	if err := utils.NotBlank(account, pswd, name, inviteCode); err != nil {
		return 0, err
	}

	// query codeDao
	if !codeDao.Exist(inviteCode) {
		return 0, errors.New("验证码不存在")
	}

	// 查询账户
	user := userDao.QueryUser(&model.Users{Account: account})
	if user.ID > 0 {
		return 0, errors.New("账户已存在,换一个吧")
	}

	user = userDao.QueryUser(&model.Users{Name: name})
	if user.ID > 0 {
		return 0, errors.New("用户昵称已存在,换一个吧")
	}
	pwd, err := GetPwd(pswd)
	if err != nil {
		return 0, err
	}
	// 保存用户
	id := userDao.CreateUser(account, name, string(pwd), inviteCode)
	// 修改code状态
	var c CodeService
	c.SetState(inviteCode)

	return id, nil
}

type UserMenu struct {
	Path       string                 `json:"path"`
	Redirect   string                 `json:"redirect,omitempty"`
	Name       string                 `json:"name"`
	Components string                 `json:"components,omitempty"`
	Meta       map[string]interface{} `json:"meta,omitempty"`
	Children   []UserMenu             `json:"children,omitempty"`
}

func (t *UserService) GetUserMenu() []*UserMenu {
	rootMenu := typeDao.List(0)
	parentIds := make([]int, len(rootMenu))
	userMenu := make(map[int]*UserMenu)
	for i, item := range rootMenu {
		parentIds[i] = item.ID
		path := "/article/"
		if item.ID == 1 {
			path = "/qa/"
		}
		userMenu[int(item.ID)] = &UserMenu{
			Path:     path,
			Name:     item.FlagName,
			Children: []UserMenu{},
			Meta: map[string]interface{}{
				"locale":       item.Title,
				"requiresAuth": true,
				"icon":         "icon-dashboard",
				"order":        1,
				"id":           item.ID,
			},
		}
	}
	children := []model.Types{}
	model.Type().Where("parent_id in (?)", parentIds).Find(&children)
	for _, item := range children {
		um := userMenu[int(item.ParentId)]
		um.Children = append(um.Children, UserMenu{
			Path: um.Path + item.FlagName,
			Name: item.FlagName,
			Meta: map[string]interface{}{
				"locale":       item.Title,
				"requiresAuth": true,
				"icon":         "icon-dashboard",
				"id":           item.ID,
			},
		})
	}
	result := []*UserMenu{}

	for _, um := range userMenu {
		result = append(result, um)
	}
	return result
}

func (s *UserService) CheckCodeUsed(code int) bool {
	var count int64
	model.User().Where("invite_code = ?", code).Count(&count)
	return count == 1
}

func (s *UserService) PageUsers(p, limit int) (users []model.Users, count int64) {
	model.User().Offset((p - 1) * limit).Limit(limit).Find(&users)
	model.User().Count(&count)
	return users, count
}

func (s *UserService) ListByNameSelectEmailAndId(usernames []string) (emails []string, id []int) {
	var users []model.Users
	model.User().Where("name in ? ", usernames).Select("account", "id").Find(&users)

	for i := range users {
		u := users[i]
		emails = append(emails, u.Account)
		id = append(id, u.ID)
	}

	return emails, id
}

func (s *UserService) ListBySelect(user model.Users) (users []model.Users) {
	model.User().Where(user).Find(&users)
	return
}

func (s *UserService) Statistics(userId, types int) (m map[string]interface{}) {
	m = make(map[string]interface{})
	var articleS ArticleService
	// 获取被点赞次数,获取用户发布的所有文章
	var articleCount int64
	if types == 1 {
		articleCount = articleS.PublishArticleCount(userId)
	} else {
		articleCount = articleS.QAArticleCount(userId)
	}
	ids := articleS.PublishArticlesSelectId(userId)
	likeCount := articleS.ArticlesLikeCount(ids)
	// 获取发布文章

	m["articleCount"] = articleCount
	m["likeCount"] = likeCount
	return
}

func (s *UserService) SearchNameSelectId(name string) (ids []int) {

	model.User().Where("name like ?", "%"+name+"%").Select("id").Find(&ids)
	return
}

func Login(login model.LoginForm) (*model.Users, error) {
	key := constant.LIMIT_LOGIN + login.Account
	if !cache.CountLimit(key, 5, constant.TTL_LIMIT_lOGIN) {
		return &model.Users{}, errors.New("操作次数过多,请稍后重试")
	}
	user := userDao.QueryUser(&model.Users{Account: login.Account})
	if user.ID == 0 {
		return &model.Users{}, errors.New("登录失败！账号不存在")
	}
	if !ComparePswd(user.Password, login.Password) {
		return &model.Users{}, errors.New("登录失败！密码错误")

	}
	return user, nil
}

func ComparePswd(oldPwsd, newPswd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(oldPwsd), []byte(newPswd))
	if err != nil {
		return false
	} else {
		return true
	}
}
func GetPwd(pwd string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return hash, err
}
