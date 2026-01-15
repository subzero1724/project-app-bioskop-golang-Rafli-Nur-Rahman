-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create cinemas table
CREATE TABLE IF NOT EXISTS cinemas (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    location VARCHAR(255) NOT NULL,
    description TEXT,
    total_seats INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create seats table
CREATE TABLE IF NOT EXISTS seats (
    id SERIAL PRIMARY KEY,
    cinema_id INTEGER NOT NULL REFERENCES cinemas(id) ON DELETE CASCADE,
    seat_number VARCHAR(10) NOT NULL,
    row_number VARCHAR(5) NOT NULL,
    seat_type VARCHAR(20) DEFAULT 'regular',
    price DECIMAL(10, 2) NOT NULL,
    UNIQUE(cinema_id, seat_number)
);

-- Create payment_methods table
CREATE TABLE IF NOT EXISTS payment_methods (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create bookings table
CREATE TABLE IF NOT EXISTS bookings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    cinema_id INTEGER NOT NULL REFERENCES cinemas(id) ON DELETE CASCADE,
    seat_id INTEGER NOT NULL REFERENCES seats(id) ON DELETE CASCADE,
    booking_date DATE NOT NULL,
    booking_time TIME NOT NULL,
    payment_method_id INTEGER REFERENCES payment_methods(id),
    payment_status VARCHAR(20) DEFAULT 'pending',
    total_amount DECIMAL(10, 2) NOT NULL,
    booking_status VARCHAR(20) DEFAULT 'reserved',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(cinema_id, seat_id, booking_date, booking_time)
);

-- Create tokens table for session management
CREATE TABLE IF NOT EXISTS tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(500) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_cinema_date_time ON bookings(cinema_id, booking_date, booking_time);
CREATE INDEX idx_seats_cinema_id ON seats(cinema_id);
CREATE INDEX idx_tokens_user_id ON tokens(user_id);
CREATE INDEX idx_tokens_token ON tokens(token);

-- Insert sample cinemas
INSERT INTO cinemas (name, location, description, total_seats) VALUES
('Cinema XXI Grand Indonesia', 'Jakarta Pusat', 'Premium cinema with IMAX and Dolby Atmos', 150),
('Cinema XXI Pondok Indah', 'Jakarta Selatan', 'Luxurious cinema with VIP lounges', 120),
('CGV Blitz Pacific Place', 'Jakarta Selatan', 'Modern cinema with 4DX technology', 180),
('Cinepolis Lippo Mall Puri', 'Jakarta Barat', 'Family-friendly cinema with junior seats', 100);

-- Insert seats for Cinema 1 (XXI Grand Indonesia)
DO $$
DECLARE
    cinema_id INTEGER := 1;
    row_letter CHAR(1);
    seat_num INTEGER;
BEGIN
    FOR row_letter IN SELECT chr(i) FROM generate_series(65, 74) i LOOP -- A to J (10 rows)
        FOR seat_num IN 1..15 LOOP -- 15 seats per row
            INSERT INTO seats (cinema_id, seat_number, row_number, seat_type, price)
            VALUES (
                cinema_id,
                row_letter || seat_num,
                row_letter,
                CASE 
                    WHEN row_letter IN ('A', 'B') THEN 'vip'
                    ELSE 'regular'
                END,
                CASE 
                    WHEN row_letter IN ('A', 'B') THEN 75000.00
                    ELSE 50000.00
                END
            );
        END LOOP;
    END LOOP;
END $$;

-- Insert seats for Cinema 2 (XXI Pondok Indah)
DO $$
DECLARE
    cinema_id INTEGER := 2;
    row_letter CHAR(1);
    seat_num INTEGER;
BEGIN
    FOR row_letter IN SELECT chr(i) FROM generate_series(65, 72) i LOOP -- A to H (8 rows)
        FOR seat_num IN 1..15 LOOP
            INSERT INTO seats (cinema_id, seat_number, row_number, seat_type, price)
            VALUES (
                cinema_id,
                row_letter || seat_num,
                row_letter,
                CASE 
                    WHEN row_letter IN ('A', 'B', 'C') THEN 'vip'
                    ELSE 'regular'
                END,
                CASE 
                    WHEN row_letter IN ('A', 'B', 'C') THEN 80000.00
                    ELSE 55000.00
                END
            );
        END LOOP;
    END LOOP;
END $$;

-- Insert seats for Cinema 3 (CGV Blitz Pacific Place)
DO $$
DECLARE
    cinema_id INTEGER := 3;
    row_letter CHAR(1);
    seat_num INTEGER;
BEGIN
    FOR row_letter IN SELECT chr(i) FROM generate_series(65, 75) i LOOP -- A to K (11 rows)
        FOR seat_num IN 1..16 LOOP
            INSERT INTO seats (cinema_id, seat_number, row_number, seat_type, price)
            VALUES (
                cinema_id,
                row_letter || seat_num,
                row_letter,
                'regular',
                60000.00
            );
        END LOOP;
    END LOOP;
END $$;

-- Insert seats for Cinema 4 (Cinepolis Lippo Mall Puri)
DO $$
DECLARE
    cinema_id INTEGER := 4;
    row_letter CHAR(1);
    seat_num INTEGER;
BEGIN
    FOR row_letter IN SELECT chr(i) FROM generate_series(65, 74) i LOOP -- A to J (10 rows)
        FOR seat_num IN 1..10 LOOP
            INSERT INTO seats (cinema_id, seat_number, row_number, seat_type, price)
            VALUES (
                cinema_id,
                row_letter || seat_num,
                row_letter,
                CASE 
                    WHEN row_letter IN ('A') THEN 'junior'
                    ELSE 'regular'
                END,
                CASE 
                    WHEN row_letter IN ('A') THEN 35000.00
                    ELSE 45000.00
                END
            );
        END LOOP;
    END LOOP;
END $$;

-- Insert payment methods
INSERT INTO payment_methods (name, description, is_active) VALUES
('Credit Card', 'Visa, MasterCard, American Express', true),
('Debit Card', 'ATM cards from major banks', true),
('E-Wallet', 'GoPay, OVO, Dana, ShopeePay', true),
('Bank Transfer', 'Transfer via internet banking', true),
('Cash', 'Pay at the counter', true);
