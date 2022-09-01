package mysql

type Conf struct {
	Addr        string // mysql服务器地址
	Port        string // mysql服务器端口号
	User        string // mysql用户名
	Pwd         string // mysql密码
	Name        string // mysql使用库名
	MaxIdleCoon int    // 最大空闲连接数
	MaxOpenCoon int    // 最大开放连接数
}
