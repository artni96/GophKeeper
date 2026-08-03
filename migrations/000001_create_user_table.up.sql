CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL UNIQUE,
    hashed_password VARCHAR(60) NOT NULL
);

CREATE TABLE IF NOT EXISTS logins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    user_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    number INTEGER NOT NULL,
    url VARCHAR(255),
    description TEXT,
    hashed_login VARCHAR(255) NOT NULL,
    hashed_password VARCHAR (60) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT constraint_unique_user_title UNIQUE (user_id, title),
    CONSTRAINT constraint_unqiue_user_number UNIQUE (user_id, number)
);

CREATE TABLE IF NOT EXISTS user_record_number (
    user_id UUID NOT NULL,
    current_number INTEGER,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);