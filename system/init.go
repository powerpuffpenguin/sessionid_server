package system

func Init() {
	initProvider()
	manager.init(provider)
}
