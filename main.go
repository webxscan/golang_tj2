package main

//golang网站流量统计  中  消息队列+多线程+orm+sql 存库
//QQ:29295842  欢迎技术交流
//http://blog.csdn.net/webxscan
//里面包含了数据库 整个工程GIT有下载
//github  https://github.com/webxscan/golang_tj2
//bee api apiPro -driver=mysql -conn="root:29295842@tcp(127.0.0.1:3306)/seo?charset=utf8"
import (
	"fmt"
	"sync"

	"encoding/json"
	"log"
	//	"strconv"
	"strings"
	"time"

	"net/url"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/session"

	"github.com/Damnever/goqueue"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

//==========================
var (
	Url_queue = goqueue.New(300000)
	wg        = &sync.WaitGroup{}
)

func Queue_size(queue *goqueue.Queue) int { //数量
	return queue.Size()
}

func Queue_get(queue *goqueue.Queue) string { //获取数据
	if !queue.IsEmpty() {
		val, err := queue.Get(5)
		if err != nil {
			return ""
		} else {
			wg.Wait()
			return fmt.Sprintf("%v", val)
		}
	}
	return ""
}

func Queue_put(queue *goqueue.Queue, data string) { //写入数据
	defer wg.Done()
	queue.PutNoWait(data)
	wg.Add(1)
}

//==========================

type Ip struct {
	Id        int    //  `orm:"column(id)"`
	Time      int    `orm:"column(time)" description:"请求时间"`
	Ip        string `orm:"column(ip);size(100)" description:"请求IP"`
	Www_host  string `orm:"column(www_host);size(100)" description:"请求域名"`
	Www_url   string `orm:"column(www_url);size(100);null" description:"请求路径"`
	Referer   string `orm:"column(Referer);size(100);null" description:"来路"`
	Method    string `orm:"column(Method);size(100);null" description:"请求方式"`
	UserAgent string `orm:"column(User_Agent);size(200);null" description:"请求头"`
	Cs        string `orm:"column(cs);size(100);null" description:"请求次数"`
}

func init() {
	maxIdle := 30 //设置最大空闲连接
	maxConn := 30 //设置最大数据库连接
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterModel(new(Ip))
	orm.RegisterDataBase("default", "mysql", "root:316118740@tcp(127.0.0.1:3306)/seo?charset=utf8", maxIdle, maxConn)
}

func add_sql() {
	i := Queue_size(Url_queue)
	if i == 0 {
		return
	}
	for index := 0; index < i; index++ {
		sql_data := Queue_get(Url_queue) //获取数据
		if sql_data == "" {
			continue //跳过
		}
		//================
		//ORM 方法1
		//	orm.Debug = true
		//	o := orm.NewOrm()
		//	o.Using("default") //指定数据库
		//	user := Ip{Time: int(time_Unix), Ip: "ddddddd"}
		//	id, err := o.Insert(&user) //添加
		//	//	o.Delete(&user)            //删除
		//	//	o.Update(&user)            //更新
		//	if err == nil {
		//		fmt.Println(id)
		//	}

		//ORM SQL语句法
		//	orm.Debug = true
		//	o := orm.NewOrm()
		//	o.Using("default") //指定数据库
		//	sSql := fmt.Sprintf("insert into ip (time, ip) values(%d,'%s')", int(time_Unix), "ddddddd")
		//ORM SQL语句法
		orm.Debug = true
		orm_sql := orm.NewOrm()
		orm_sql.Using("default") //指定数据库
		_, err := orm_sql.Raw(sql_data).Exec()
		if err != nil {
			fmt.Println("插入出错")
		}
		//================
	}
	return
}

//===================
var globalSessions *session.Manager

type Iindex struct {
	beego.Controller
}

func init() {
	config := fmt.Sprintf(`{"cookieName":"gosessionid","gclifetime":%d,"enableSetCookie":true}`, 3600*24) //
	conf := new(session.ManagerConfig)
	if err := json.Unmarshal([]byte(config), conf); err != nil {
		log.Fatal("json decode error", err)
	}
	globalSessions, _ = session.NewManager("memory", conf)
	go globalSessions.GC()
}

//网站访问计数器
func (this *Iindex) Count() {
	path_url := this.Ctx.Request.URL.String()
	fmt.Println("get url:", path_url)
	if path_url == "/favicon.ico" { //忽略此路由地址请求
		this.Ctx.WriteString("")
		this.Ctx.ResponseWriter.Header().Set("Content-Type", "text/html")
		return
	}

	//this.Ctx.Request  //这里面有大家所需要一切客户信息
	//fmt.Printf("===%v===\n", this.Ctx.Request)

	time_Unix := time.Now().Unix()
	Client_Host := this.Ctx.Request.Host                           //访问域名
	Client_Method := this.Ctx.Request.Method                       //请求方式
	Client_User_Agent := this.Ctx.Request.Header.Get("User-Agent") //请求头
	Client_IP := this.Ctx.Request.Header.Get("Remote_addr")        //客户端IP
	Client_Referer := this.Ctx.Request.Header.Get("Referer")       //来源
	if len(Client_IP) <= 7 {
		Client_IP = this.Ctx.Request.RemoteAddr //获取客户端IP
	}
	if strings.Contains(Client_IP, ":") {
		ip_boolA, ip_dataA := For_IP(string(Client_IP)) //获取IP
		if ip_boolA {
			Client_IP = ip_dataA
		}
	}

	this.Ctx.ResponseWriter.Header().Set("Content-Type", "text/html")
	this.Ctx.WriteString("golang网站流量统计  中  消息队列+多线程+orm+sql 存库</br>\n")
	this.Ctx.WriteString("QQ:29295842</br>\n")

	this.Ctx.WriteString(fmt.Sprintf("=====客户端IP:%v======</br>\n", Client_IP))
	this.Ctx.WriteString(fmt.Sprintf("=====访问域名:%v======</br>\n", Client_Host))
	this.Ctx.WriteString(fmt.Sprintf("=====请求路径:%v======</br>\n", path_url))
	this.Ctx.WriteString(fmt.Sprintf("=====来源来路:%v======</br>\n", Client_Referer))
	this.Ctx.WriteString(fmt.Sprintf("=====请求方式:%v======</br>\n", Client_Method))
	this.Ctx.WriteString(fmt.Sprintf("=====请求头:%v======</br>\n", Client_User_Agent))
	this.Ctx.WriteString(fmt.Sprintf("=====访问次数:%v======</br>\n", this.Cookie_session()))
	//后面就是数据存贮  可以多种模式
	//消息队列+多线程+orm+sql 存库
	Client_User_Agent = strings.Replace(Client_User_Agent, ",", "_", -1)
	Client_User_Agent = strings.Replace(Client_User_Agent, "，", "_", -1)
	Client_User_Agent = strings.Replace(Client_User_Agent, "'", "_", -1)
	Client_User_Agent = strings.Replace(Client_User_Agent, "\"", "_", -1)
	Client_User_Agent = strings.Replace(Client_User_Agent, "/", "_", -1)
	//Client_User_Agent = strings.Replace(Client_User_Agent, ";", "_", -1)
	//Client_User_Agent = strings.Replace(Client_User_Agent, " ", "_", -1)
	resUri, pErr := url.Parse(path_url) //转义URL地址
	if pErr != nil {
		path_url = "/"
	} else {
		path_url = resUri.Path
	}

	data_sql := fmt.Sprintf("insert into ip (time,ip,www_host,www_url,Referer,Method,User_Agent,cs) values(%d,'%s','%s','%s','%s','%s','%s','%d')",
		int(time_Unix), Client_IP, Client_Host, path_url, Client_Referer, Client_Method, Client_User_Agent, this.Cookie_session())

	Queue_put(Url_queue, data_sql) //写入队列

	return
}

func For_IP(valuex string) (bool, string) {
	data_list := strings.Split(valuex, ":")
	if len(data_list) >= 2 {
		return true, data_list[0]
	}
	return false, ""
}

func (this *Iindex) Cookie_session() int { //id统计  PV  这样统计只能针对单个浏览器有效
	pv := 0
	//=====================
	//Cookie 统计法
	//	cook := this.Ctx.GetCookie("countnum") //获取Cookie
	//	if cook == "" {
	//		this.Ctx.SetCookie("countnum", "1", "/")
	//		pv = 1
	//	} else {
	//		xx, err := strconv.Atoi(cook)
	//		if err == nil {
	//			pv = xx + 1
	//			this.Ctx.SetCookie("countnum", strconv.Itoa(pv), "/")
	//		} else {
	//			//this.Ctx.SetCookie("countnum", "1", "/")
	//			pv = 0
	//		}
	//	}
	//	return pv
	//=====================
	//session 统计法
	sess, _ := globalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)
	ct := sess.Get("countnum")
	if ct == nil {
		sess.Set("countnum", 1)
		pv = 1
	} else {
		pv = ct.(int) + 1
		sess.Set("countnum", pv)
	}
	return pv
}

