package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/util/gmode"
	"golershop.cn/internal/controller/account"
	"golershop.cn/internal/controller/admin"
	"golershop.cn/internal/controller/analytics"
	"golershop.cn/internal/controller/cms"
	"golershop.cn/internal/controller/marketing"
	"golershop.cn/internal/controller/pay"
	"golershop.cn/internal/controller/pt"
	"golershop.cn/internal/controller/shop"
	"golershop.cn/internal/controller/sys"
	"golershop.cn/internal/controller/trade"
	"golershop.cn/internal/service"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start shopsuite http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()

			// HOOK, 开发阶段禁止浏览器缓存,方便调试
			if gmode.IsDevelop() {
				s.BindHookHandler("/*", ghttp.HookBeforeServe, func(r *ghttp.Request) {
					r.Response.Header().Set("Cache-Control", "no-store")
				})
			}

			// 跨域处理
			s.Use(service.Middleware().CORS)

			//需要传播给异步流程或者保持和之前逻辑兼容
			s.Use(service.Middleware().NeverDoneCtx)

			s.Use(service.Middleware().Ctx)

			s.Use(service.Middleware().MiddlewareErrorHandler)
			s.Use(service.Middleware().MiddlewareHandlerResponse)

			s.Use(service.Middleware().CheckLogin)

			// 通过s.Group的分组路由方式定义一组路由注册，在其回调方法中注册的所有路由，都会带有其定义的分组路由前缀/。
			s.Group("/", func(group *ghttp.RouterGroup) {
				//钩子记录异步日志
				group.Hook("/manage/*", ghttp.HookAfterOutput, service.LogAction().OperateLog)
				group.Hook("/front/*", ghttp.HookAfterOutput, service.AccessHistory().OperateAccess)

				/*
					// Group middlewares.
					group.Middleware(
						service.Middleware().Ctx,
						service.Middleware().CheckLogin,
					)

				*/
				// 通过group.Bind方法注册路由对象，该方法将会遍历路由对象的所有公开方法，读取方法的输入输出结构体定义，并对其执行路由注册。
				group.Bind(
					admin.Menu,
					admin.UserAdmin,
					admin.UserRole,
					account.User,
					account.UserInvoice,
					account.UserDeliveryAddress,
					account.DeliveryAddress,
					account.UserLevel,
					account.UserInfo,
					account.UserMessage,
					account.UserTagGroup,
					account.UserBindConnect,

					cms.ArticleBase,
					cms.ArticleTag,
					cms.ArticleCategory,
					cms.ArticleComment,
					cms.Article,
					account.Login,
					marketing.ActivityBase,
					pay.Resource,
					pay.UserResource,
					pay.ConsumeTrade,
					pay.ConsumeDeposit,
					pay.ConsumeRecord,
					pay.ConsumeWithdraw,
					pay.UserPointsHistory,
					pay.PaymentCallback,
					pay.PaymentIndex,
					pay.UserPay,
					pay.Record,
					pt.Product,
					pt.ProductCategory,
					pt.ProductComment,
					pt.ProductTag,
					pt.ProductBrand,
					pt.ProductItem,
					pt.ProductType,
					pt.ProductBase,
					pt.ProductSpec,
					pt.ProductSpecItem,
					pt.ProductAssist,
					shop.Shop,
					shop.FavoritesItem,
					shop.StoreExpressLogistics,
					shop.StoreShippingAddress,
					shop.StoreTransportType,
					shop.StoreTransportItem,
					shop.UserProductBrowse,
					shop.UserVoucher,
					shop.Voucher,
					sys.Config,
					sys.Dict,
					sys.ExpressBase,
					sys.ContractType,
					sys.CrontabBase,
					sys.FeedbackBase,
					sys.Feedback,
					sys.FeedbackCategory,
					sys.FeedbackType,
					sys.LogAction,
					sys.Material,
					sys.MessageTemplate,
					sys.Captcha,
					sys.Upload,
					sys.Page,
					sys.PageModule,
					sys.Release,
					sys.PageCategoryNav,
					sys.PagePcNav,
					trade.Order,
					trade.Cart,
					trade.OrderBase,
					trade.OrderInvoice,
					trade.OrderLogistics,
					trade.OrderReturn,
					trade.Return,
					trade.OrderReturnReason,
					analytics.Analytics,
					analytics.AnalyticsReturn,
				)
			})

			//config init
			service.ConfigBase().Init(ctx)

			s.Run()
			return nil
		},
	}
)
