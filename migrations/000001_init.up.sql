CREATE TABLE managers (
    id SERIAL PRIMARY KEY,
    login VARCHAR NOT NULL UNIQUE,
    password VARCHAR NOT NULL,
    name VARCHAR,
    surname VARCHAR,
    phone VARCHAR,
    is_active BOOL,
    is_admin BOOL,
    full_access BOOL DEFAULT true
);

CREATE TABLE clients (
    id SERIAL PRIMARY KEY,
    manager_id INT,
    name VARCHAR,
    phone VARCHAR,
    email VARCHAR,
    new BOOL,
    results JSONB,
    in_archive BOOL,
    date BIGINT,
    FOREIGN KEY (manager_id) REFERENCES managers (id)
);

CREATE TABLE diagnostics (
    id SERIAL PRIMARY KEY,
    title VARCHAR,
    link VARCHAR
);