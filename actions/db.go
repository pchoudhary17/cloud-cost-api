package action

import (
	"time"

	log "github.com/sirupsen/logrus"
	model "github.infra.cloudera.com/computeops/cloud-cost-dashboard-api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
	db.AutoMigrate(&model.Instance{})
	//db.Migrator().CreateConstraint(&Cluster{}, "clusterKey")
	return db
}

func GetAllInstancesFromDB(db *gorm.DB) []model.Instance {
	var instancesReceived []model.Instance
	result := db.Find(&instancesReceived)
	log.Infof("Rows fetched from DB %v", result.RowsAffected)
	return instancesReceived
}

func GetClustersFromDB(db *gorm.DB, sd string, ed string) []model.Clusters {
	var clustersReceived []model.Clusters
	var result *gorm.DB

	startDate, err := time.Parse("2006-01-02", sd)
	if err != nil {
		log.Fatal(err)
	}
	endDate, err := time.Parse("2006-01-02", ed)
	if err != nil {
		log.Fatal(err)
	}

	result = db.Model(&model.Instance{}).Select(
		"cluster_name, cluster_state, cluster_type, owner, region, account, manager, l4, sum(cost) as cost, sum(estimated_cost) as estimated_cost, sum(lifetime_cost) as lifetime_cost, avg(cpu_usage) as cpu_usage, max(max_cpu_usage) as max_cpu_usage, sum(execution_time) as execution_time").Where(
		"launch_time > ? AND exit_time < ?", startDate, endDate).Group(
		"cluster_name, cluster_state, cluster_type, owner, region, account, manager, l4").Find(&clustersReceived)

	log.Infof("Instance Rows fetched from DB %v", result.RowsAffected)
	return clustersReceived
}
