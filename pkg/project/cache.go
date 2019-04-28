package project

func GetCurrentCacheDir() string {
	if ActiveTarget != nil {
		return ActiveTarget.CacheDir
	}
	return ActiveProject.CacheDir
}
