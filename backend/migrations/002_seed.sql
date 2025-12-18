-- Seed-данные

-- Администратор (email: admin@example.com, пароль: admin123)
INSERT INTO users (email, name, password_hash, role)
VALUES (
    'admin@example.com',
    'Admin User',
    -- bcrypt hash для пароля "admin123"
    -- bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
    '$2a$10$Vq5M8kGkAPqf.ZPuG1UFXOqkqTDPH5bJ8VfG5l9uP6XUOZx3QVG1e',
    'admin'
)
ON CONFLICT (email) DO NOTHING;

-- Тестовые товары
INSERT INTO products (name, description, price, stock)
VALUES
    ('Go T-Shirt', 'Comfortable Go-branded T-shirt', 2500, 100),
    ('Go Mug', 'Ceramic mug with Go gopher', 1500, 50),
    ('Sticker Pack', 'Set of Go/gopher stickers', 500, 200)
ON CONFLICT DO NOTHING;


