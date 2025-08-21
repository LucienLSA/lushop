package initialize

import (
	"fmt"
	"lushopapi/global"
	"sync"

	v2goodsproto "lushopapi/proto/goods"
	v2inventoryproto "lushopapi/proto/inventory"
	v2orderproto "lushopapi/proto/order"
	v2userproto "lushopapi/proto/user"
	v2useropproto "lushopapi/proto/userop"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

func SrvConn() {
	InitUserSrv()   //用户srv
	InitUserOpSrv() //用户op srv
	InitGoodsSrv()  //商品srv
	InitOrderSrv()  //订单srv
	InitInvSrv()    //库存srv
}

func NewGrpcPool(addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	pool := &myConnPool{
		addr: addr,
		opts: opts,
	}
	return pool.Get()
}

type myConnPool struct {
	addr string
	opts []grpc.DialOption
}

var (
	connCache = make(map[string]*grpc.ClientConn)
	connMutex sync.Mutex
)

func (p *myConnPool) Get() (*grpc.ClientConn, error) {
	connMutex.Lock()
	defer connMutex.Unlock()

	// 检查缓存中是否有可用连接
	if conn, ok := connCache[p.addr]; ok {
		// 检查连接状态
		if conn.GetState() != connectivity.Shutdown && conn.GetState() != connectivity.TransientFailure {
			return conn, nil
		}
		// 关闭无效连接
		conn.Close()
		delete(connCache, p.addr)
	}

	// 创建新连接
	conn, err := grpc.Dial(p.addr, p.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial gRPC server: %v", err)
	}

	// 缓存新连接
	connCache[p.addr] = conn
	return conn, nil
}

func InitUserSrv() {
	consulInfo := global.ServerConfig.ConsulInfo
	addr := fmt.Sprintf("consul://%s:%s/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	}

	conn, err := NewGrpcPool(addr, opts...)
	if err != nil {
		zap.S().Fatal("[InitUserSrv] 连接 【用户服务失败】", zap.Error(err))
		return
	}

	userSrcClient := v2userproto.NewUserClient(conn)
	global.UserSrvClient = userSrcClient
}
func InitUserOpSrv() {
	consulInfo := global.ServerConfig.ConsulInfo
	addr := fmt.Sprintf("consul://%s:%s/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserOpSrvInfo.Name)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	}

	conn, err := NewGrpcPool(addr, opts...)
	if err != nil {
		zap.S().Fatal("[InitUserOpSrv] 连接 【用户op服务失败】", zap.Error(err))
		return
	}

	userOpSrcClient := v2useropproto.NewUserOpClient(conn)
	global.UserOpSrvClient = userOpSrcClient
}
func InitGoodsSrv() {
	consulInfo := global.ServerConfig.ConsulInfo
	addr := fmt.Sprintf("consul://%s:%s/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	}

	conn, err := NewGrpcPool(addr, opts...)
	if err != nil {
		zap.S().Fatal("[InitGoodsSrv] 连接 【商品服务失败】", zap.Error(err))
		return
	}

	goodsSrcClient := v2goodsproto.NewGoodsClient(conn)
	global.GoodsSrvClient = goodsSrcClient
}
func InitOrderSrv() {
	consulInfo := global.ServerConfig.ConsulInfo
	addr := fmt.Sprintf("consul://%s:%s/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.OrderSrvInfo.Name)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	}

	conn, err := NewGrpcPool(addr, opts...)
	if err != nil {
		zap.S().Fatal("[InitOrderSrv] 连接 【订单服务失败】", zap.Error(err))
		return
	}

	orderSrcClient := v2orderproto.NewOrderClient(conn)
	global.OrderSrvClient = orderSrcClient
}
func InitInvSrv() {
	consulInfo := global.ServerConfig.ConsulInfo
	addr := fmt.Sprintf("consul://%s:%s/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.InventorySrvInfo.Name)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	}

	conn, err := NewGrpcPool(addr, opts...)
	if err != nil {
		zap.S().Fatal("[InitInvSrv] 连接 【库存服务失败】", zap.Error(err))
		return
	}

	invSrcClient := v2inventoryproto.NewInventoryClient(conn)
	global.InventorySrvClient = invSrcClient
}

// func SrvConn() {
// 	consulInfo := global.ServerConfig.ConsulInfo
// 	var opts []grpc.DialOption
// 	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()),
// 		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
// 		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
// 	)
// 	retryOpts := []grpc_retry.CallOption{
// 		grpc_retry.WithMax(3),                           // 最大重试次数
// 		grpc_retry.WithPerRetryTimeout(1 * time.Second), // 每次超时最大时间
// 		grpc_retry.WithCodes(codes.Unknown, codes.DeadlineExceeded, codes.Unavailable),
// 	}
// 	opts = append(opts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)))
// 	goodsConn, err := grpc.Dial(
// 		fmt.Sprintf("consul://%s:%s/%s?wait=14s", consulInfo.Host, consulInfo.Port,
// 			global.ServerConfig.GoodsSrvInfo.Name),
// 		opts...,
// 	)
// 	if err != nil {
// 		zap.S().Fatalf("[Init SrvConn] 连接 [商品服务失败]", err.Error())
// 	}
// 	GoodsSrvClient := proto.NewGoodsClient(goodsConn)
// 	global.GoodsSrvClient = GoodsSrvClient
// }

// func SrvConn2() {
// 	// 从注册中心获取到底层groc服务的信息
// 	cfg := api.DefaultConfig()
// 	consulInfo := global.ServerConfig.ConsulInfo
// 	cfg.Address = fmt.Sprintf("%s:%s", consulInfo.Host, consulInfo.Port)
// 	client, err := api.NewClient(cfg)
// 	if err != nil {
// 		panic(err)
// 	}
// 	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`,
// 		global.ServerConfig.UserSrvInfo.Name))
// 	if err != nil {
// 		panic(err)
// 	}
// 	var userSrvHost string
// 	var userSrvPort int
// 	for _, value := range data {
// 		userSrvHost = value.Address
// 		userSrvPort = value.Port
// 		break
// 	}

// 	if userSrvHost == "" {
// 		zap.S().Fatal("[Init SrvConn] 连接 [用户服务失败]")
// 		return
// 	}

// 	// ip := global.ServerConfig.UserSrvInfo.Host
// 	// port := global.ServerConfig.UserSrvInfo.Port
// 	// 拨号连接grpc服务
// 	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
// 	if err != nil {
// 		zap.S().Errorw("[GetUserList] 连接【用户服务失败】",
// 			"msg", err.Error(),
// 		)
// 	}

// 	UserSrvClient := proto.NewUserClient(userConn)
// 	global.UserSrvClient = UserSrvClient
// }
