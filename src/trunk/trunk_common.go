package trunk

func (self Trunk) Test() {
	println("Test")
}

func (self Trunk) GetRootPath() string {
	return self.RootPath
}
