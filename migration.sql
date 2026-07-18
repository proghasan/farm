-- Farm Management System Database Schema
-- MySQL 8.x+

CREATE DATABASE IF NOT EXISTS farm CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE farm;

CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    email VARCHAR(255) NULL UNIQUE,
    phone VARCHAR(20) NULL UNIQUE,
    username VARCHAR(100) NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    avatar VARCHAR(255) NULL,
    role ENUM('Owner','Manager','Veterinarian','Worker','Accountant') DEFAULT 'Worker',
    status ENUM('Active','Inactive','Suspended') DEFAULT 'Active',
    last_login_at TIMESTAMP NULL,
    remember_token VARCHAR(100) NULL,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE species (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    updated_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE breeds (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    species_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(150) NOT NULL,
    created_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    updated_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    CONSTRAINT fk_breed_species FOREIGN KEY (species_id) REFERENCES species(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE animals (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tag_no VARCHAR(50) NOT NULL UNIQUE,
    species_id BIGINT UNSIGNED NOT NULL,
    breed_id BIGINT UNSIGNED NOT NULL,
    gender ENUM('Male','Female') NOT NULL,
    birth_date DATE NULL,
    purchase_date DATE NULL,
    purchase_price DECIMAL(12,2) DEFAULT 0,
    current_weight DECIMAL(8,2) NULL,
    color VARCHAR(100) NULL,
    status ENUM('Healthy','Pregnant','Sick','Sold','Dead') DEFAULT 'Healthy',
    last_vaccine DATE NULL,
    remarks TEXT NULL,
    created_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    updated_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_animal_species FOREIGN KEY (species_id) REFERENCES species(id),
    CONSTRAINT fk_animal_breed FOREIGN KEY (breed_id) REFERENCES breeds(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE animal_weight_histories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    animal_id BIGINT UNSIGNED NOT NULL,
    weight DECIMAL(8,2) NOT NULL,
    record_date DATE NOT NULL,
    remarks TEXT NULL,
    created_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    updated_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    CONSTRAINT fk_weight_animal FOREIGN KEY (animal_id) REFERENCES animals(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE vaccines (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    species_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(150) NOT NULL,
    description TEXT NULL,
    dose VARCHAR(100) NULL,
    minimum_age_value INT UNSIGNED NOT NULL,
    minimum_age_unit ENUM('Day','Week','Month','Year') NOT NULL,
    interval_value INT NOT NULL,
    interval_unit ENUM('Day','Week','Month','Year') NOT NULL,
    is_repeatable BOOLEAN DEFAULT TRUE,
    created_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    updated_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    CONSTRAINT fk_vaccine_species FOREIGN KEY (species_id) REFERENCES species(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE animal_vaccinations (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    animal_id BIGINT UNSIGNED NOT NULL,
    vaccine_id BIGINT UNSIGNED NOT NULL,
    vaccination_date DATE NOT NULL,
    next_due_date DATE NULL,
    doctor_name VARCHAR(150) NULL,
    remarks TEXT NULL,
    created_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    updated_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    CONSTRAINT fk_av_animal FOREIGN KEY (animal_id) REFERENCES animals(id) ON DELETE CASCADE,
    CONSTRAINT fk_av_vaccine FOREIGN KEY (vaccine_id) REFERENCES vaccines(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE inventory_categories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(150) NOT NULL UNIQUE,
    created_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    updated_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE inventory_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    category_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(200) NOT NULL,
    sku VARCHAR(100) NULL,
    unit VARCHAR(50) NOT NULL,
    purchase_price DECIMAL(12,2) DEFAULT 0,
    selling_price DECIMAL(12,2) DEFAULT 0,
    created_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    updated_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    CONSTRAINT fk_inventory_category FOREIGN KEY (category_id) REFERENCES inventory_categories(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE inventory_transactions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    inventory_item_id BIGINT UNSIGNED NOT NULL,
    transaction_type ENUM('Purchase','Sale','Consumption','Adjustment','Return','Damage') NOT NULL,
    quantity DECIMAL(12,2) NOT NULL,
    unit_price DECIMAL(12,2) DEFAULT 0,
    transaction_date DATE NOT NULL,
    remarks TEXT NULL,
    created_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    updated_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    CONSTRAINT fk_inventory_transaction FOREIGN KEY (inventory_item_id) REFERENCES inventory_items(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE account_heads (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    type ENUM('Income','Expense') NOT NULL,
    name VARCHAR(150) NOT NULL,
    description TEXT NULL,
    created_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    updated_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE account_transactions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    account_head_id BIGINT UNSIGNED NOT NULL,
    transaction_date DATE NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    payment_method ENUM('Cash','Bank','Mobile Banking','Other') DEFAULT 'Cash',
    reference_no VARCHAR(100) NULL,
    description TEXT NULL,
    created_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    updated_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    CONSTRAINT fk_account_head FOREIGN KEY (account_head_id) REFERENCES account_heads(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Default admin user (password: admin123)
INSERT INTO users (name, email, username, password, role, status) VALUES (
    'Admin',
    'admin@farm.com',
    'admin',
    '$2a$10$dummy-hash-placeholder', -- replace with actual bcrypt hash of 'admin123'
    'Owner',
    'Active'
);
