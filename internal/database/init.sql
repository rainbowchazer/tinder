CREATE TABLE users (
    email TEXT PRIMARY KEY,
    password TEXT NOT NULL,
    username TEXT NOT NULL,
    age INT,
    photo TEXT
);

CREATE TABLE likes (
    user_email TEXT REFERENCES users(email),
    liked_email TEXT REFERENCES users(email),
    likes BOOLEAN,
    PRIMARY KEY (user_email, liked_email)
);

CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    user1 TEXT REFERENCES users(email),
    user2 TEXT REFERENCES users(email)
);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    sender_email TEXT REFERENCES users(email),
    receiver_email TEXT REFERENCES users(email),
    message TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
