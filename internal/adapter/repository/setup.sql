CREATE TABLE IF NOT EXISTS users (
    user_id   SERIAL PRIMARY KEY,
    user_name VARCHAR(100) NOT NULL UNIQUE,
    password  VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
    session_id SERIAL PRIMARY KEY,
    user_id    INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    token      VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS habits (
    habit_id     SERIAL PRIMARY KEY,
    user_id      INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    habit_name   VARCHAR(100) NOT NULL,
    category     VARCHAR(50),
    target_value INT,
    unit         VARCHAR(20),
    p_cost       INT CHECK (p_cost BETWEEN 1 AND 10),
    m_cost       INT CHECK (m_cost BETWEEN 1 AND 10),
    must         BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE (user_id, habit_name)
);

CREATE TABLE IF NOT EXISTS conditions (
    condition_id SERIAL PRIMARY KEY,
    user_id      INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    p_condition  INT NOT NULL CHECK (p_condition BETWEEN 1 AND 5),
    m_condition  INT NOT NULL CHECK (m_condition BETWEEN 1 AND 5),
    date         DATE NOT NULL,
    UNIQUE (user_id, date)
);

CREATE TABLE IF NOT EXISTS logs (
    log_id       SERIAL PRIMARY KEY,
    habit_id     INT NOT NULL REFERENCES habits(habit_id) ON DELETE CASCADE,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    date         DATE NOT NULL,
    UNIQUE (habit_id, date)
);
