package models

import (
	"cloudiac/runner"
	"path"
)

type BaseTask struct {
	SoftDeleteModel

	/* 通用任务参数 */
	Type string `json:"type" gorm:"not null;enum('plan','apply','destroy','scan')" enums:"'plan','apply','destroy','scan'"` // 任务类型。1. plan: 计划 2. apply: 部署 3. destroy: 销毁

	Pipeline string       `json:"-" gorm:"type:json"`        // 用户自定义 pipeline 内容
	Flow     PipelineTask `json:"-" gorm:"type:json"`        // 实际生成的任务执行流程
	CurrStep int          `json:"currStep" gorm:"default:0"` // 当前在执行的流程步骤

	ContainerId string `json:"-" gorm:"size:64"`

	// 任务每一步的执行超时(整个任务无超时控制)
	StepTimeout int `json:"stepTimeout" gorm:"default:600;comment:执行超时"`

	RunnerId string `json:"runnerId" gorm:"not null"` // 部署通道

	Status  string `json:"status" gorm:"type:enum('pending','running','approving','rejected','failed','complete','timeout');default:'pending'" enums:"'pending','running','approving','rejected','failed','complete','timeout'"`
	Message string `json:"message" gorm:"type:text"` // 任务的状态描述信息，如失败原因等

	StartAt *Time `json:"startAt" gorm:"type:datetime;comment:任务开始时间"` // 任务开始时间
	EndAt   *Time `json:"endAt" gorm:"type:datetime;comment:任务结束时间"`   // 任务结束时间
}

// ScanTask 合规扫描任务
type ScanTask struct {
	BaseTask

	// 模板扫描任务参数
	OrgId     Id `json:"orgId" gorm:"size:32;not null"` // 组织ID
	ProjectId Id `json:"projectId" gorm:"size:32"`      // 项目ID
	TplId     Id `json:"tplId" gorm:"size:32"`          // 模板ID
	EnvId     Id `json:"envId" gorm:"size:32"`          // 环境ID

	Name      string `json:"name" gorm:"not null;comment:任务名称"` // 任务名称
	CreatorId Id     `json:"creatorId" gorm:"size:32;not null"` // 创建人ID

	RepoAddr string `json:"repoAddr" gorm:""`
	Revision string `json:"revision" gorm:""`
	CommitId string `json:"commitId" gorm:""` // 创建任务时 revision 对应的 commit id

	Workdir string `json:"workdir" gorm:"default:''"`

	Mirror       bool `json:"mirror"`       // 是否属于部署任务的扫描任务
	MirrorTaskId Id   `json:"mirrorTaskId"` // 部署任务ID

	PolicyStatus string `json:"policyStatus" gorm:"size:16" enums:"'passed','violated','pending','failed'"` // 策略检查结果

	// 扩展属性，包括 source, transitionId 等
	Extra TaskExtra `json:"extra" gorm:"type:json"` // 扩展属性
}

func (ScanTask) TableName() string {
	return "iac_scan_task"
}

func (t *ScanTask) TfParseJsonPath() string {
	if t.EnvId != "" {
		return path.Join(t.ProjectId.String(), t.EnvId.String(), t.Id.String(), runner.TerrascanJsonFile)
	} else {
		return path.Join(t.TplId.String(), t.Id.String(), runner.TerrascanJsonFile)
	}
}

func (t *ScanTask) TfResultJsonPath() string {
	if t.EnvId != "" {
		return path.Join(t.ProjectId.String(), t.EnvId.String(), t.Id.String(), runner.TerrascanResultFile)
	} else {
		return path.Join(t.TplId.String(), t.Id.String(), runner.TerrascanResultFile)
	}
}
