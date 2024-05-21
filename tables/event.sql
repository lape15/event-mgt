
-- CREATE TABLE IF NOT EXISTS  events (
--     event_id INT AUTO_INCREMENT PRIMARY KEY,
--     name VARCHAR(255) NOT NULL,
--     description TEXT,
--     start TIMESTAMP NOT NULL,
--     end TIMESTAMP NOT NULL,
--     location VARCHAR(255) NOT NULL,
--     organizer_id BIGINT NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
--     event_limit INT,
--     FOREIGN KEY (organizer_id) REFERENCES users(user_id)
-- );


-- CREATE TABLE IF NOT EXISTS events (
--     event_id INT AUTO_INCREMENT PRIMARY KEY,
--     name VARCHAR(255) NOT NULL,
--     description TEXT,
--     start TIMESTAMP NOT NULL,
--     end TIMESTAMP NOT NULL,
--     location VARCHAR(255) NOT NULL,
--     event_limit INT,
--     organizer_id BIGINT NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--      updated_at TIMESTAMP,
--     FOREIGN KEY (organizer_id) REFERENCES users (user_id) 
-- );

CREATE TABLE IF NOT EXISTS events (
    event_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    start TIMESTAMP NOT NULL DEFAULT '0000-00-00 00:00:00',
    end TIMESTAMP NOT NULL DEFAULT '0000-00-00 00:00:00',
    location VARCHAR(255) NOT NULL,
    event_limit INT,
    organizer_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP,
    FOREIGN KEY (organizer_id) REFERENCES users (user_id) 
);

