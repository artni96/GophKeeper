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

    login_value BYTEA,
    login_nonce BYTEA,
    login_key_id INTEGER,

    password_value BYTEA,
    password_nonce BYTEA,
    password_key_id INTEGER,

    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,

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
    pan_value BYTEA,
    pan_nonce BYTEA,
    pan_key_id INTEGER,

    holder_value BYTEA,
    holder_nonce BYTEA,
    holder_key_id INTEGER,

    expiry_date_value BYTEA,
    expiry_date_nonce BYTEA,
    expiry_date_key_id INTEGER,

    cvv_value BYTEA,
    cvv_nonce BYTEA,
    cvv_key_id INTEGER,

    pin_value BYTEA,
    pin_nonce BYTEA,
    pin_key_id INTEGER,

    bank VARCHAR(255),
    brand VARCHAR(255),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    CONSTRAINT cardss_fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT cards_unique_constraint_user_title UNIQUE (user_id, title),
    CONSTRAINT cards_unqiue_constraint_user_number UNIQUE (user_id, number),
    CONStRAINT cards_unqiue_constraint_user_pan UNIQUE (user_id, pan_value)
);

CREATE TABLE IF NOT EXISTS texts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    number INTEGER NOT NULL,

    text_value BYTEA,
    text_nonce BYTEA,
    text_key_id INTEGER,

    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    CONSTRAINT texts_fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT texts_unique_constraint_user_title UNIQUE (user_id, title),
    CONSTRAINT texts_unique_constraint_user_number UNIQUE (user_id, number)
)