//===================

func main() {

	fmt.Println("----------------")

	go func() { //多线程任务
		for { //死循环
			time.Sleep(time.Second * 1)
			add_sql() //定时更新器
		}
	}()
	go func() { //
		for { //死循环
			time.Sleep(time.Second * 5)
			add_sql() //定时更新器
		}
	}()
	go func() { //
		for { //死循环
			time.Sleep(time.Second * 10)
			add_sql() //定时更新器
		}
	}()

	//===========================================
	beego.BConfig.Listen.ServerTimeOut = 10 //设置 HTTP 的超时时间，默认是 0，不超时。
	//beego.BConfig.Listen.EnableHTTP = true //是否启用 HTTP 监听，默认是 true。

	beego.BConfig.Listen.HTTPPort = 1000 //应用监听端口，默认为 8080。

	//beego.SetLogger("file", `{"filename":"logs/admin.log","maxlines":10000}`)
	//	beego.BConfig.EnableErrorsShow = false //是否显示系统错误信息，默认为 true。
	//	//是否将错误信息进行渲染，默认值为 true，即出错会提示友好的出错页面，对于 API 类型的应用可能需要将该选项设置为 false 以阻止在 dev 模式下不必要的模板渲染信息返回
	//	beego.BConfig.EnableErrorsRender = false
	//	//运行模式，可选值为 prod, dev 或者 test. 默认是 dev, 为开发模式
	//	//	beego.BConfig.RunMode = "prod"
	//	//运行模式，可选值为 prod, dev 或者 test. 默认是 dev, 为开发模式
	//	beego.BConfig.RunMode = "test"

	beego.BConfig.AppName = "斗转星移"           //应用名称，默认是 beego。通过 bee new 创建的是创建的项目名。
	beego.BConfig.ServerName = "QQ:29295842" //beego 服务器默认在请求的时候输出 server 为 beego。

	beego.BConfig.WebConfig.Session.SessionName = "sessionID"         //存在客户端的 cookie 名称，默认值是 beegosessionID。
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 3600 * 24  //session 过期时间，默认值是 3600 秒。
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime = 3600 * 24 //session 默认存在客户端的 cookie 的时间，默认值是 3600 秒。

	//beego.BConfig.WebConfig.Session.SessionDomain = "" //session cookie 存储域名, 默认空。
	//beego.BConfig.WebConfig.ViewsPath = "admin" //模板路径，默认值是 views。

	beego.Router("/*", &Iindex{}, "*:Count")
	go beego.Run()
	//=============================================
	for { //死循环
		time.Sleep(10 * time.Second)
	}

	//make一个chan用于阻塞主线程,避免程序退出
	//	blockMainRoutine := make(chan bool)
	//	<-blockMainRoutine
	//===========================================

	//	Queue_put(Url_file_queue, "xxxxxxxxxxxxxxxxxxxxxxxxx") //写入数据

	//	fmt.Printf("===%v===\n", Queue_get(Url_file_queue))

}
