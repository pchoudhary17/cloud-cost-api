package model

import "time"

type GroupedByData struct {
	GroupedBy          string
	Cost, LifetimeCost float64
	Clusters           int64
}

type Instance struct {
	InstanceId                                                                string `gorm:"primaryKey"`
	InstanceName, InstanceType, InstanceState, Region, Account                string
	LaunchTime, ExitTime                                                      time.Time
	ClusterName, ClusterState, ClusterType                                    string
	ExecutionTime, Cost, EstimatedCost, LifetimeCost, CPU_Usage, MaxCPU_Usage float64
	Owner, Manager, L4                                                        string
}

type Clusters struct {
	ClusterName, ClusterState, ClusterType, Owner, Account, Region, Manager, L4 string
	ExecutionTime, Cost, EstimatedCost, LifetimeCost, CpuUsage, MaxCpuUsage     float64
}

type Tabler interface {
	TableName() string
}

// TableName overrides the table name used by Instance to `instances`
func (Instance) TableName() string {
	return "instances"
}
