
CREATE TABLE IF NOT EXISTS roles (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  username TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  provider TEXT,
  password TEXT,
  reset_password_token TEXT,
  confirmation_token TEXT,
  confirmed BOOLEAN DEFAULT FALSE,
  blocked BOOLEAN DEFAULT FALSE,
  role_id INT REFERENCES roles(id),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS licenses (
  id SERIAL PRIMARY KEY,
  photo TEXT,
  name TEXT,
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS main_carusel (
  id SERIAL PRIMARY KEY,
  image TEXT,
  header TEXT,
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS specializations (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS schedules (
  id SERIAL PRIMARY KEY,
  day INTEGER NOT NULL CHECK (day BETWEEN 1 AND 7),
  time_from TIME NOT NULL,
  time_to TIME NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS service_categories (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  category_photo TEXT,
  favorite BOOLEAN DEFAULT FALSE,
  specialization_id INT UNIQUE REFERENCES specializations(id),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS services (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  specific_photo TEXT,
  price INTEGER,
  service_category_id INT REFERENCES service_categories(id),
  specialization_id INT REFERENCES specializations(id),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS doctors (
  id SERIAL PRIMARY KEY,
  fullname TEXT NOT NULL,
  description TEXT,
  doctor_photo TEXT,
  schedule_id INT UNIQUE REFERENCES schedules(id),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS doctor_specializations (
  doctor_id INT REFERENCES doctors(id) ON DELETE CASCADE,
  specialization_id INT REFERENCES specializations(id) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (doctor_id, specialization_id)
);


CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id);
CREATE INDEX IF NOT EXISTS idx_services_category_id ON services(service_category_id);
CREATE INDEX IF NOT EXISTS idx_services_specialization_id ON services(specialization_id);
CREATE INDEX IF NOT EXISTS idx_service_categories_specialization_id ON service_categories(specialization_id);
CREATE INDEX IF NOT EXISTS idx_schedules_day ON schedules(day);
CREATE INDEX IF NOT EXISTS idx_doctor_specializations_doctor_id ON doctor_specializations(doctor_id);
CREATE INDEX IF NOT EXISTS idx_doctor_specializations_specialization_id ON doctor_specializations(specialization_id);

-- Insert default roles
-- INSERT INTO roles (name) VALUES 
--     ('admin'),
--     ('doctor'),
--     ('user')
-- ON CONFLICT (name) DO NOTHING;