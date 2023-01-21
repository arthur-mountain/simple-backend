-- 新增表
CREATE TABLE table_name(
  id INTEGER PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  age INTEGER NOT NULL
  -- PRIMARY KEY(id)
);

-- 新增表中的欄位
ALTER TABLE table_name ADD time DATETIME NOT NULL;

-- 刪除表中的欄位
ALTER TABLE table_name DROP time;

-- 修改表中的欄位，且預設當下 timestamp
ALTER TABLE table_name MODIFY time TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- 刪除表
DROP TABLE table_name;

-- 顯示全部的表
SHOW TABLES;

-- 顯示表的欄位資訊
DESCRIBE table_name;