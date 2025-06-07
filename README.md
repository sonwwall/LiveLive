


# 🎓 LiveLive - 教育直播平台后端服务

> LiveLive 是一个专为教育场景打造的高性能直播平台后端，支持直播管理、互动聊天、答题统计、AI 课堂分析等功能，适用于在线课堂、远程教学等多种场景。

---

## ✨ 项目特点

- 📺 支持 **RTMP 推流** 与 **HTTP-FLV 播放**（基于 LiveGo）
- 💬 实时聊天室系统（基于 WebSocket）
- 🧠 AI 自动分析课堂内容，生成学习报告（基于火山引擎大模型）
- 📊 在线答题系统，自动统计正确率和分布
- 🔁 直播录制与回放功能
- 🏫 支持多教室并发直播管理
- ✅ 用户鉴权、课程管理等完整教学功能
- ⚙️ 分布式服务架构，支持水平扩展

---

## 🛠 技术栈

| 类别 | 技术 |
|------|------|
| 语言 | Go (Golang) |
| 框架 | [Kitex](https://www.cloudwego.io/), [hz](https://www.cloudwego.io/zh/docs/hertz/), [Eino](https://www.cloudwego.io/zh/docs/eino/) |
| 数据库 | MySQL, Redis, Etcd |
| 消息队列 | Kafka |
| 流媒体服务 | LiveGo |
| 长连接 | WebSocket |
| 容器化部署 | Docker, Docker Compose |
| AI 服务 | 火山引擎语音识别、大模型，阿里云 OSS |

---

## 📚 已实现功能

1. 👤 用户注册、登录、鉴权、资料查询
2. 📡 发起/观看直播，支持录制与回放
3. ✅ 答题功能（选择题/判断题），自动统计分布与正确率
4. 🕐 实时签到，支持截止时间后自动统计
5. 💬 聊天室系统（WebSocket 长连接）
6. 📖 AI 自动分析课程内容，生成结构化学习报告

---

## 🔗 接口文档

📘 [点击查看接口文档（Apipost）](https://doc.apipost.net/docs/detail/486962bd70b1000?target_id=2f4d22e9bd00a4)

---

## 🚀 快速开始（Quick Start）

### 1️⃣ 克隆项目

```bash
git clone https://github.com/sonwwall/LiveLive.git
cd LiveLive
````

---

### 2️⃣ 修改配置文件

编辑 `config/` 文件夹下的两个配置文件：

* `db.yaml`：配置 MySQL、Redis 等基础组件地址
* `ai.yaml`：配置阿里云 OSS、火山引擎语音识别与 AI 大模型的 AK/SK

---

### 3️⃣ 启动依赖服务

使用 Docker 启动 MySQL、Redis、Kafka、Etcd、LiveGo 等服务：

```bash
docker-compose up -d
```

启动的服务包括：

* MySQL
* Redis
* Etcd
* Kafka + Zookeeper
* LiveGo（推流/播放服务）

---

### 4️⃣ 安装 FFmpeg

#### macOS：

```bash
brew install ffmpeg
```

#### Windows：

访问官网下载安装包：[https://ffmpeg.org/download.html](https://ffmpeg.org/download.html)

---

### 5️⃣ 使用 OBS 推流测试

使用接口文档获取：

* 推流地址（RTMP）
* 推流密钥（或直播间 ID）

填入 OBS 配置中，即可发起直播。

---

### 6️⃣ 观看直播流

从接口文档中获取播放地址（HTTP-FLV），支持以下方式播放：

* 浏览器访问（会触发下载或使用 flv.js 播放）
* 使用 VLC 播放器(修改前缀为rtmp://)

 




---


## 🙋‍♂️ 联系作者

如果你对项目有兴趣或问题，欢迎 Issue 或 PR！

```
