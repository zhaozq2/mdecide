package main

import (
	"fmt"
	"log"
	"net/http"

	conditionHandler "mdecide/app/condition/handler"
	optionsHandler "mdecide/app/options/handler"
	resultHandler "mdecide/app/result/handler"
	roundHandler "mdecide/app/round/handler"
	topicHandler "mdecide/app/topic/handler"
	"mdecide/common/response"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var defaultTopics = []struct {
	Title       string
	Description string
	Options     []string
}{
	{
		Title:       "职业选择",
		Description: "选择工作时考虑的因素",
		Options:     []string{"工资待遇", "发展空间", "工作氛围", "公司稳定性", "通勤时间", "工作与生活平衡", "行业前景", "团队氛围"},
	},
	{
		Title:       "租房选择",
		Description: "租房时考虑的因素",
		Options:     []string{"租金价格", "地理位置", "房屋面积", "采光通风", "周边配套", "交通便利", "小区安全", "装修状况"},
	},
	{
		Title:       "购车选择",
		Description: "买车时考虑的因素",
		Options:     []string{"预算价格", "油耗表现", "品牌口碑", "空间大小", "安全性能", "外观设计", "维修保养", "保值率"},
	},
	{
		Title:       "教育选择",
		Description: "教育方面的决策",
		Options:     []string{"学校声誉", "学费成本", "专业前景", "地理位置", "师资力量", "教学设施", "就业率", "学制时长"},
	},
	{
		Title:       "理财方式",
		Description: "个人理财投资选择",
		Options:     []string{"收益高低", "风险程度", "流动性", "门槛高低", "安全性", "操作简便", "持有期限", "历史表现"},
	},
	{
		Title:       "健身方式",
		Description: "选择健身方式",
		Options:     []string{"健身房", "居家锻炼", "户外运动", "瑜伽/普拉提", "游泳", "跑步", "球类运动", "私教课程"},
	},
	{
		Title:       "旅游目的地",
		Description: "选择旅游目的地",
		Options:     []string{"预算花费", "景点丰富", "美食体验", "交通便利", "安全程度", "气候舒适", "人流多少", "住宿条件"},
	},
	{
		Title:       "手机选择",
		Description: "购买手机考虑的因素",
		Options:     []string{"预算价格", "品牌偏好", "性能配置", "拍照效果", "续航能力", "屏幕大小", "系统流畅", "外观设计"},
	},
	{
		Title:       "电脑选择",
		Description: "购买电脑考虑的因素",
		Options:     []string{"使用场景", "预算价格", "性能配置", "品牌口碑", "屏幕素质", "续航能力", "便携性", "售后服务"},
	},
	{
		Title:       "留学国家",
		Description: "选择留学国家",
		Options:     []string{"教育质量", "留学费用", "语言环境", "就业前景", "安全程度", "移民政策", "文化体验", "申请难度"},
	},
	{
		Title:       "工作地点",
		Description: "选择工作城市",
		Options:     []string{"工资水平", "生活成本", "发展机会", "气候环境", "教育资源", "医疗条件", "交通便利", "城市氛围"},
	},
	{
		Title:       "装修风格",
		Description: "选择装修风格",
		Options:     []string{"预算成本", "美观程度", "实用功能", "环保材料", "施工难度", "维护保养", "风格持久", "个性化程度"},
	},
	{
		Title:       "婚姻时机",
		Description: "考虑是否结婚",
		Options:     []string{"经济基础", "感情稳定", "年龄合适", "职业发展", "家庭催促", "生子计划", "心理准备", "社会责任"},
	},
	{
		Title:       "副业选择",
		Description: "选择做副业",
		Options:     []string{"投入成本", "时间投入", "技能要求", "收入潜力", "发展空间", "风险程度", "兴趣爱好", "可持续性"},
	},
	{
		Title:       "保险配置",
		Description: "选择保险种类",
		Options:     []string{"保费预算", "保障范围", "赔付额度", "保险公司", "等待期", "续保条件", "理赔速度", "附加服务"},
	},
	{
		Title:       "运动目标",
		Description: "设定运动目标",
		Options:     []string{"减脂塑形", "增肌力量", "提升耐力", "改善柔韧", "释放压力", "社交互动", "培养习惯", "比赛竞技"},
	},
	{
		Title:       "学习方式",
		Description: "选择学习方式",
		Options:     []string{"自学成才", "线上课程", "线下培训", "一对一辅导", "学习效率", "时间灵活", "费用成本", "证书认可"},
	},
	{
		Title:       "餐饮习惯",
		Description: "选择餐饮方式",
		Options:     []string{"健康营养", "口味偏好", "制作难度", "时间成本", "成本控制", "食材获取", "多样化", "社交属性"},
	},
	{
		Title:       "社交方式",
		Description: "选择社交方式",
		Options:     []string{"线上社交", "线下聚会", "兴趣社群", "职场人脉", "时间投入", "情感投入", "效率优先", "深度交流"},
	},
	{
		Title:       "时间管理",
		Description: "时间分配决策",
		Options:     []string{"工作优先", "学习提升", "休息放松", "家庭陪伴", "社交活动", "健康运动", "兴趣爱好", "个人成长"},
	},
}

func main() {
	db, err := gorm.Open(sqlite.Open("mdecide.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	topicHandler.InitDB(db)
	roundHandler.InitRoundDB(db)
	conditionHandler.InitCondDB(db)
	resultHandler.InitResultDB(db)

	r := mux.NewRouter()

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		response.Ok(map[string]string{"status": "ok"}).WriteTo(w)
	})

	optionsHandler.RegisterOptionsRoutes(r)
	conditionHandler.RegisterConditionRoutes(r)
	roundHandler.RegisterRoundRoutes(r)
	resultHandler.RegisterResultRoutes(r)
	topicHandler.RegisterTopicRoutes(r)

	fmt.Println("Starting server on :8888...")
	log.Fatal(http.ListenAndServe(":8888", cors(r)))
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
