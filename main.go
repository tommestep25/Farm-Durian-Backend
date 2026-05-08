package main

import (
    "os"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // 1. เชื่อมต่อ Database (ฟังก์ชันจาก database.go)
    initDB()

    // 2. ตั้งค่า CORS
    r.Use(cors.New(cors.Config{
        AllowOrigins: []string{
            "http://localhost:5173",
            "https://farm-durian.netlify.app",
        },
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    // 3. สร้าง API Route (ต้องประกาศให้เสร็จก่อนสั่ง Run)
    api := r.Group("/api")
    {
        api.GET("/trees", getTrees)
        api.POST("/trees", createTree)
        api.GET("/dashboard/stats", getDashboardStats)
        api.POST("/fertilizer", addFertilizer)
        api.GET("/fertilizer/history", getFertilizerHistory)
        api.GET("/fertilizer/inventory", getFertilizerInventory)
        api.POST("/fertilizer/inventory", addFertilizerToInventory)
        api.PUT("/fertilizer/inventory/:id", updateFertilizerStock)
        api.POST("/pest-disease", addPestLog)
        api.GET("/pest-disease/history", getPestHistory)
        api.GET("/reports/summary", getReportSummary)
        api.GET("/reports/monthly-usage", getMonthlyUsage)
    }

    // 4. เตรียมค่า Port
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // 5. สั่งรัน Server ไว้ที่บรรทัดสุดท้ายเพียงที่เดียว
    r.Run(":" + port)
}