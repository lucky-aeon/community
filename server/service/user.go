package services

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"xhyovo.cn/community/pkg/cache"
	"xhyovo.cn/community/pkg/constant"
	"xhyovo.cn/community/server/dao"

	"golang.org/x/crypto/bcrypt"
	"xhyovo.cn/community/pkg/email"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/server/service/event"

	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/model"
)

type UserService struct {
}

func (s UserService) DeleteUser(id int) {
	var user model.Users
	model.User().Where("id = ?", id).Find(&user)
	
	// 删除用户的 GitHub 绑定（如果存在）
	var githubAuthService GitHubAuthService
	githubAuthService.UnbindGitHubAccount(id)
	
	// 删除用户
	mysql.GetInstance().Where("id = ?", id).Delete(&model.Users{})
	
	// 删除邀请码
	mysql.GetInstance().Delete(&model.InviteCodes{}, user.InviteCode)
}

func (s UserService) ActiveUsers(page, limit int) (users []*model.UserSimple, count int64) {
	tx := model.User().Select("users.id", "users.name", "users.avatar,users.desc").Joins("left join articles a ON users.id = a.user_id").
		Group("users.id, users.name").
		Having("COUNT(a.id) > 0").
		Order("COUNT(a.id) DESC")
	tx.Count(&count)
	if count == 0 {
		return []*model.UserSimple{}, 0
	}
	tx.Limit(limit).Offset((page - 1) * limit).Find(&users)
	return
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
	user.InviteCode = ""
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

func (*UserService) ResetPwd(account string) bool {
	characters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()")
	length := 10
	password := make([]rune, length)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < length; i++ {
		password[i] = characters[rand.Intn(len(characters))]
	}
	newPwd := string(password)
	newPwdHash, _ := GetPwd(newPwd)
	tx := model.User().Where("account = ?", account).Limit(1).Update("password", newPwdHash)
	if tx.RowsAffected == 0 {
		return false
	}
	email.Send([]string{account}, "您的密码已重置为: "+newPwd, "敲鸭社区:重置密码")
	return true
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

func Register(account, pswd, name, inviteCode string) (int, error) {

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

	// 自动订阅 xhyovo
	var subscription model.Subscriptions
	subscription.SubscriberId = id
	subscription.EventId = event.UserFollowingEvent
	subscription.BusinessId = 13
	var su SubscriptionService
	su.Subscribe(&subscription)

	// 自动订阅会议
	var subscription2 model.Subscriptions
	subscription2.SubscriberId = id
	subscription2.EventId = event.Meeting
	subscription2.BusinessId = constant.MeetingId
	su.Subscribe(&subscription2)

	// 生成账单
	var mS = MemberInfoService{}

	// 查出邀请码
	codeObject := c.GetByCode(inviteCode)
	member := mS.GetById(codeObject.MemberId)
	order := model.Orders{
		InviteCode:      inviteCode,
		Price:           member.Money,
		Purchaser:       id,
		AcquisitionType: codeObject.AcquisitionType,
		Creator:         codeObject.Creator,
	}
	orderDao := dao.OrderDao{}
	orderDao.Save(order)

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

func (s *UserService) PageUsers(p, limit int, condition model.UserSimple, email, inviteCode string) (users []model.Users, count int64) {
	tx := model.User().Where("account like ?", condition.Account)

	if condition.UName != "%%" {
		tx = tx.Where("name like ?", condition.UName)
	}

	if email != "" {
		tx = tx.Where("account like ?", "%"+email+"%")
	}

	if inviteCode != "" {
		tx = tx.Where("invite_code like ?", "%"+inviteCode+"%")
	}

	tx.Count(&count)
	tx.Offset((p - 1) * limit).Limit(limit).Find(&users)
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

	// 先检查是否已经达到限流上限
	v, exists := cache.GetInstance().Get(key)
	if exists {
		count := v.(int)
		if count >= 5 {
			return &model.Users{}, errors.New("操作次数过多,请稍后重试")
		}
	}

	user := userDao.QueryUser(&model.Users{Account: login.Account})

	// 合并登录失败的判断
	if user.ID == 0 || !ComparePswd(user.Password, login.Password) {
		// 登录失败，增加限流计数
		if !cache.CountLimit(key, 5, constant.TTL_LIMIT_lOGIN) {
			return &model.Users{}, errors.New("操作次数过多,请稍后重试")
		}

		// 根据失败原因返回不同的错误信息
		if user.ID == 0 {
			return &model.Users{}, errors.New("登录失败！账号不存在")
		}
		return &model.Users{}, errors.New("登录失败！密码错误")
	}

	// 登录成功，清除限流计数
	cache.GetInstance().Delete(key)
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

func (s *UserService) ListBlackUser(p, limit int) ([]model.Users, int64) {
	var users []model.Users
	var count int64
	tx := model.User().Where("state = 2").Select("id", "name", "account")
	tx.Offset((p - 1) * limit).Limit(limit).Find(&users)
	tx.Count(&count)
	return users, count
}

// 是否被拉黑
func (s *UserService) IsBlack(userId int) bool {
	var count int64
	model.User().Where("id = ? and state = 2", userId).Count(&count)
	return count == 1
}

func (s *UserService) BanByUserId(id int) {
	model.User().Where("id = ?", id).Update("state", 2)
	
	// 可选：禁用用户时同时解绑 GitHub 账号（防止通过 GitHub 登录绕过禁用）
	var githubAuthService GitHubAuthService
	githubAuthService.UnbindGitHubAccount(id)
}

// ban 用户
func (s *UserService) BanByUserAccount(account string) {
	// 先获取用户ID
	var user model.Users
	model.User().Where("account = ?", account).First(&user)
	
	// 更新状态
	model.User().Where("account = ?", account).Update("state", 2)
	
	// 解绑 GitHub 账号
	if user.ID > 0 {
		var githubAuthService GitHubAuthService
		githubAuthService.UnbindGitHubAccount(user.ID)
	}
}

func (s *UserService) UnBanByUserId(id int) {
	model.User().Where("id = ?", id).Update("state", 1)
}

func (s UserService) ExistUserByAccount(account string) bool {
	var count int64
	model.User().Where("account", account).Count(&count)
	return count > 0
}

// CreateUserFromGitHub 通过GitHub信息和邀请码创建用户
func (s *UserService) CreateUserFromGitHub(githubUser *model.GitHubUser, inviteCode string) (*model.Users, error) {
	// 验证邀请码
	var codeService CodeService
	if !codeService.Exist(inviteCode) {
		return nil, errors.New("邀请码不存在或已失效")
	}
	
	// 生成用户账号，优先使用GitHub邮箱，如果没有则生成唯一账号
	account := githubUser.Email
	if account == "" {
		account = fmt.Sprintf("github_%s_%d@github.local", githubUser.Login, githubUser.ID)
	}
	
	// 检查账号是否已存在
	existingUser := userDao.QueryUser(&model.Users{Account: account})
	if existingUser.ID > 0 {
		return nil, errors.New("该邮箱已被注册")
	}
	
	// 生成用户昵称，如果GitHub有name则使用，否则使用login
	name := githubUser.Name
	if name == "" {
		name = githubUser.Login
	}
	
	// 检查昵称是否已存在，如果存在则添加后缀
	originalName := name
	counter := 1
	for {
		existingUser := userDao.QueryUser(&model.Users{Name: name})
		if existingUser.ID == 0 {
			break
		}
		name = fmt.Sprintf("%s_%d", originalName, counter)
		counter++
	}
	
	// 使用邀请码作为密码（方便用户使用邮箱+邀请码登录）
	hashedPassword, err := GetPwd(inviteCode)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}
	
	// 创建用户
	userID := userDao.CreateUser(account, name, string(hashedPassword), inviteCode)
	if userID == 0 {
		return nil, errors.New("用户创建失败")
	}
	
	// 获取创建的用户
	user := userDao.QueryUser(&model.Users{ID: userID})
	
	// 如果有头像URL，更新用户头像
	if githubUser.AvatarURL != "" {
		user.Avatar = githubUser.AvatarURL
		userDao.UpdateUser(user)
	}
	
	// 修改邀请码状态
	codeService.SetState(inviteCode)
	
	// 自动订阅默认事件（复用现有逻辑）
	s.autoSubscribeEvents(userID, inviteCode)
	
	return user, nil
}

// GetGitHubBinding 获取用户的GitHub绑定信息
func (s *UserService) GetGitHubBinding(userID int) (*model.UserGitHubBinding, error) {
	var githubAuthService GitHubAuthService
	return githubAuthService.GetGitHubBinding(userID)
}

// generateRandomPassword 生成随机密码
func (s *UserService) generateRandomPassword() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 16
	
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = charset[rand.Intn(len(charset))]
	}
	return string(bytes)
}

// autoSubscribeEvents 自动订阅默认事件（从Register函数中提取的逻辑）
func (s *UserService) autoSubscribeEvents(userID int, inviteCode string) {
	// 自动订阅 xhyovo
	var subscription model.Subscriptions
	subscription.SubscriberId = userID
	subscription.EventId = event.UserFollowingEvent
	subscription.BusinessId = 13
	var su SubscriptionService
	su.Subscribe(&subscription)

	// 自动订阅会议
	var subscription2 model.Subscriptions
	subscription2.SubscriberId = userID
	subscription2.EventId = event.Meeting
	subscription2.BusinessId = constant.MeetingId
	su.Subscribe(&subscription2)

	// 生成账单
	var mS = MemberInfoService{}
	var cS = CodeService{}
	codeObject := cS.GetByCode(inviteCode)
	member := mS.GetById(codeObject.MemberId)
	order := model.Orders{
		InviteCode:      inviteCode,
		Price:           member.Money,
		Purchaser:       userID,
		AcquisitionType: codeObject.AcquisitionType,
		Creator:         codeObject.Creator,
	}
	orderDao := dao.OrderDao{}
	orderDao.Save(order)
}
