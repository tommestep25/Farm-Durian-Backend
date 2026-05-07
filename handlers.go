package main

import (
	"log"
    "os"
	"net/http"
	"github.com/gin-gonic/gin"
)

type Tree struct {
	ID           int    `db:"id" json:"id"`
	TreeID       string `db:"tree_id" json:"treeId"`
	Species      string `db:"species" json:"species"`
	PlantDate    string `db:"plant_date" json:"plantDate"` // รับจาก TO_CHAR ใน SQL
	Status       string `db:"status" json:"status"`
	CurrentStage string `db:"current_stage" json:"currentStage"`
}

func getTrees(c *gin.Context) {
	trees := []Tree{}
	// ใช้ TO_CHAR เพื่อแปลงประเภท 'date' ใน Postgres ให้เป็น string ที่ Go เข้าใจง่าย
	query := `
		SELECT 
			id, 
			tree_id, 
			species, 
			TO_CHAR(plant_date, 'YYYY-MM-DD') as plant_date, 
			status, 
			current_stage 
		FROM public.trees 
		ORDER BY tree_id ASC`
	
	err := db.Select(&trees, query)
	if err != nil {
		log.Println("❌ Database Query Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูลได้: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, trees)
}

func createTree(c *gin.Context) {
	var newTree Tree
	if err := c.ShouldBindJSON(&newTree); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ถูกต้อง"})
		return
	}

	query := `
		INSERT INTO public.trees (tree_id, species, plant_date, status, current_stage) 
		VALUES (:tree_id, :species, :plant_date, :status, :current_stage)`
	
	_, err := db.NamedExec(query, newTree)
	if err != nil {
		log.Println("❌ Insert Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "บันทึกข้อมูลล้มเหลว: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "ลงทะเบียนต้นไม้สำเร็จ"})
}

type FertilizerLog struct {
	ID        int     `db:"id" json:"id"`
	TreeID    string  `db:"tree_id" json:"treeId"`
	Date      string  `db:"date" json:"date"`
	Formula   string  `db:"formula" json:"formula"`
	Amount    float64 `db:"amount" json:"amount"`
	Target    string  `db:"target" json:"target"`
}

// ฟังก์ชันบันทึกการใส่ปุ๋ย
func addFertilizer(c *gin.Context) {
	var logData FertilizerLog
	if err := c.ShouldBindJSON(&logData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ถูกต้อง"})
		return
	}

	query := `
		INSERT INTO public.fertilizer_logs (tree_id, date, formula, amount, target) 
		VALUES (:tree_id, :date, :formula, :amount, :target)`
	
	_, err := db.NamedExec(query, logData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "บันทึกปุ๋ยล้มเหลว: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "บันทึกข้อมูลปุ๋ยสำเร็จ"})
}

type PestDiseaseLog struct {
    ID          int     `db:"id" json:"id"`
    TreeID      string  `db:"tree_id" json:"treeId"`
    FoundDate   string  `db:"found_date" json:"foundDate"`
    Symptom     string  `db:"symptom" json:"symptom"`
    DiseaseName string  `db:"disease_name" json:"diseaseName"`
    Treatment   string  `db:"treatment" json:"treatment"`
    Medicine    string  `db:"medicine" json:"medicine"`
    Status      string  `db:"status" json:"status"`
    ImageUrl    *string `db:"image_url" json:"imageUrl"` // ✅ เติม * หน้า string
}

// โครงสร้างสำหรับรายงานสรุป (Summary Report)
type SummaryReport struct {
    TreeID          string  `db:"tree_id" json:"treeId"`
    Species         string  `db:"species" json:"species"`
    // แก้ไข 2 บรรทัดข้างล่างนี้ให้เติม _kg ตามในรูป
    TotalFertilizer float64 `db:"total_fertilizer_kg" json:"totalFertilizer"`
    LatestYield     float64 `db:"latest_yield_kg" json:"latestYield"`
    HealthStatus    string  `db:"health_status" json:"healthStatus"`
}

// 1. ฟังก์ชันบันทึกการพบโรคและแมลง
func addPestLog(c *gin.Context) {
    var logData PestDiseaseLog
    if err := c.ShouldBindJSON(&logData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ถูกต้อง"})
        return
    }

    // เพิ่ม image_url และ :image_url เข้าไปใน Query
    query := `
        INSERT INTO public.pest_disease_logs (
            tree_id, 
            found_date, 
            symptom, 
            disease_name, 
            treatment, 
            medicine, 
            status,
            image_url
        ) 
        VALUES (
            :tree_id, 
            :found_date, 
            :symptom, 
            :disease_name, 
            :treatment, 
            :medicine, 
            :status,
            :image_url
        )`
    
    _, err := db.NamedExec(query, logData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "บันทึกข้อมูลโรคล้มเหลว: " + err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"message": "บันทึกข้อมูลโรค/แมลงสำเร็จ"})
}

// 2. ฟังก์ชันดึงรายงานสรุป (จาก View farm_summary_report ที่เราสร้างไว้ใน Postgres)
func getReportSummary(c *gin.Context) {
	reports := []SummaryReport{}
	// ดึงข้อมูลจาก View ที่เราเคยเขียนไว้ใน SQL Editor
	err := db.Select(&reports, "SELECT * FROM farm_summary_report ORDER BY tree_id ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ดึงรายงานล้มเหลว: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, reports)
}
func getFertilizerHistory(c *gin.Context) {
	treeID := c.Query("treeId")
	logs := []FertilizerLog{}
	
	// สำหรับ PostgreSQL ต้องใช้ $1 แทน ? 
	query := `SELECT id, tree_id, TO_CHAR(date, 'YYYY-MM-DD') as date, formula, amount, target 
	          FROM public.fertilizer_logs 
	          WHERE tree_id = $1 
	          ORDER BY date DESC`
	
	err := db.Select(&logs, query, treeID)
	if err != nil {
		log.Println("❌ History Query Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}

// ดึงประวัติโรคเฉพาะต้น
func getPestHistory(c *gin.Context) {
    treeID := c.Query("treeId")
    logs := []PestDiseaseLog{}
    
    // เพิ่ม image_url เข้าไปใน Query หลัง medicine, status
    query := `SELECT id, tree_id, TO_CHAR(found_date, 'YYYY-MM-DD') as found_date, 
              symptom, disease_name, treatment, medicine, status, image_url 
              FROM public.pest_disease_logs 
              WHERE tree_id = $1 
              ORDER BY found_date DESC`
    
    err := db.Select(&logs, query, treeID)
    if err != nil {
        log.Println("❌ Pest History Error:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, logs)
}

type DashboardStats struct {
	TotalTrees int `db:"total_trees" json:"totalTrees"`
	Normal     int `db:"normal_count" json:"normal"`
	Risk       int `db:"risk_count" json:"risk"`
	Urgent     int `db:"urgent_count" json:"urgent"`
}

func getDashboardStats(c *gin.Context) {
	var stats DashboardStats
	query := `
		SELECT 
			COUNT(*) as total_trees,
			COUNT(*) FILTER (WHERE status = 'NORMAL') as normal_count,
			COUNT(*) FILTER (WHERE status = 'RISK') as risk_count,
			COUNT(*) FILTER (WHERE status = 'URGENT') as urgent_count
		FROM public.trees`
	
	err := db.Get(&stats, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

type FertilizerInventory struct {
    ID            int     `db:"id" json:"id"`
    BrandName     string  `db:"brand_name" json:"brandName"`
    Formula       string  `db:"formula" json:"formula"`
    PricePerKg    float64 `db:"price_per_kg" json:"pricePerKg"`
    StockQuantity float64 `db:"stock_quantity" json:"stockQuantity"`
}

// ฟังก์ชันดึงรายการปุ๋ยทั้งหมด
func getFertilizerInventory(c *gin.Context) {
    items := []FertilizerInventory{}
    // ตรวจสอบว่ามี public. นำหน้าชื่อตาราง
    err := db.Select(&items, "SELECT id, brand_name, formula, price_per_kg, stock_quantity FROM public.fertilizer_inventory ORDER BY brand_name")
    if err != nil {
        log.Println("❌ Inventory Query Error:", err) // เพิ่มบรรทัดนี้เพื่อดู Error ใน Terminal
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, items)
}
func addFertilizerToInventory(c *gin.Context) {
    var newItem FertilizerInventory
    if err := c.ShouldBindJSON(&newItem); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ถูกต้อง"})
        return
    }

    query := `
        INSERT INTO public.fertilizer_inventory (brand_name, formula, price_per_kg, stock_quantity) 
        VALUES (:brand_name, :formula, :price_per_kg, :stock_quantity)`
    
    _, err := db.NamedExec(query, newItem)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "เพิ่มข้อมูลคลังล้มเหลว: " + err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"message": "เพิ่มปุ๋ยเข้าคลังสำเร็จ"})
}

// 2. ฟังก์ชันอัปเดตสต็อกหรือราคา (เมื่อซื้อของเพิ่ม)
func updateFertilizerStock(c *gin.Context) {
    id := c.Param("id")
    var updateData FertilizerInventory
    if err := c.ShouldBindJSON(&updateData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ถูกต้อง"})
        return
    }

    query := `
        UPDATE public.fertilizer_inventory 
        SET price_per_kg = $1, stock_quantity = $2 
        WHERE id = $3`
    
    _, err := db.Exec(query, updateData.PricePerKg, updateData.StockQuantity, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "อัปเดตสต็อกล้มเหลว: " + err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "อัปเดตข้อมูลสำเร็จ"})
}
type InventorySummary struct {
    TotalKg       float64 `db:"total_kg" json:"totalKg"`
    TotalValue    float64 `db:"total_value" json:"totalValue"`
}
// โครงสร้างสำหรับสรุปยอดรายเดือน
type MonthlyUsage struct {
    Month          string  `db:"usage_month" json:"month"`
    TotalKg        float64 `db:"total_kg" json:"totalKg"`
    EstimatedCost  float64 `db:"estimated_cost" json:"estimatedCost"`
}

// 3. ฟังก์ชันสรุปยอดการใช้ปุ๋ยรายเดือน (กี่กิโล / กี่บาท)
func getMonthlyUsage(c *gin.Context) {
    var summary InventorySummary
    
    // คำนวณ SUM จากตาราง public.fertilizer_inventory ตรงๆ
    query := `
        SELECT 
            COALESCE(SUM(stock_quantity), 0) as total_kg,
            COALESCE(SUM(stock_quantity * price_per_kg), 0) as total_value
        FROM public.fertilizer_inventory`

    err := db.Get(&summary, query)
    if err != nil {
        log.Println("❌ Inventory Summary Error:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // ส่งข้อมูลกลับไปในรูปแบบที่ Frontend รอรับ (ใช้ชื่อ Field เดิมเพื่อให้ UI ไม่พัง)
    c.JSON(http.StatusOK, []gin.H{
        {
            "month":          "ปัจจุบัน",
            "totalKg":        summary.TotalKg,
            "estimatedCost":  summary.TotalValue,
        },
    })
}