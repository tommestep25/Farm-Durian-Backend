func main() {
    r := gin.Default()

    initDB()

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

    // --- 1. เตรียมค่า Port (แค่ประกาศตัวแปร ห้ามสั่ง r.Run ตรงนี้!) ---
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // --- 2. สร้าง API Route (ต้องอยู่ก่อน r.Run) ---
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

    // --- 3. สั่งรัน Server ไว้ที่บรรทัดสุดท้ายของฟังก์ชัน main เพียงที่เดียว! ---
    r.Run(":" + port)
}