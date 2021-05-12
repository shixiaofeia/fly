package mysql

type Conf struct {
	Address     string // mysql服务器地址
	Port        string // mysql服务器端口号
	User        string // mysql用户名
	Password    string // mysql密码
	DBName      string // mysql使用库名
	MaxIdleCoon int    // 最大空闲连接数
	MaxOpenCoon int    // 最大开放连接数
}
