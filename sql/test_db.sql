CREATE DATABASE test_snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE USER 'test_web'@'%' IDENTIFIED BY 'pass';

GRANT CREATE, DROP, ALTER, INDEX, SELECT, UPDATE, DELETE, INSERT ON test_snippetbox.* TO 'test_web'@'%';
