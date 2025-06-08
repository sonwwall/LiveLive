
---

# 🎓 LiveLive - 教育直播平台后端服务

> **LiveLive** 是一个专为教育场景打造的高性能直播平台后端，支持直播管理、互动聊天、答题统计、AI 课堂分析等功能，适用于在线课堂、远程教学等多种场景。

---

## ✨ 项目特点

* 📺 支持 **RTMP 推流** 与 **HTTP-FLV 播放**（基于 LiveGo）
* 💬 实时聊天室系统（基于 WebSocket）
* 🧠 AI 自动分析课堂内容，生成学习报告（基于火山引擎大模型）
* 📊 在线答题系统，自动统计正确率和分布
* 🕘 实时签到，自动统计出勤情况
* 🎥 直播录制与回放功能
* 🏫 支持多教室并发直播管理
* ✅ 用户鉴权、课程管理等完整教学功能
* 🧠 使用缓存减少数据库访问，提高响应速度
* 📩 使用消息队列提升系统吞吐能力
* ⚙️ 分布式架构，支持服务水平扩展

---

## 🛠 技术栈

| 类别            | 技术                                                                                                                       |
|---------------|--------------------------------------------------------------------------------------------------------------------------|
| 语言            | Go (Golang)                                                                                                              |
| 框架            | [Kitex](https://www.cloudwego.io/), [hz](https://www.cloudwego.io/zh/docs/hertz/), [Eino](https://www.cloudwego.io/zh/docs/eino/) |
| 数据存储          | MySQL, Redis, Etcd                                                                                                       |
| 消息队列          | Kafka                                                                                                                    |
| 流媒体           | LiveGo                                                                                                                   |
| 长连接           | WebSocket                                                                                                                |
| 部署            | Docker, Docker Compose                                                                                                   |
| AI 服务         | 火山引擎语音识别 & 大模型, 阿里云 OSS                                                                                                  |
| 视频流拉取，转码      | FFmpeg                                                                                                                   |

---

---

## 架构图
![img_5.png](img_5.png)


---

## 📚 已实现功能

* 👤 用户注册、登录、鉴权、资料查询
* 📡 发起/观看直播，支持录制与回放
* ✅ 答题功能（选择题/判断题），自动统计分布与正确率
* 🕐 实时签到，支持截止时间后自动统计
* 💬 聊天室系统（WebSocket 长连接）
* 📖 AI 自动分析课程内容，生成结构化学习报告
* 🥞 完整的课程管理系统，支持课程创建与邀请码加入

---

## 🔗 接口文档

📘 [👉 点击查看接口文档（Apipost）](https://doc.apipost.net/docs/487fb69e1cb1000?locale=zh-cn)

---

## 🚀 快速开始（Quick Start）

### 1️⃣ 克隆项目

```bash
git clone https://github.com/sonwwall/LiveLive.git
cd LiveLive
```

---

### 2️⃣ 修改配置文件

编辑 `config/` 目录下的配置文件：

* `db.yaml`：配置 MySQL、Redis 等组件地址
* `ai.yaml`：配置阿里云 OSS、火山引擎语音识别与 AI 大模型的 AK/SK

---

### 3️⃣ 启动依赖服务

使用 Docker 启动所有依赖服务：

```bash
docker-compose up -d
```

> 启动组件：
>
> * MySQL
> * Redis
> * Etcd
> * Kafka + Zookeeper
> * LiveGo（推流/播放服务）

---

### 4️⃣ 安装 FFmpeg

#### macOS：

```bash
brew install ffmpeg
```

#### Windows：

👉 [官网下载 FFmpeg](https://ffmpeg.org/download.html)

---

### 5️⃣ 启动项目服务

分别执行以下命令：

```bash
go run ./rpc/user
go run ./rpc/ai
go run ./rpc/live
go run ./rpc/course
go run ./rpc/websocket
go run ./rpc/quiz
go run ./job
go run ./api
```

---

### 6️⃣ 使用 OBS 推流测试

1. 参考接口文档，获取 RTMP 推流地址 & 推流密钥（或直播间 ID）
2. 在 OBS 中填入后开始推流测试

---

### 7️⃣ 观看直播流

通过接口文档获取播放地址（HTTP-FLV）：

* 浏览器播放（推荐使用 flv.js）
* VLC 播放器（地址前缀需改为 `rtmp://`）

---

## 🤖 AI 总结功能使用说明

1. **建立连接**：老师与学生需通过 WebSocket 接口加入课堂（详见接口文档）
2. **老师开始直播**：调用「开始录制」接口
3. **直播结束**：调用「结束录制」接口
4. **生成报告**：AI 分析完成后，报告会通过 WebSocket 实时推送给老师与学生

---

## 💡 其他说明

* 聊天室、问答、签到均依赖 WebSocket 实时推送机制
* 相关接口参数与使用方式详见接口文档
* **FFmpeg 连接问题解决方案**：
  如果遇到录制回放，AI功能异常，修改以下文件：

  文件：`rpc/live/handler.go`

    * **取消注释**第 293 行，注释第 295 行
    * **取消注释**第 374-381 行，注释第 383-390 行
* 视频录制默认放在recordings目录下,转码的音频默认放在audios目录下

---

## 🙋‍♂️ 联系作者

欢迎 Issue / PR，或通过 GitHub 联系作者！


---


