CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL UNIQUE,
    hashed_password VARCHAR(60) NOT NULL
);

CREATE TABLE IF NOT EXISTS user_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    encrypted_key BYTEA NOT NULL,
    salt BYTEA NOT NULL,
    key_id INTEGER GENERATED ALWAYS AS IDENTITY,
    is_active BOOLEAN,
    created_at TIMESTAMP NOT NULL,
    CONSTRAINT user_keys_fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT user_keys_unique_constraint_user_key_id UNIQUE (user_id, key_id)
);

CREATE TABLE IF NOT EXISTS logins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    number INTEGER NOT NULL,
    url VARCHAR(255),
    description TEXT,
    hashed_login BYTEA,
    hashed_password BYTEA,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    nonce BYTEA NOT NULL,
    key_id INTEGER NOT NULL,
    CONSTRAINT logins_fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT logins_unique_constraint_user_title UNIQUE (user_id, title),
    CONSTRAINT logins_unique_constraint_user_number UNIQUE (user_id, number)
);

CREATE TABLE IF NOT EXISTS user_record_number (
    user_id UUID NOT NULL UNIQUE,
    current_number INTEGER,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    number INTEGER NOT NULL,
    hashed_pan BIGINT,
    hashed_holder VARCHAR(255),
    hashed_expiry_date VARCHAR(4),
    hashed_cvv VARCHAR(4),
    hashed_pin VARCHAR(4),
    bank VARCHAR(255),
    brand VARCHAR(255),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    CONSTRAINT cardss_fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT cards_unique_constraint_user_title UNIQUE (user_id, title),
    CONSTRAINT cards_unqiue_constraint_user_number UNIQUE (user_id, number),
    CONStRAINT cards_unqiue_constraint_user_pan UNIQUE (user_id, hashed_pan)
);

CREATE TABLE IF NOT EXISTS texts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    number INTEGER NOT NULL,
    hashed_text VARCHAR(255),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    CONSTRAINT texts_fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT texts_unique_constraint_user_title UNIQUE (user_id, title),
    CONSTRAINT texts_unique_constraint_user_number UNIQUE (user_id, number)
)
