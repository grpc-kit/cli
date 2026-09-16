# 业务后台模版

基于 `@grpc-kit/adm` 的最小化业务后台骨架，用于新服务初始化管理后台。

## 初始化步骤

1. 将本目录复制到目标服务的 `web/admin/` 下：

```bash
cp -r template/web/admin /path/to/your-service/web/admin
```

2. 修改 `package.json`：

- `name` 改为你的服务名，如 `"myproduct-myservice-admin"`
- `@grpc-kit/adm` 的路径改为实际的模版包路径：

```json
"@grpc-kit/adm": "file:/absolute/path/to/github.com/grpc-kit/adm/vite-antd"
```

3. 修改 `index.html` 中的 `<title>` 和 `src/main.tsx` 中的 `title` 为你的服务名称。

4. 安装依赖并启动：

```bash
cd web/admin
npm install
npm run dev
```

5. 构建（产物输出到 `../../public/admin`）：

```bash
npm run build
```

## 目录结构

```
web/admin/
  package.json          # 项目依赖与脚本
  vite.config.ts        # 构建配置、alias、代理
  index.html            # HTML 入口
  src/
    main.tsx            # 应用入口，调用 createApp
    routes.tsx          # 业务路由清单（集中注册）
    vite-env.d.ts       # Vite 类型声明
    pages/              # 业务页面（按领域分目录）
      example/
        demo.tsx
    api/                # 业务 API 客户端（按 proto 自动生成）
      {product}/{service}/v1/
        client.ts
        types.ts
        index.ts
    components/         # 业务自定义组件（按需创建）
```

## 核心约定

### 路由注册

所有业务路由统一在 `src/routes.tsx` 中注册，不得分散到其他文件：

```tsx
import type { RouteItem } from '@grpc-kit/adm';
import MyPage from './pages/my-domain/my-page';

export const businessRoutes: RouteItem[] = [
  { path: '/my-domain/my-page', element: <MyPage /> },
];
```

### API 客户端

业务 API 客户端统一使用模版导出的 httpClient：

```ts
import { httpClient } from '@grpc-kit/adm';
```

API 客户端文件放在 `src/api/{productCode}/{shortName}/v1/`，通常由构建工具从 `microservice.gateway.yaml` 自动生成。

### 路径别名

Vite 和 TypeScript 中已配置以下别名：

| 别名 | 指向 | 用途 |
|------|------|------|
| `@/` | `src/` | 本地源码 |
| `@api/` | `src/api/` | 本地业务 API |
| `@api/known/` | 模版包内置 API | 登录、用户、权限等通用接口 |

### 模版扩展

通过 `createApp` 的配置参数扩展模版能力：

```tsx
import { createApp } from '@grpc-kit/adm';
import { businessRoutes } from './routes';

createApp({
  title: 'My Service Admin',     // 浏览器侧边栏标题
  routes: businessRoutes,         // 业务路由
  menuIcons: {                    // 额外图标映射（可选）
    'MyIcon': <MyIconComponent />,
  },
});
```

## 新增业务页面

1. 在 `src/pages/` 下按领域创建目录和页面文件：

```
src/pages/device/
  list.tsx
```

2. 在 `src/routes.tsx` 中注册路由：

```tsx
import DeviceList from './pages/device/list';

export const businessRoutes: RouteItem[] = [
  // ...existing routes
  { path: '/device/list', element: <DeviceList /> },
];
```

3. 验证：

```bash
npx tsc --noEmit
npm run build
```

## 质量门禁

提交前必须通过：

```bash
npx tsc --noEmit   # 类型检查
npm run build       # 构建
npm run lint        # ESLint（建议）
```
