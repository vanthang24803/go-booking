-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE catalogs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE
);

CREATE TABLE listings (
    id SERIAL PRIMARY KEY,
    landlord_id INT REFERENCES users(id) ON DELETE CASCADE,

    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    location VARCHAR(255) NOT NULL,
    guests INT NOT NULL CHECK (guests >= 1),
    beds INT NOT NULL CHECK (beds >= 0),
    baths INT NOT NULL CHECK (baths >= 0),
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    cleaning_fee DECIMAL(10, 2) NOT NULL CHECK (cleaning_fee >= 0),
    service_fee DECIMAL(10, 2) NOT NULL CHECK (service_fee >= 0),
    taxes DECIMAL(10, 2) NOT NULL CHECK (taxes >= 0),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bookings (
    id SERIAL PRIMARY KEY,
    listing_id INT REFERENCES listings(id) ON DELETE CASCADE,
    guest_id INT REFERENCES users(id) ON DELETE CASCADE,

    guests INT NOT NULL CHECK (guests >= 1),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    nights INT NOT NULL CHECK (nights >= 0),
    phone_number VARCHAR(255),
    message_the_host TEXT,


    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE price_details (
    id SERIAL PRIMARY KEY,
    booking_id INT UNIQUE REFERENCES bookings(id) ON DELETE CASCADE,

    total_home_price DECIMAL(10, 2) NOT NULL CHECK (total_home_price >= 0),
    cleaning_fee DECIMAL(10, 2) NOT NULL CHECK (cleaning_fee >= 0),
    service_fee DECIMAL(10, 2) NOT NULL CHECK (service_fee >= 0),
    taxes DECIMAL(10, 2) NOT NULL CHECK (taxes >= 0),
    total_price DECIMAL(10, 2) NOT NULL CHECK (total_price >= 0),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    booking_id INT UNIQUE REFERENCES bookings(id) ON DELETE CASCADE,

    is_successful BOOLEAN NOT NULL,
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE catalogs_listings (
    catalog_id INT REFERENCES catalogs(id) ON DELETE CASCADE,
    listing_id INT REFERENCES listings(id) ON DELETE CASCADE
);

CREATE TABLE photos (
   id SERIAL PRIMARY KEY,
   listing_id INT NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
   
   public_id VARCHAR(255) NOT NULL UNIQUE,
   url TEXT NOT NULL,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    listing_id INT REFERENCES listings(id) ON DELETE CASCADE,
    author_id INT REFERENCES users(id) ON DELETE CASCADE,
    booking_id INT REFERENCES bookings(id) ON DELETE CASCADE,

    rating INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment TEXT NOT NULL,
    is_published BOOLEAN DEFAULT TRUE,
    is_edited BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE reviews;
DROP TABLE photos;
DROP TABLE catalogs_listings;
DROP TABLE catalogs;
DROP TABLE listings;
DROP TABLE bookings;
DROP TABLE price_details;
DROP TABLE payments;
-- +goose StatementEnd
