# 习惯养成应用 (Habit Tracker)

一个面向个人用户的习惯养成应用，帮助用户建立和坚持好习惯。

## 功能特性

### 🎯 习惯管理
- 创建多个习惯，支持自定义名称和颜色
- 四种打卡频率模式：
  - 每天打卡
  - 每周N次（一周内完成N次即可）
  - 指定星期（只在指定的星期几需要打卡）
  - 每月N次（一个月内完成N次）
- 可选的提醒时间设置
- 习惯列表支持拖拽排序
- 修改频率规则后连续天数重新计算（历史数据保留）
- 习惯归档和彻底删除

### ✅ 打卡机制
- 点击习惯卡片完成今日打卡
- 支持补打卡（过去7天内）
- 当天打卡可以取消
- 幂等操作：重复提交不产生重复记录
- 打卡成功显示随机鼓励文案

### 🔥 连续天数计算
- 每天类型：从今天往前数每天都有打卡记录则连续
- 每周N次类型：以自然周为单位，每周打卡次数达标则连续
- 指定星期类型：只看指定的那几天是否都打卡了
- 每月N次类型：以自然月为单位，每月打卡次数达标则连续
- 记录每个习惯的历史最长连续天数

### 🏆 成就系统
预设15个成就徽章，分为四类：
- **连续打卡类**：连续7天、21天、30天、60天、100天、365天
- **累计打卡类**：累计50次、100次、500次、1000次
- **习惯数量类**：同时维护3个、5个、10个活跃习惯
- **特殊类**：完美一周、完美一月

### 📊 统计与可视化
- 打卡热力图（过去一年）
- 单习惯统计（本周/本月完成率、总体完成率、最长连续天数等）
- 总览统计（综合完成率、今日完成进度）

## 技术栈

- **前端**: Svelte + Vite + Tailwind CSS
- **后端**: Go + Gin 框架
- **数据库**: PostgreSQL
- **部署**: Docker Compose

## 快速开始

### 使用 Docker Compose（推荐）

```bash
# 克隆项目
git clone <repository-url>
cd habit-tracker

# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

访问地址：
- 前端：http://localhost:5173
- 后端API：http://localhost:8080
- PostgreSQL：localhost:5432

### 本地开发

#### 后端
```bash
cd backend

# 安装依赖
go mod download

# 复制环境变量
cp .env.example .env

# 运行
go run cmd/server/main.go
```

#### 前端
```bash
cd frontend

# 安装依赖
npm install

# 运行
npm run dev
```

## API 接口

### 习惯管理
- `GET /api/habits` - 获取习惯列表
- `POST /api/habits` - 创建习惯
- `POST /api/habits/reorder` - 重新排序习惯
- `GET /api/habits/:id` - 获取单个习惯
- `PUT /api/habits/:id` - 更新习惯
- `DELETE /api/habits/:id` - 删除习惯

### 打卡
- `POST /api/habits/:id/checkin?date=YYYY-MM-DD` - 打卡
- `DELETE /api/habits/:id/checkin?date=YYYY-MM-DD` - 取消打卡
- `GET /api/habits/:id/checkins` - 获取打卡记录

### 统计
- `GET /api/stats/overview` - 概览统计
- `GET /api/stats/heatmap` - 热力图数据
- `GET /api/stats/achievements` - 成就列表

## 项目结构

```
habit-tracker/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go          # 应用入口
│   ├── internal/
│   │   ├── models/              # 数据模型
│   │   ├── database/            # 数据库连接
│   │   ├── handlers/            # API 处理器
│   │   └── services/            # 业务逻辑
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── components/          # Svelte 组件
│   │   └── App.svelte           # 主应用
│   ├── package.json
│   ├── vite.config.js
│   └── Dockerfile
└── docker-compose.yml
```

## 数据库模型

### Habit (习惯)
- id: 主键
- name: 名称
- color: 颜色
- frequency_type: 频率类型
- frequency_value: 频率值
- specific_days: 指定星期
- reminder_time: 提醒时间
- current_streak: 当前连续天数
- longest_streak: 最长连续天数
- total_checkins: 累计打卡次数
- sort_order: 排序
- is_archived: 是否已归档
- frequency_modified_at: 频率修改时间
- created_at, updated_at: 时间戳

### Checkin (打卡记录)
- id: 主键
- habit_id: 习惯ID
- date: 日期
- note: 备注
- is_backfill: 是否补打卡
- created_at, updated_at: 时间戳

### Achievement (成就定义)
- id: 主键
- type: 成就类型
- name: 名称
- description: 描述
- target: 目标值
- icon: 图标
- created_at: 创建时间

### UserAchievement (用户获得的成就)
- id: 主键
- achievement_id: 成就ID
- habit_id: 相关习惯ID
- earned_at: 获得时间

## 注意事项

1. 时区处理：日期判定基于用户本地时区
2. 数据一致性：补打卡和取消打卡后，相关连续天数和成就状态会同步更新
3. 成就一旦获得不可撤销
4. 修改习惯频率会重置连续天数（历史数据保留）

## License

MIT
