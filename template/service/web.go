// Copyright © 2020 The gRPC Kit Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

func (t *templateService) fileDirectoryWeb() {
	t.files = append(t.files, &templateFile{
		name: "web/README.md",
		body: `
# web

用于存放该微服务相关前端***源代码***，包含普通用户前台与管理员后台。

约定如下：

1. 目录 "./web/admin/" 面向管理员后台；
2. 目录 "./web/webroot/" 面向普通用户前台。

## admin

面向管理员后台，编译后的静态文件需存放至 "./public/admin/" 中。

## webroot

面向普通用户前台，编译后的静态文件需存放至 "./public/webroot/" 中。
`,
	})
}
