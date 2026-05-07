package main

import (
	"log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func initDB() {
	// แทนที่บรรทัดด้านล่างด้วย Connection String ของคุณ
	dsn := os.Getenv("DATABASE_URL")

    // 2. ถ้าในเครื่องไม่มี DATABASE_URL (เช่นตอนรัน Local) ให้ใช้ค่าเดิมที่คุณมี
    if dsn == "" {
        dsn = "postgres://postgres:%23Lovelove144@db.otnuuwigqsfpkrfdvgpj.supabase.co:5432/postgres"
    }
	
	
	var err error
	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalln("❌ เชื่อมต่อ Database ล้มเหลว:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalln("❌ Database ไม่ตอบสนอง:", err)
	}

	log.Println("✅ เชื่อมต่อ PostgreSQL (Supabase) สำเร็จ!")
}