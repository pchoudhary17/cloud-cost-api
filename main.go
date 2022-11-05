package main

import (
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

func main() {
	StartServices()
}

func StartServices() {
	router := gin.Default()
	router.GET("/clusters", getClusters)
	router.Run("localhost:8080")
}

func getClusters(c *gin.Context) {

	groupBy := c.DefaultQuery("grouped_by", "none")
	beginDate := c.Query("startDate")
	endDate := c.Query("endDate")

	//Do these in frontend JS
	//#############################
	// clusterState := c.PostForm("clusterState")
	// records := c.PostForm("records")
	// reportName := c.DefaultQuery("reportName", "Cloud_Cost_Report"+beginDate+"-"+endDate) .
	//#############################

	dbHandler := ConnectToDb()
	clustersFromDB := GetClustersFromDB(dbHandler, beginDate, endDate)
	if groupBy == "none" {
		c.JSON(http.StatusOK, clustersFromDB)
	} else {
		c.JSON(http.StatusOK, groupData(clustersFromDB, groupBy))
	}
}

func groupData(clusters []Clusters, groupedBy string) []GroupedByData {
	var data []GroupedByData

	groups := make(map[string]GroupedByData)
	for _, c := range clusters {
		r := reflect.ValueOf(c)
		key := reflect.Indirect(r).FieldByName(groupedBy).String()
		_, present := groups[key]

		if !present {
			groups[key] = GroupedByData{GroupedBy: key, Cost: c.Cost, LifetimeCost: c.LifetimeCost, Clusters: 1}
		} else {
			g := groups[key]
			g.Cost += c.Cost
			g.LifetimeCost += c.LifetimeCost
			g.Clusters += 1
			groups[key] = g
		}
	}

	for _, g := range groups {
		data = append(data, g)
	}
	return data
}

func ConnectToDb() *gorm.DB {

	////##################### HEROKU CONFIGS  ####################
	//############################################################
	// HOST := "ec2-176-34-215-248.eu-west-1.compute.amazonaws.com"
	// DB := "dd5sa470g55vms"
	// PORT := "5432"
	// PASSWORD := "7fc395047bf2fe8a7415c2c3da94f3ffada9acec10f10bf1603bf96171bc9aa1"
	// USER := "efzyknuynecnce"
	dsn := "postgres://efzyknuynecnce:7fc395047bf2fe8a7415c2c3da94f3ffada9acec10f10bf1603bf96171bc9aa1@ec2-176-34-215-248.eu-west-1.compute.amazonaws.com:5432/dd5sa470g55vms"
	//############################################################

	//dsn := "host=localhost user=pchoudhary password=password dbname=pchoudhary port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&Instance{})
	//db.Migrator().CreateConstraint(&Cluster{}, "clusterKey")
	return db
}

func GetAllInstancesFromDB(db *gorm.DB) []Instance {
	var instancesReceived []Instance
	result := db.Find(&instancesReceived)
	log.Infof("Rows fetched from DB %v", result.RowsAffected)
	return instancesReceived
}

func GetClustersFromDB(db *gorm.DB, sd string, ed string) []Clusters {
	var clustersReceived []Clusters
	var result *gorm.DB

	startDate, err := time.Parse("2006-01-02", sd)
	if err != nil {
		log.Fatal(err)
	}
	endDate, err := time.Parse("2006-01-02", ed)
	if err != nil {
		log.Fatal(err)
	}

	result = db.Model(&Instance{}).Select(
		"cluster_name, cluster_state, cluster_type, owner, region, account, manager, l4, sum(cost) as cost, sum(estimated_cost) as estimated_cost, sum(lifetime_cost) as lifetime_cost, avg(cpu_usage) as cpu_usage, max(max_cpu_usage) as max_cpu_usage, sum(execution_time) as execution_time").Where(
		"launch_time > ? AND exit_time < ?", startDate, endDate).Group(
		"cluster_name, cluster_state, cluster_type, owner, region, account, manager, l4").Find(&clustersReceived)

	log.Infof("Instance Rows fetched from DB %v", result.RowsAffected)
	return clustersReceived
}
