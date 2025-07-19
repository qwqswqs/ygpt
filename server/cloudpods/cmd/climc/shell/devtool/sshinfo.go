// Copyright 2019 Yunion
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

package devtool

import (
	"yunion.io/x/onecloud/cmd/climc/shell"
	modules "yunion.io/x/onecloud/pkg/mcclient/modules/devtool"
	options "yunion.io/x/onecloud/pkg/mcclient/options/devtool"
)

func init() {
	cmd := shell.NewResourceCmd(&modules.DevToolSshInfos).WithKeyword("devtool-sshinfo")
	cmd.List(new(options.SshInfoListOptions))
	cmd.Create(new(options.SshInfoCreateOptions))
	cmd.Show(new(options.SshInfoOptions))
	cmd.Delete(new(options.SshInfoOptions))
}
