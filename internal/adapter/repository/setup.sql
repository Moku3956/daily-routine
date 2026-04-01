CREATE TABLE IF NOT EXISTS habits (
    habit_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    habit_name VARCHAR(100) NOT NULL,
    category VARCHAR(50),
    value INT,
    unit VARCHAR(20),
    p_cost INT,
    m_cost INT,
    must BOOLEAN
);