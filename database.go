package main

import (
	"log"
	"os"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func initDB() {
    // dsn := os.Getenv("DATABASE_URL")
	dsn == ""
    if dsn == "" {
        // ใช้ Port 6543 และชื่อ Host แบบ Pooler พร้อมรหัสผ่านใหม่ของคุณ
        // สังเกต: ชื่อ User ต้องมี .otnuuwigqsfpkrfdvgpj ต่อท้ายด้วยครับ
        dsn = "postgres://postgres.otnuuwigqsfpkrfdvgpj:Luckygamer144@aws-0-ap-southeast-1.pooler.supabase.com:6543/postgres"
    }
    
    var err error
    // ใช้ sqlx.Open แทน Connect เพื่อไม่ให้มันค้างตอนเริ่มต้น
    db, err = sqlx.Open("postgres", dsn)
    if err != nil {
        log.Fatalln("❌ Database Open Error:", err)
    }

    // ตั้งค่าเวลาเชื่อมต่อให้ไม่รอนานเกินไป
    db.SetMaxOpenConns(10)
    
    // ลอง Ping ดูว่าใช้ได้จริงไหม
    if err = db.Ping(); err != nil {
        log.Fatalln("❌ Database Ping Error (ต่อไม่ได้):", err)
    }

    log.Println("✅ เชื่อมต่อ PostgreSQL (Supabase) สำเร็จ!")
}