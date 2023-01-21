-- Create(新增欄位資料，可多筆新增)
-- 未定義欄位資訊，須依照「表建立時的欄位順序依序填寫」
-- 因此即使第一欄 id 是 AUTO_INCREMENT 也必須填寫...
-- 故可使用下方第二種(有定義欄位資訊)的方式，進行新增
INSERT INTO table_name 
VALUES
  (3, "test3",22);

-- 有定義欄位資訊
INSERT INTO table_name (name, age)
VALUES
  ("test1", 20),
  ("test2", 21),
  ("test3", 22);

-- Read(取得欄位資訊)
SELECT * FROM table_name;
SELECT name, age FROM table_name;
SELECT name, age FROM table_name WHERE id='1';

-- Update(更新欄位資料，可多欄更新)
UPDATE table_name SET name='test4', age='23' WHERE id=1;

-- Delete(刪除欄位資料)
DELETE FROM table_name WHERE name='test4';