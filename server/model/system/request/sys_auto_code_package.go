package request

import (
	model "ygpt/server/model/system"
)

type SysAutoCodePackageCreate struct {
	Desc        string `json:"desc" example:"描述"`
	Label       string `json:"label" example:"展示名"`
	Template    string `json:"template"  example:"模版"`
	PackageName string `json:"packageName" example:"包名"`
}

func (r *SysAutoCodePackageCreate) AutoCode() AutoCode {
	return AutoCode{
		Package: r.PackageName,
	}
}

func (r *SysAutoCodePackageCreate) Create() model.SysAutoCodePackage {
	return model.SysAutoCodePackage{
		Desc:        r.Desc,
		Label:       r.Label,
		Template:    r.Template,
		PackageName: r.PackageName,
	}
}
