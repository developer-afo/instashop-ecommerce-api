-- remove category_id column from products
ALTER TABLE products DROP COLUMN category_id;

-- drop categories table
DROP TABLE categories;