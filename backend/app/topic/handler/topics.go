package handler

import (
	"encoding/json"
	"net/http"

	topicModel "mdecide/app/topic/model"
	"mdecide/common/response"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(database *gorm.DB) {
	db = database
	db.AutoMigrate(&topicModel.Topic{}, &topicModel.Option{}, &topicModel.Round{})
}

func RegisterTopicRoutes(r *mux.Router) {
	r.HandleFunc("/api/topics/seed", SeedTopics).Methods("POST")
	r.HandleFunc("/api/topics/templates", ListTemplateTopics).Methods("GET")
	r.HandleFunc("/api/topics/templates/{id}/options", GetTemplateOptions).Methods("GET")
	r.HandleFunc("/api/topics/import", ImportTemplates).Methods("POST")
	r.HandleFunc("/api/topics/{id}/save-as-template", SaveAsTemplate).Methods("POST")
	r.HandleFunc("/api/topics", ListTopics).Methods("GET")
	r.HandleFunc("/api/topics", CreateTopic).Methods("POST")
	r.HandleFunc("/api/topics/{id}", GetTopic).Methods("GET")
	r.HandleFunc("/api/topics/{id}", UpdateTopic).Methods("PUT")
	r.HandleFunc("/api/topics/{id}", DeleteTopic).Methods("DELETE")
	r.HandleFunc("/api/topics/{id}/options", ListOptions).Methods("GET")
	r.HandleFunc("/api/topics/{id}/options", SaveOptions).Methods("POST")
	r.HandleFunc("/api/topics/{id}/rounds", ListRounds).Methods("GET")
	r.HandleFunc("/api/topics/{id}/rounds", SaveRounds).Methods("POST")
}

func ListTopics(w http.ResponseWriter, r *http.Request) {
	var topics []topicModel.Topic
	db.Find(&topics)
	response.Ok(topics).WriteTo(w)
}

func CreateTopic(w http.ResponseWriter, r *http.Request) {
	var topic topicModel.Topic
	if err := json.NewDecoder(r.Body).Decode(&topic); err != nil {
		response.Error(http.StatusBadRequest, "invalid request body").WriteTo(w)
		return
	}
	db.Create(&topic)
	response.Ok(topic).WriteTo(w)
}

func GetTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var topic topicModel.Topic
	if err := db.First(&topic, vars["id"]).Error; err != nil {
		response.Error(http.StatusNotFound, "topic not found").WriteTo(w)
		return
	}
	response.Ok(topic).WriteTo(w)
}

func UpdateTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var topic topicModel.Topic
	if err := json.NewDecoder(r.Body).Decode(&topic); err != nil {
		response.Error(http.StatusBadRequest, "invalid request body").WriteTo(w)
		return
	}
	db.Model(&topicModel.Topic{}).Where("id = ?", vars["id"]).Updates(topic)
	response.Ok(topic).WriteTo(w)
}

func DeleteTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	db.Delete(&topicModel.Topic{}, vars["id"])
	response.OkMsg("deleted").WriteTo(w)
}

func ListOptions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var options []topicModel.Option
	db.Where("topic_id = ?", vars["id"]).Find(&options)
	response.Ok(options).WriteTo(w)
}

func SaveOptions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var options []topicModel.Option
	if err := json.NewDecoder(r.Body).Decode(&options); err != nil {
		response.Error(http.StatusBadRequest, "invalid request body").WriteTo(w)
		return
	}
	db.Where("topic_id = ?", vars["id"]).Delete(&topicModel.Option{})
	topicID := parseUint(vars["id"])
	for i := range options {
		options[i].TopicID = topicID
	}
	db.Create(&options)
	response.Ok(options).WriteTo(w)
}

func ListRounds(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var rounds []topicModel.Round
	topicID := parseUint(vars["id"])
	db.Where("topic_id = ?", topicID).Find(&rounds)
	response.Ok(rounds).WriteTo(w)
}

func SaveRounds(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var reqRounds []map[string]any
	if err := json.NewDecoder(r.Body).Decode(&reqRounds); err != nil {
		response.Error(http.StatusBadRequest, "invalid request body").WriteTo(w)
		return
	}
	topicID := parseUint(vars["id"])
	db.Where("topic_id = ?", topicID).Delete(&topicModel.Round{})

	for _, reqR := range reqRounds {
		round := topicModel.Round{
			TopicID: topicID,
		}
		if v, ok := reqR["roundNumber"].(float64); ok {
			round.RoundNumber = int(v)
		}
		if v, ok := reqR["importanceStatus"].(string); ok {
			round.ImportanceStatus = v
		}
		if v, ok := reqR["necessityStatus"].(string); ok {
			round.NecessityStatus = v
		}
		if v, ok := reqR["status"].(string); ok {
			round.Status = v
		}
		if v, ok := reqR["matchScore"].(float64); ok {
			round.MatchScore = int(v)
		}
		if v, ok := reqR["results"].([]any); ok {
			if data, err := json.Marshal(v); err == nil {
				round.Results = string(data)
			}
		}
		if v, ok := reqR["options"].([]any); ok {
			if data, err := json.Marshal(v); err == nil {
				round.Options = string(data)
			}
		}
		db.Create(&round)
	}
	response.Ok(reqRounds).WriteTo(w)
}

