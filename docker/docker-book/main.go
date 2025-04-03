package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type Task struct {
    Id          string `json:"id" binding:"required"`
    Name        string `json:"name" binding:"required"`
    Description string `json:"description"`
    Timestamp   int64  `json:"timestamp"`
}

var (
    client = redis.NewClient(&redis.Options{
        Addr:     getStrEnv("REDIS_HOST", "localhost:6379"),
        Password: getStrEnv("REDIS_PASSWORD", ""),
        DB:       getIntEnv("REDIS_DB", 0),
    })
)

func setupRouter() *gin.Engine {
    r := gin.Default()

    r.GET("/ping3131313", func(c *gin.Context) {
        if _, err := client.Ping().Result(); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis connection failed"})
            return
        }
        c.String(http.StatusOK, "pong")
    })

    r.GET("/task", func(c *gin.Context) {
        tasks := []Task{}
        keys, err := client.Keys("task:*").Result()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
            return
        }
        for _, key := range keys {
            if taskJSON, err := client.Get(key).Result(); err == nil {
                var task Task
                if json.Unmarshal([]byte(taskJSON), &task) == nil {
                    tasks = append(tasks, task)
                }
            }
        }
        c.JSON(http.StatusOK, gin.H{"tasks": tasks})
    })

    r.GET("/task/:id", func(c *gin.Context) {
        id := c.Param("id")
        taskJSON, err := client.Get("task:" + id).Result()
        if err == redis.Nil {
            c.JSON(http.StatusNotFound, gin.H{"id": id, "message": "Task not found"})
            return
        } else if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task"})
            return
        }
        var task Task
        if err := json.Unmarshal([]byte(taskJSON), &task); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse task"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"task": task})
    })

    r.POST("/task", func(c *gin.Context) {
        var task Task
        if err := c.ShouldBindJSON(&task); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task format", "details": err.Error()})
            return
        }
        taskJSON, _ := json.Marshal(task)
        if err := client.Set("task:"+task.Id, taskJSON, 0).Err(); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save task"})
            return
        }
        c.JSON(http.StatusCreated, gin.H{"task": task, "created": true, "message": "Task created successfully"})
    })

    r.PUT("/task/:id", func(c *gin.Context) {
        id := c.Param("id")
        var task Task
        if err := c.ShouldBindJSON(&task); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task format", "details": err.Error()})
            return
        }
        task.Id = id
        taskJSON, _ := json.Marshal(task)
        if err := client.Set("task:"+id, taskJSON, 0).Err(); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"task": task, "updated": true, "message": "Task updated successfully"})
    })

    r.DELETE("/task/:id", func(c *gin.Context) {
        id := c.Param("id")
        if err := client.Del("task:" + id).Err(); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"id": id, "message": "Task deleted successfully"})
    })

    return r
}
 func waitForRedis(client *redis.Client) {
    for {
        _, err := client.Ping().Result()
        if err == nil {
            break
        }
        time.Sleep(1 * time.Second)
    }
}
func main() {
    waitForRedis(client)
    r := setupRouter()
    r.Run(getStrEnv("TASK_MANAGER_HOST", ":8080"))
}

func getIntEnv(key string, defaultValue int) int {
    value := os.Getenv(key)
    if len(value) == 0 {
        return defaultValue
    }
    if i, err := strconv.Atoi(value); err == nil {
        return i
    }
    return defaultValue
}

func getStrEnv(key string, defaultValue string) string {
    value := os.Getenv(key)
    if len(value) == 0 {
        return defaultValue
    }
    return value
}