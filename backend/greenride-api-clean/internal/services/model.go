package services

func SetupService() {
	GetFirebaseService()
	// 初始化用户任务处理器
	InitUserTaskHandlers()
	InitPaymentChannelHandlers()
	InitOrderTaskHandlers()
}