func parseUint(s string) uint {
	var id uint
	for _, c := range s {
		if c >= '0' && c <= '9' {
			id = id*10 + uint(c-'0')
		}
	}
	return id
}

func ListTemplateTopics(w http.ResponseWriter, r *http.Request) {
	var topics []topicModel.Topic
	db.Where("is_template = ?", true).Find(&topics)
	response.Ok(topics).WriteTo(w)
}

func GetTemplateOptions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID := parseUint(vars["id"])
	var options []topicModel.Option
	db.Where("topic_id = ? AND is_template = ?", topicID, true).Find(&options)
	if len(options) == 0 {
		db.Where("topic_id = ?", topicID).Find(&options)
	}
	response.Ok(options).WriteTo(w)
}

func ImportTemplates(w http.ResponseWriter, r *http.Request) {
	var req []struct {
		TemplateID uint   `json:"templateId"`
		Title      string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(http.StatusBadRequest, "invalid request body").WriteTo(w)
		return
	}
	if len(req) == 0 {
		response.Error(http.StatusBadRequest, "no template ids").WriteTo(w)
		return
	}

	var templateIDs []uint
	for _, r := range req {
		templateIDs = append(templateIDs, r.TemplateID)
	}

	var templates []topicModel.Topic
	db.Where("id IN ? AND is_template = ?", templateIDs, true).Find(&templates)
	if len(templates) == 0 {
		response.Error(http.StatusNotFound, "templates not found").WriteTo(w)
		return
	}

	templateMap := make(map[uint]string)
	for _, r := range req {
		templateMap[r.TemplateID] = r.Title
	}

	var importedTopics []topicModel.Topic
	for _, t := range templates {
		newTitle := templateMap[t.ID]
		if newTitle == "" {
			newTitle = t.Title
		}
		newTopic := topicModel.Topic{
			Title:       newTitle,
			Description: t.Description,
			IsActive:    true,
			IsTemplate:  false,
		}
		db.Create(&newTopic)

		var sourceOpts []topicModel.Option
		db.Where("topic_id = ?", t.ID).Find(&sourceOpts)
		var newOpts []topicModel.Option
		for i, opt := range sourceOpts {
			newOpts = append(newOpts, topicModel.Option{
				Title:      opt.Title,
				TopicID:    newTopic.ID,
				Importance: 0,
				Necessity:  0,
				SortOrder:  i,
				IsActive:   true,
			})
		}
		db.Create(&newOpts)
		importedTopics = append(importedTopics, newTopic)
	}
	response.Ok(importedTopics).WriteTo(w)
}

func SaveAsTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID := parseUint(vars["id"])

	var sourceTopic topicModel.Topic
	if err := db.First(&sourceTopic, topicID).Error; err != nil {
		response.Error(http.StatusNotFound, "topic not found").WriteTo(w)
		return
	}

	var req struct {
		Title string `json:"title"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	newTitle := req.Title
	if newTitle == "" {
		newTitle = sourceTopic.Title
	}

	newTopic := topicModel.Topic{
		Title:       newTitle,
		Description: sourceTopic.Description,
		IsActive:    true,
		IsTemplate:  true,
	}
	db.Create(&newTopic)

	var sourceOpts []topicModel.Option
	db.Where("topic_id = ?", topicID).Find(&sourceOpts)
	var newOpts []topicModel.Option
	for i, opt := range sourceOpts {
		newOpts = append(newOpts, topicModel.Option{
			Title:      opt.Title,
			TopicID:    newTopic.ID,
			Importance: 0,
			Necessity:  0,
			SortOrder:  i,
			IsActive:   true,
		})
	}
	db.Create(&newOpts)
	response.Ok(newTopic).WriteTo(w)
}

func SeedTopics(w http.ResponseWriter, r *http.Request) {
	var count int64
	db.Model(&topicModel.Topic{}).Where("is_template = ?", true).Count(&count)
	if count >= 40 {
		response.Ok(map[string]interface{}{"count": int(count), "message": "templates already exist"}).WriteTo(w)
		return
	}
	var ids []uint
	db.Model(&topicModel.Topic{}).Where("is_template = ?", true).Pluck("id", &ids)
	if len(ids) > 0 {
		db.Where("topic_id IN ?", ids).Delete(&topicModel.Option{})
		db.Where("id IN ?", ids).Delete(&topicModel.Topic{})
	}
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
		{
			Title:       "餐厅选择",
			Description: "出去吃饭如何选餐厅",
			Options:     []string{"菜品口味", "人均价格", "地理位置", "环境氛围", "服务质量", "卫生状况", "排队时间", "特色推荐"},
		},
		{
			Title:       "服装购买",
			Description: "买衣服考虑的因素",
			Options:     []string{"款式设计", "价格预算", "面料质量", "尺码合身", "品牌知名度", "颜色搭配", "穿着场合", "耐穿程度"},
		},
		{
			Title:       "礼物选择",
			Description: "送礼物如何挑选",
			Options:     []string{"对方喜好", "礼物寓意", "价格范围", "实用价值", "包装精美", "独特程度", "场合适宜", "收藏价值"},
		},
		{
			Title:       "宠物选择",
			Description: "养什么宠物好",
			Options:     []string{"饲养成本", "陪伴时间", "活动空间", "喂养难度", "互动性", "寿命长短", "气味大小", "掉毛程度"},
		},
		{
			Title:       "假期安排",
			Description: "假期怎么安排",
			Options:     []string{"放松休息", "外出旅游", "学习充电", "陪伴家人", "社交聚会", "运动健身", "兴趣爱好", "加班工作"},
		},
		{
			Title:       "视频平台",
			Description: "选择视频平台",
			Options:     []string{"内容丰富", "会员价格", "画质清晰", "广告多少", "独家资源", "使用体验", "兼容设备", "社区氛围"},
		},
		{
			Title:       "音乐App",
			Description: "选择音乐播放器",
			Options:     []string{"曲库规模", "音质效果", "会员费用", "个性化推荐", "离线下载", "界面设计", "版权问题", "社交功能"},
		},
		{
			Title:       "购物平台",
			Description: "选择网购平台",
			Options:     []string{"商品价格", "商品质量", "物流速度", "售后服务", "商品种类", "用户评价", "优惠活动", "支付安全"},
		},
		{
			Title:       "银行选择",
			Description: "选择哪家银行",
			Options:     []string{"网点便利", "利率高低", "手续费多少", "服务质量", "安全可靠", "产品丰富", "科技体验", "品牌口碑"},
		},
		{
			Title:       "宽带选择",
			Description: "办理宽带选哪家",
			Options:     []string{"网络速度", "套餐价格", "覆盖范围", "稳定性", "安装便捷", "售后服务", "合约期限", "赠送设备"},
		},
		{
			Title:       "健身房",
			Description: "选择健身房",
			Options:     []string{"距离远近", "价格费用", "器械设备", "环境氛围", "教练水平", "营业时间", "会员人数", "附加服务"},
		},
		{
			Title:       "在线课程平台",
			Description: "选择学习平台",
			Options:     []string{"课程质量", "价格实惠", "师资力量", "互动答疑", "证书认可", "播放流畅", "内容更新", "社区氛围"},
		},
		{
			Title:       "阅读选择",
			Description: "选择看什么书",
			Options:     []string{"知识价值", "趣味性", "阅读难度", "篇幅长短", "作者声誉", "豆瓣评分", "实用程度", "口碑推荐"},
		},
		{
			Title:       "周末娱乐",
			Description: "周末怎么玩",
			Options:     []string{"看电影", "逛商场", "玩游戏", "运动健身", "读书学习", "朋友聚会", "户外活动", "宅家休息"},
		},
		{
			Title:       "家具选购",
			Description: "买家具考虑什么",
			Options:     []string{"价格预算", "材质质量", "风格搭配", "尺寸合适", "环保安全", "耐用程度", "品牌口碑", "配送安装"},
		},
		{
			Title:       "家电购买",
			Description: "买家电考虑因素",
			Options:     []string{"品牌质量", "价格预算", "功能实用", "能耗等级", "售后保障", "外观设计", "操作便捷", "尺寸大小"},
		},
		{
			Title:       "找对象方式",
			Description: "如何认识TA",
			Options:     []string{"相亲介绍", "自由恋爱", "社交软件", "朋友介绍", "工作认识", "兴趣社群", "线下活动", "旅行偶遇"},
		},
		{
			Title:       "买房位置",
			Description: "房子买在哪里",
			Options:     []string{"交通便利", "学区资源", "商业配套", "环境绿化", "升值潜力", "价格预算", "开发商品牌", "户型设计"},
		},
		{
			Title:       "投资方向",
			Description: "钱投到哪里",
			Options:     []string{"股票基金", "银行理财", "债券投资", "房产投资", "黄金避险", "数字货币", "保险理财", "定期存款"},
		},
		{
			Title:       "考证规划",
			Description: "考什么证书",
			Options:     []string{"含金量高", "考试难度", "学习时间", "费用成本", "就业帮助", "有效期限", "专业相关", "个人兴趣"},
		},
	}
	for _, t := range defaultTopics {
		topic := topicModel.Topic{Title: t.Title, Description: t.Description, IsActive: true, IsTemplate: true}
		db.Create(&topic)
		var opts []topicModel.Option
		for i, title := range t.Options {
			opts = append(opts, topicModel.Option{
				Title:      title,
				TopicID:    topic.ID,
				Importance: 0,
				Necessity:  0,
				SortOrder:  i,
				IsActive:   true,
			})
		}
		db.Create(&opts)
	}
	response.Ok(map[string]int{"count": len(defaultTopics)}).WriteTo(w)
}
