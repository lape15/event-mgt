CREATE TABLE IF NOT EXISTS event_rsvps(
    event_id INT,
    user_id BIGINT, 
    PRIMARY KEY (event_id, user_id)
    FOREIGN KEY (event_id) REFERENCES events(event_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);