
# wifi-cracker

轻量级 Wi-Fi 密码测试工具，支持**单目标**、**多目标**、**自动排序**三种破解模式，内置**失败密码库**避免重复尝试，优先验证系统已保存密码，全程日志可追溯。

---

## ✨ 核心功能

| 功能 | 描述 |
| ---- | ---- |
| 🔍 扫描 | 实时扫描周围 Wi-Fi 并按信号强度排序 |
| 🔑 优先验证 | 自动尝试系统已保存密码，成功即跳过 |
| 📚 失败库 | JSON 持久化记录失败/异常/成功密码，重启不复扫 |
| ⚙️ 多模式 | single / multi / auto / verify / show-saved |
| ⏱️ 超时 | 单密码连接超时自定义（1-60 s） |
| 🗂️ 自定义库 | 指定字典文件、失败库路径 |
| 📃 日志 | 本地文件记录全部尝试与结果 |

---

## 🚀 一键安装

```bash
# 克隆
git clone https://github.com/yearnming/wifi.git
cd wifi-cracker

# 运行（Windows 需管理员权限）
go run main.go --help
```

---

## 📖 使用示例

```bash
# 1. 显示帮助
go run main.go

# 2. 破解单个指定 Wi-Fi
go run main.go --mode single --ssid home_5G -l top1000.txt

# 3. 破解多个（交互选编号）
go run main.go --mode multi -l top1000.txt

# 4. 自动按信号强度全扫
go run main.go --mode auto -t 5

# 5. 仅验证已保存密码（不爆破）
go run main.go --mode verify -t 3

# 6. 查看系统保存的密码
go run main.go --mode show-saved

# 7. 自定义失败库路径
go run main.go --mode auto -d /tmp/mydb.json
```

---

## 📂 文件说明

| 文件 | 作用 |
| ---- | ---- |
| `main.go` | 入口、参数解析、生命周期 |
| `pkg/wifi/` | 核心破解逻辑、状态归类 |
| `pkg/filedb/` | 失败/异常/成功密码持久化 |
| `pkg/scanpasswd/` | 读取系统已保存密码 |
| `filedb.json` | **默认失败库**（生成在同目录） |
| `logs/wifi-crack.log` | 完整运行日志 |

---

## ⚙️ 参数一览

```
Flags:
   -t, --timeout int      单密码连接超时（秒，1-60，默认 4）
   -l, --dict string      密码字典文件（默认内置）
   -m, --mode string      运行模式：single/multi/auto/verify/show-saved
   -d, --fail-db string   失败库 JSON 路径（默认./filedb.json）
   -h, --help             显示本帮助
```

---

## 🛡️ 安全声明

* 仅用于**合法授权**的渗透测试或自我审计  
* 成功密码默认保存在**当前用户目录**，请妥善保管  
* 开发者不对任何非法使用承担责任

---

## 🤝 参与贡献

欢迎提 Issue / PR，一起让工具更强大！

---

## 📄 License

MIT © 2025 yearnming
```