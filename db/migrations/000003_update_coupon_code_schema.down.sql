-- delete new version
DROP TABLE IF EXISTS coupon_codes;
DROP TABLE IF EXISTS files;

-- restore old coupon codes table and files table
CREATE TABLE coupon_code_files
(
    id         SERIAL PRIMARY KEY,
    file_name  TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE coupon_codes
(
    id         SERIAL PRIMARY KEY,
    file_id    INTEGER NOT NULL,
    code       TEXT    NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
