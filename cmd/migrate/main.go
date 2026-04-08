package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/lite-cms/cms/internal/config"
	"github.com/lite-cms/cms/internal/pkg/database"
)

func main() {
	// 1. 加载配置
	cfg := config.Load()

	// 2. 连接数据库
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	log.Println("数据库连接成功")

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取 SQL DB 失败: %v", err)
	}

	// 3. 读取 migrations 目录下的所有 .up.sql 文件
	files, err := os.ReadDir("migrations")
	if err != nil {
		log.Fatalf("读取 migrations 目录失败: %v", err)
	}

	var migrationFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".up.sql") {
			migrationFiles = append(migrationFiles, f.Name())
		}
	}

	// 4. 按文件名排序确保顺序执行
	sort.Strings(migrationFiles)

	log.Printf("发现 %d 个迁移文件，准备执行...", len(migrationFiles))

	// 5. 逐个执行 SQL
	for _, filename := range migrationFiles {
		log.Printf("执行迁移: %s ...", filename)
		
		path := filepath.Join("migrations", filename)
		content, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("读取文件 %s 失败: %v", filename, err)
		}

		// 执行 SQL 语句
		// 注意：对于复杂的 SQL，可能需要按分号拆分执行，
		// 但这里的 SQL 比较标准，直接通过 Exec 执行整个文件内容
		_, err = sqlDB.Exec(string(content))
		if err != nil {
			// 如果是 "already exists" 错误，通常可以忽略（取决于 SQL 里的 IF NOT EXISTS）
			if strings.Contains(err.Error(), "already exists") {
				log.Printf("跳过 %s: 部分对象已存在", filename)
				continue
			}
			log.Fatalf("执行 %s 失败: %v", filename, err)
		}
		
		log.Printf("迁移成功: %s", filename)
	}

	fmt.Println("\n所有数据库迁移任务已完成！")
}
