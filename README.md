# 多轮打分筛选决策系统mdecide

## 项目简介

基于"多轮打分筛选法"的决策辅助工具，帮助用户通过结构化、可重复的评估流程做出更理性的选择。系统支持创建任意决策场景，对筛选条件进行重要性（1-5分）和必要性（0/1）打分，并通过多轮（支持跨天）权重调整，最终生成带总结的决策结果。

**核心决策逻辑**：
1. 罗列筛选条件
2. 重要性打分（1-5）→ 排序
3. 必要性打分（0/1）
4. 30分钟内得出初步结果
5. 第二天/多轮重新打分，获得更稳定结论

## 技术栈

| 模块 | 技术选型 |
|------|----------|
| 前端Web | Vue 3 + Vite + Element Plus + Axios |
| 小程序 | 微信小程序原生框架 + WeUI |
| 后端 | Go + go-zero微服务框架 + MySQL + Redis |

## 项目结构
├── README.md
├── backend/ # go-zero后端服务
│ ├── app/ # 应用模块（按业务领域划分）
│ │ ├── topic/ # 决策主题管理
│ │ │ ├── handler/ # HTTP处理器
│ │ │ └── model/ # 数据模型
│ │ ├── condition/ # 筛选条件管理
│ │ │ ├── builder/ # 条件构建器
│ │ │ ├── handler/ # HTTP处理器
│ │ │ ├── logic/ # 业务逻辑
│ │ │ ├── model/ # 数据模型
│ │ │ └── validate/ # 参数校验
│ │ ├── options/ # 选项配置（重要性/必要性选项）
│ │ │ ├── builder/ # 选项构建器
│ │ │ ├── desc/ # 选项描述
│ │ │ ├── handler/ # HTTP处理器
│ │ │ ├── logic/ # 业务逻辑
│ │ │ ├── model/ # 数据模型
│ │ │ └── validate/ # 参数校验
│ │ ├── round/ # 轮次管理
│ │ │ ├── builder/ # 轮次构建器
│ │ │ ├── handler/ # HTTP处理器
│ │ │ ├── logic/ # 业务逻辑
│ │ │ ├── model/ # 数据模型
│ │ │ └── validate/ # 参数校验
│ │ └── result/ # 结果分析
│ │ ├── builder/ # 结果构建器
│ │ ├── handler/ # HTTP处理器
│ │ ├── logic/ # 业务逻辑
│ │ ├── model/ # 数据模型
│ │ └── validate/ # 参数校验
│ ├── common/ # 公共组件
│ │ ├── ctxdata/ # 上下文数据传递
│ │ ├── errorx/ # 错误定义与处理
│ │ ├── response/ # 统一响应格式
│ │ └── verify/ # 通用验证器
│ └── etc/ # 配置文件
├── frontend/ # Vue3前端
│ └── src/
│ ├── api/ # API接口封装
│ ├── components/ # 公共组件
│ ├── router/ # 路由配置
│ ├── stores/ # Pinia状态管理
│ ├── utils/ # 工具函数
│ └── views/ # 页面视图
│ ├── topic/ # 决策主题页
│ ├── condition/ # 条件管理页
│ ├── scoring/ # 打分页
│ └── result/ # 结果页
└── miniprogram/ # 微信小程序
└── pages/
├── index/ # 首页
├── topic/ # 决策主题
├── condition/ # 条件管理
├── scoring/ # 打分页面
└── result/ # 结果页面
