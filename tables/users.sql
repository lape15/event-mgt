-- CREATE TABLE IF NOT EXISTS users (
--     user_id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
--     username TEXT NOT NULL,
--     email VARCHAR(100) NOT NULL,
--     password VARCHAR(255) NOT NULL,
--     full_name VARCHAR(100),
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
--     role VARCHAR(50) DEFAULT 'user'
--     )

-- ALTER TABLE users ADD COLUMN role VARCHAR(50) DEFAULT 'creator';

CREATE TABLE IF NOT EXISTS users (
    user_id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    role VARCHAR(50) DEFAULT 'creator'
);
