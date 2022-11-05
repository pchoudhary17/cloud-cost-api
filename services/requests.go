package services

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	action "github.infra.cloudera.com/computeops/cloud-cost-dashboard-api/actions"
	model "github.infra.cloudera.com/computeops/cloud-cost-dashboard-api/models"
)

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

	dbHandler := action.ConnectToDb()
	clustersFromDB := action.GetClustersFromDB(dbHandler, beginDate, endDate)
	if groupBy == "none" {
		c.JSON(http.StatusOK, clustersFromDB)
	} else {
		c.JSON(http.StatusOK, groupData(clustersFromDB, groupBy))
	}
}

func groupData(clusters []model.Clusters, groupedBy string) []model.GroupedByData {
	var data []model.GroupedByData

	groups := make(map[string]model.GroupedByData)
	for _, c := range clusters {
		r := reflect.ValueOf(c)
		key := reflect.Indirect(r).FieldByName(groupedBy).String()
		_, present := groups[key]

		if !present {
			groups[key] = model.GroupedByData{GroupedBy: key, Cost: c.Cost, LifetimeCost: c.LifetimeCost, Clusters: 1}
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
