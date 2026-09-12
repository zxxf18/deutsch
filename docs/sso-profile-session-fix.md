# SSO 资料与退出修复（2026-09-12）

## 资料归属

Casdoor 继续负责邮箱验证、密码、认证。Deutsch 只保存业务资料，不恢复密码登录或邀请码认证。

固定 Casdoor issuer 下，`sub` 为不可变 UUID。业务 `users.id` 使用同一 UUID，保持学习记录、偏好和错题归属不变，不按邮箱或用户名关联账户。

认证后的 `/api/v1/auth/me` 和受保护请求首次访问时自动建立本地用户；重复访问保留昵称、简介。并发首次创建按主键重读；停用、软删账户不恢复；数据库失败不放行。新增资料的 `password_encrypted` 留空，原有表和数据不删除。

`PATCH /api/v1/user/profile` 的 `nickname`、`description` 都可省略：省略或 null 保留原值，空字符串清空。最大长度分别为 50、500 个 Unicode 字符。修改成功后 `/api/v1/auth/me` 读取已保存昵称和简介，刷新不回退到旧 Cookie 的昵称。

## 前端与退出

四个应用的按钮显示“登录”。已登录按钮打开用户菜单，显式提供“用户资料”“退出登录”，长名称单行省略。Deutsch 保留本地个人中心；其余三个跳转 `https://sso.yebuluo.com.cn/account`，由 Casdoor 自带页面保存统一资料。

退出并行调用业务 POST logout 和 Casdoor `POST /api/logout`。中央请求必须 `credentials: include`，依赖现网按业务域名允许凭据的 CORS。两边都成功后才回首页；失败显示重试提示，仍尝试清除业务 Cookie，不能把中央失败当作退出成功。

范围：清除当前应用与当前浏览器 Casdoor 会话。其他已打开应用各自的业务 Cookie 不会由此被服务器撤销，跨设备/全部应用实时单点注销需要另外接 Casdoor Back-Channel Logout。本次没有实现或声称这种全局注销。

## 验证与发布

- Go：`go test ./...`、`go vet ./...`；含实际 Cookie 验证、中间件建档、HTTP PATCH 解析、保存、`me` 回读集成测试（隔离内存存储）。
- 前端：`node --experimental-strip-types --test tests/*.test.mjs`、`npm run lint`、`npm run build`。
- 独立 Chrome：320、390、768、1024、1440px；模拟专用测试身份，验证菜单打开不退出、Escape 关闭、资料链接、保存回显和双端退出请求。生产普通用户资料端到端验收需要已验证的专用测试账号，不修改验证标记或冒用现有用户。
- 现网中央注销：使用授权管理员在独立浏览器建立真实 Casdoor 会话，仅将业务 `/me` 替换为菜单展示用测试响应；四个业务和 Casdoor 的 logout 接口都走真实请求。四站均确认 Casdoor `get-account` 变为未认证，移除测试响应并重载真实主页后，再次点击登录显示密码认证表单。没有用未验证管理员绕过业务认证或修改现有用户资料。
- 本地构建 Linux amd64 后，仅替换 Deutsch 后端及四个前端镜像，保留环境变量、网络、数据挂载、重启策略。无新增服务/数据库。
- 回滚：使用部署时备份的容器配置与旧镜像恢复对应容器，重启 `public-nginx` 刷新上游地址。保留创建的业务资料，不回滚/删除用户数据。

## 审计记录

2026-09-12 生产依赖审计：Kidstar、Poetry 无告警；Deutsch 的旧 React Router 有 3 项 moderate（本次固定站内跳转不使用可控 `//` 路径）；Phonics 的旧 Next.js/Sharp 有 high/critical，但现网 vinext standalone 镜像不包含这两个包，且构建不启动其图片优化 HTTP 服务，告警路径不在本次部署的运行链路。未强制升级依赖或变更锁文件；依赖独立升级应在 2026-09-19 前复查。
