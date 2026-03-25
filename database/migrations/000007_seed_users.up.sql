-- Seed: default users (password: 12345679, bcrypt hashed)
INSERT INTO users (id, password, username, email, token, created_at, updated_at) VALUES
    ('fadel',  '$2a$10$3vpAR9B6R4/kp3aVxW/OputnakVHimZeq3exqRV/UDCFahM7UkCY2', 'fadel',  'fadel@wan.local',  '', UNIX_TIMESTAMP(NOW(3)) * 1000, UNIX_TIMESTAMP(NOW(3)) * 1000),
    ('nanang', '$2a$10$5ZuDScaSKHV5ZKmv0hnIvucCkpDwRFyxi3COlwML.VTQ7Jy0U31ru', 'nanang', 'nanang@wan.local', '', UNIX_TIMESTAMP(NOW(3)) * 1000, UNIX_TIMESTAMP(NOW(3)) * 1000),
    ('try',    '$2a$10$w5AG658Bi41QFzbPUsE/aOpwmIIYla6dmjGAFD32JysSE0wEN0xVe', 'try',    'try@wan.local',    '', UNIX_TIMESTAMP(NOW(3)) * 1000, UNIX_TIMESTAMP(NOW(3)) * 1000),
    ('daffa',  '$2a$10$smlLuNcd9eU7bzpbNe7VbO2miK3JMZimNWvfE.X/FmV4BDQw8lO1C', 'daffa',  'daffa@wan.local',  '', UNIX_TIMESTAMP(NOW(3)) * 1000, UNIX_TIMESTAMP(NOW(3)) * 1000),
    ('dani',   '$2a$10$2DFjPcL58ZWuqyVy4PH6WemCybnb5O95j32865DLic65tNVN50qja', 'dani',   'dani@wan.local',   '', UNIX_TIMESTAMP(NOW(3)) * 1000, UNIX_TIMESTAMP(NOW(3)) * 1000);
