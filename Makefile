include ../../app_makefile
# docker仓库用户名
DOCKER_USER := fulltimelink
# 当前应用版本
VERSION := 0.0.4
# 服务端口 & grpcui
PORT := 9090
# gorm自动生成模型
GORM_DSN := "root:root@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
GORM_TABLES := "rockscache"