package common

import (
	"context"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

const (
	//ContextKey 上下文变量存储键名,前后端系统共享
	ContextKey = "ContextKey"
)

//ContextUser 请求上下文中的用户信息
type ContextUser struct {
	Id       uint64 //用户ID
	Account  string //用户账号
	NickName string //用户名称
	Avatar   string //用户头像
}

//Context 请求上下文结构
type Context struct {
	Session *ghttp.Session //当前Session管理对象
	User    *ContextUser   //上下文用户信息
	Data    g.Map          //自定KV变量,业务模块根据需要设置,不固定
}

type contextService struct{}

//Init 初始化上下文对象中,以便后续的请求流程中可以修改
func (s *contextService) Init(r *ghttp.Request, ctx *Context) {
	r.SetCtxVar(ContextKey, ctx)
}

func (s *contextService) Middleware(r *ghttp.Request, user *ContextUser) {
	//初始化,务必最开始执行
	ctx := &Context{
		Session: r.Session,
		Data:    make(g.Map),
		User:    user,
	}
	s.Init(r, ctx)
	//将自定义的上下文对象传递到模板变量中使用
	r.Assigns(g.Map{
		"Context": ctx,
	})
	//执行下一步请求逻辑
	r.Middleware.Next()
}

//Get 获取上下文变量,如果没有设置,那么返回nil
func (s *contextService) Get(ctx context.Context) *Context {
	val := ctx.Value(ContextKey)
	if val == nil {
		return nil
	}
	if localCtx, ok := val.(*Context); ok {
		return localCtx
	}
	return nil
}

//SetUser 将上下文信息设置到上下文请求中,注意是完整覆盖
func (s *contextService) SetUser(ctx context.Context, user *ContextUser) {
	s.Get(ctx).User = user
}

//GetUser 根据上下文获取当前用户
func (s *contextService) GetUser(ctx context.Context) *ContextUser {
	if cu := s.Get(ctx); cu != nil {
		return cu.User
	}
	return nil
}
