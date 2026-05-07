package main

import (
	"log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// เรียกใช้งานฟังก์ชันเชื่อมต่อ DB ที่อยู่ใน database.go
	initDB()

	// ตั้งค่า CORS ให้ละเอียดขึ้นเพื่อรองรับ PUT/DELETE ในอนาคต
r.Use(cors.New(cors.Config{
        AllowOrigins: []string{
            "http://localhost:5173",          // สำหรับรันในเครื่องตัวเอง
            "https://farm-durian.netlify.app", // URL หน้าเว็บจริงของคุณ
        },
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

	port := os.Getenv("PORT") // ดึงค่า Port ที่ Server กำหนดมาให้
    if port == "" {
        port = "8080" // ถ้าไม่มี (เช่นรันในเครื่อง) ให้ใช้ 8080 เหมือนเดิม
    }

	// สร้าง API Route
	api := r.Group("/api")
	{
		// --- 🌳 จัดการต้นไม้ & Dashboard ---
		api.GET("/trees", getTrees)
		api.POST("/trees", createTree)
		api.GET("/dashboard/stats", getDashboardStats)

		// --- 🧪 ระบบใส่ปุ๋ย (ประวัติการใช้งาน) ---
		api.POST("/fertilizer", addFertilizer)
		api.GET("/fertilizer/history", getFertilizerHistory)

		// --- 📦 ระบบคลังปุ๋ย (ยี่ห้อ, ราคา, สต็อก) ---
		api.GET("/fertilizer/inventory", getFertilizerInventory)    // ดูรายการปุ๋ยทั้งหมด
		api.POST("/fertilizer/inventory", addFertilizerToInventory) // เพิ่มยี่ห้อปุ๋ยใหม่
		api.PUT("/fertilizer/inventory/:id", updateFertilizerStock) // อัปเดตราคาหรือจำนวนสต็อก

		// --- 🐛 โรค & แมลง ---
		api.POST("/pest-disease", addPestLog)
		api.GET("/pest-disease/history", getPestHistory)

		// --- 📈 รายงานสรุป ---
		api.GET("/reports/summary", getReportSummary)
		api.GET("/reports/monthly-usage", getMonthlyUsage) // 🔥 เพิ่มสำหรับสรุปยอดรายเดือน (กี่กิโล/กี่บาท)
	}

	// รัน Server ที่ Port 8080
	r.Run(":8080")
}