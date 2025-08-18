package handlers

import (
	"net/http"
	"os"
	"expense-control-service/configs"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func Alive(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "alive",
		"message": "Service is up and running",
	})
}

func Status(c *gin.Context) {
	mysqlStatus := "healthy"
	if configs.GlobalConfig.MySqlLatency >= 100 {
		mysqlStatus = "unhealthy"
	}

	arrayStatus := []string{
		mysqlStatus,
	}

	applicationStatus := "healthy"
	for _, status := range arrayStatus {
		{
			if status == "unhealthy" {
				applicationStatus = "unhealthy"
			}
		}
	}

	cpuUsage, _ := cpu.Percent(0, false)
	memory, _ := mem.VirtualMemory()

	response := gin.H{
		"name":      os.Getenv("APP_NAME"),
		"status":    applicationStatus,
		"version":   "1.0",
		"timestamp": time.Now().Format(time.RFC3339),
		"system": gin.H{
			"cpu": gin.H{
				"utilization": cpuUsage,
			},
			"memory": gin.H{
				"total": memory.Total / 1024 / 1024,
				"used":  memory.Used / 1024 / 1024,
			},
		},
	}

	response["dependencies"] = []gin.H{
		{
			"name":     "MySQL",
			"kind":     "mysql",
			"status":   mysqlStatus,
			"latency":  configs.GlobalConfig.MySqlLatency,
			"optional": false,
			"internal": true,
		},
	}

	c.JSON(http.StatusOK, response)
}

func RedirectToAlive(c *gin.Context) {
	c.Redirect(http.StatusFound, "/health-check/alive")
}
