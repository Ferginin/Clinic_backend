-- =============================================================
-- Тестовые данные для Медлайф
-- Запускать ПОСЛЕ init.sql
-- Требует расширение pgcrypto для хеширования паролей
--
-- Тестовые пароли:
--   admin@admin.ru     → admin
--   doctor*@medlife.ru → Doctor123
--   patient*@test.ru   → Patient123
-- =============================================================

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- -------------------------------------------------------------
-- Роли (на случай если init.sql ещё не выполнялся)
-- -------------------------------------------------------------
INSERT INTO roles (name) VALUES ('admin'), ('doctor'), ('user')
ON CONFLICT (name) DO NOTHING;

-- -------------------------------------------------------------
-- Пользователи
-- -------------------------------------------------------------

-- Администратор
INSERT INTO users (username, email, password, confirmed, role_id) VALUES
  ('Администратор', 'admin@admin.ru',
   crypt('admin', gen_salt('bf', 10)), true,
   (SELECT id FROM roles WHERE name = 'admin'))
ON CONFLICT (email) DO NOTHING;

-- Врачи (учётные записи, будут привязаны к записям doctors)
INSERT INTO users (username, email, password, confirmed, role_id) VALUES
  ('Андреев Виктор Иванович',      'andreev@medlife.ru',       crypt('Doctor123', gen_salt('bf', 10)), true, (SELECT id FROM roles WHERE name = 'doctor')),
  ('Бобов Юрий Львович',           'bobov@medlife.ru',         crypt('Doctor123', gen_salt('bf', 10)), true, (SELECT id FROM roles WHERE name = 'doctor')),
  ('Бутерина Алина Александровна', 'buterina@medlife.ru',      crypt('Doctor123', gen_salt('bf', 10)), true, (SELECT id FROM roles WHERE name = 'doctor')),
  ('Бурков Анатолий Николаевич',   'burkov@medlife.ru',        crypt('Doctor123', gen_salt('bf', 10)), true, (SELECT id FROM roles WHERE name = 'doctor')),
  ('Гасдинова Наталья Николаевна', 'gasdinova@medlife.ru',     crypt('Doctor123', gen_salt('bf', 10)), true, (SELECT id FROM roles WHERE name = 'doctor')),
  ('Грищенко Елена Юрьевна',       'grishchenko@medlife.ru',   crypt('Doctor123', gen_salt('bf', 10)), true, (SELECT id FROM roles WHERE name = 'doctor')),
  ('Ломакова Ирина Юрьевна',       'lomakova@medlife.ru',      crypt('Doctor123', gen_salt('bf', 10)), true, (SELECT id FROM roles WHERE name = 'doctor')),
  ('Музыченко Анастасия Сергеевна','muzychenko@medlife.ru',    crypt('Doctor123', gen_salt('bf', 10)), true, (SELECT id FROM roles WHERE name = 'doctor'))
ON CONFLICT (email) DO NOTHING;

-- Пациенты
INSERT INTO users (username, email, password, confirmed, role_id) VALUES
  ('Иван Петров',    'patient1@test.ru', crypt('Patient123', gen_salt('bf', 10)), true, (SELECT id FROM roles WHERE name = 'user')),
  ('Мария Сидорова', 'patient2@test.ru', crypt('Patient123', gen_salt('bf', 10)), true, (SELECT id FROM roles WHERE name = 'user')),
  ('Алексей Козлов', 'patient3@test.ru', crypt('Patient123', gen_salt('bf', 10)), true, (SELECT id FROM roles WHERE name = 'user'))
ON CONFLICT (email) DO NOTHING;

-- -------------------------------------------------------------
-- Специализации
-- -------------------------------------------------------------
INSERT INTO specializations (name) VALUES
  ('Онкология (маммология)'),
  ('Кардиология'),
  ('Неврология'),
  ('Гастроэнтерология'),
  ('Хирургия'),
  ('Дерматология'),
  ('Терапия'),
  ('Урология'),
  ('Пульмонология'),
  ('Эндокринология')
ON CONFLICT DO NOTHING;

-- -------------------------------------------------------------
-- Врачи (привязка к user_id)
-- -------------------------------------------------------------
INSERT INTO doctors (fullname, description, user_id) VALUES
  ('Андреев Виктор Иванович',
   'Врач-онколог-маммолог высшей категории. Стаж более 20 лет. Специализируется на ранней диагностике и лечении заболеваний молочной железы.',
   (SELECT id FROM users WHERE email = 'andreev@medlife.ru')),

  ('Бобов Юрий Львович',
   'Хирург высшей категории. Выполняет весь спектр плановых и экстренных хирургических вмешательств. Стаж 18 лет.',
   (SELECT id FROM users WHERE email = 'bobov@medlife.ru')),

  ('Бутерина Алина Александровна',
   'Акушер-гинеколог. Ведёт беременность, занимается диагностикой и лечением гинекологических заболеваний. Стаж 12 лет.',
   (SELECT id FROM users WHERE email = 'buterina@medlife.ru')),

  ('Бурков Анатолий Николаевич',
   'Кардиолог, специалист по ЭКГ и УЗИ сердца. Кандидат медицинских наук. Стаж 15 лет.',
   (SELECT id FROM users WHERE email = 'burkov@medlife.ru')),

  ('Гасдинова Наталья Николаевна',
   'Терапевт высшей категории. Ведёт первичный приём, занимается профилактикой и лечением хронических заболеваний. Стаж 22 года.',
   (SELECT id FROM users WHERE email = 'gasdinova@medlife.ru')),

  ('Грищенко Елена Юрьевна',
   'Невролог. Специализируется на диагностике и лечении головных болей, остеохондроза, нарушений сна. Стаж 10 лет.',
   (SELECT id FROM users WHERE email = 'grishchenko@medlife.ru')),

  ('Ломакова Ирина Юрьевна',
   'Гастроэнтеролог. Диагностика и лечение заболеваний желудочно-кишечного тракта, печени. Стаж 14 лет.',
   (SELECT id FROM users WHERE email = 'lomakova@medlife.ru')),

  ('Музыченко Анастасия Сергеевна',
   'Эндокринолог. Занимается патологией щитовидной железы, сахарным диабетом, нарушениями обмена веществ. Стаж 9 лет.',
   (SELECT id FROM users WHERE email = 'muzychenko@medlife.ru'))
ON CONFLICT DO NOTHING;

-- -------------------------------------------------------------
-- Расписание (schedules.doctor_id — один врач, несколько дней)
-- day: 1=Пн, 2=Вт, 3=Ср, 4=Чт, 5=Пт, 6=Сб, 7=Вс
-- -------------------------------------------------------------

-- Андреев: Пн-Ср-Пт
INSERT INTO schedules (doctor_id, day, time_from, time_to) VALUES
  ((SELECT id FROM doctors WHERE fullname = 'Андреев Виктор Иванович'), 1, '09:00', '15:00'),
  ((SELECT id FROM doctors WHERE fullname = 'Андреев Виктор Иванович'), 3, '09:00', '15:00'),
  ((SELECT id FROM doctors WHERE fullname = 'Андреев Виктор Иванович'), 5, '09:00', '13:00');

-- Бобов: Вт-Чт
INSERT INTO schedules (doctor_id, day, time_from, time_to) VALUES
  ((SELECT id FROM doctors WHERE fullname = 'Бобов Юрий Львович'), 2, '08:00', '14:00'),
  ((SELECT id FROM doctors WHERE fullname = 'Бобов Юрий Львович'), 4, '08:00', '14:00');

-- Бутерина: Пн-Вт-Чт-Пт
INSERT INTO schedules (doctor_id, day, time_from, time_to) VALUES
  ((SELECT id FROM doctors WHERE fullname = 'Бутерина Алина Александровна'), 1, '14:00', '19:00'),
  ((SELECT id FROM doctors WHERE fullname = 'Бутерина Алина Александровна'), 2, '14:00', '19:00'),
  ((SELECT id FROM doctors WHERE fullname = 'Бутерина Алина Александровна'), 4, '14:00', '19:00'),
  ((SELECT id FROM doctors WHERE fullname = 'Бутерина Алина Александровна'), 5, '14:00', '18:00');

-- Бурков: Пн-Ср-Пт
INSERT INTO schedules (doctor_id, day, time_from, time_to) VALUES
  ((SELECT id FROM doctors WHERE fullname = 'Бурков Анатолий Николаевич'), 1, '10:00', '16:00'),
  ((SELECT id FROM doctors WHERE fullname = 'Бурков Анатолий Николаевич'), 3, '10:00', '16:00'),
  ((SELECT id FROM doctors WHERE fullname = 'Бурков Анатолий Николаевич'), 5, '10:00', '14:00');

-- Гасдинова: Пн-Вт-Ср-Чт-Пт
INSERT INTO schedules (doctor_id, day, time_from, time_to) VALUES
  ((SELECT id FROM doctors WHERE fullname = 'Гасдинова Наталья Николаевна'), 1, '08:00', '14:00'),
  ((SELECT id FROM doctors WHERE fullname = 'Гасдинова Наталья Николаевна'), 2, '08:00', '14:00'),
  ((SELECT id FROM doctors WHERE fullname = 'Гасдинова Наталья Николаевна'), 3, '08:00', '14:00'),
  ((SELECT id FROM doctors WHERE fullname = 'Гасдинова Наталья Николаевна'), 4, '08:00', '14:00'),
  ((SELECT id FROM doctors WHERE fullname = 'Гасдинова Наталья Николаевна'), 5, '08:00', '13:00');

-- Грищенко: Вт-Чт-Сб
INSERT INTO schedules (doctor_id, day, time_from, time_to) VALUES
  ((SELECT id FROM doctors WHERE fullname = 'Грищенко Елена Юрьевна'), 2, '10:00', '17:00'),
  ((SELECT id FROM doctors WHERE fullname = 'Грищенко Елена Юрьевна'), 4, '10:00', '17:00'),
  ((SELECT id FROM doctors WHERE fullname = 'Грищенко Елена Юрьевна'), 6, '10:00', '15:00');

-- Ломакова: Пн-Ср-Пт
INSERT INTO schedules (doctor_id, day, time_from, time_to) VALUES
  ((SELECT id FROM doctors WHERE fullname = 'Ломакова Ирина Юрьевна'), 1, '12:00', '18:00'),
  ((SELECT id FROM doctors WHERE fullname = 'Ломакова Ирина Юрьевна'), 3, '12:00', '18:00'),
  ((SELECT id FROM doctors WHERE fullname = 'Ломакова Ирина Юрьевна'), 5, '12:00', '17:00');

-- Музыченко: Вт-Чт
INSERT INTO schedules (doctor_id, day, time_from, time_to) VALUES
  ((SELECT id FROM doctors WHERE fullname = 'Музыченко Анастасия Сергеевна'), 2, '09:00', '15:00'),
  ((SELECT id FROM doctors WHERE fullname = 'Музыченко Анастасия Сергеевна'), 4, '09:00', '15:00');

-- -------------------------------------------------------------
-- Специализации врачей
-- -------------------------------------------------------------
INSERT INTO doctor_specializations (doctor_id, specialization_id) VALUES
  ((SELECT id FROM doctors WHERE fullname = 'Андреев Виктор Иванович'),
   (SELECT id FROM specializations WHERE name = 'Онкология (маммология)')),

  ((SELECT id FROM doctors WHERE fullname = 'Бобов Юрий Львович'),
   (SELECT id FROM specializations WHERE name = 'Хирургия')),

  ((SELECT id FROM doctors WHERE fullname = 'Бутерина Алина Александровна'),
   (SELECT id FROM specializations WHERE name = 'Эндокринология')),

  ((SELECT id FROM doctors WHERE fullname = 'Бурков Анатолий Николаевич'),
   (SELECT id FROM specializations WHERE name = 'Кардиология')),

  ((SELECT id FROM doctors WHERE fullname = 'Гасдинова Наталья Николаевна'),
   (SELECT id FROM specializations WHERE name = 'Терапия')),

  ((SELECT id FROM doctors WHERE fullname = 'Грищенко Елена Юрьевна'),
   (SELECT id FROM specializations WHERE name = 'Неврология')),

  ((SELECT id FROM doctors WHERE fullname = 'Ломакова Ирина Юрьевна'),
   (SELECT id FROM specializations WHERE name = 'Гастроэнтерология')),

  ((SELECT id FROM doctors WHERE fullname = 'Музыченко Анастасия Сергеевна'),
   (SELECT id FROM specializations WHERE name = 'Эндокринология'))
ON CONFLICT DO NOTHING;

-- -------------------------------------------------------------
-- Категории услуг
-- -------------------------------------------------------------
INSERT INTO service_categories (name, description, favorite) VALUES
  ('Консультации специалистов', 'Первичные и повторные консультации врачей всех специальностей', true),
  ('Диагностика',               'УЗИ, ЭКГ, рентген, лабораторные исследования и другие методы диагностики', true),
  ('Маммология',                'Диагностика и лечение заболеваний молочной железы, маммография 3D', true),
  ('Хирургия',                  'Плановые и экстренные хирургические операции и манипуляции', false)
ON CONFLICT DO NOTHING;

-- -------------------------------------------------------------
-- Услуги
-- -------------------------------------------------------------
INSERT INTO services (name, description, price, service_category_id, specialization_id) VALUES
  -- Консультации
  ('Консультация кардиолога',       'Первичный приём, оценка состояния сердечно-сосудистой системы, назначение лечения',         1500,
   (SELECT id FROM service_categories WHERE name = 'Консультации специалистов'),
   (SELECT id FROM specializations WHERE name = 'Кардиология')),

  ('Консультация невролога',        'Диагностика неврологических расстройств, головных болей, нарушений сна',                    1400,
   (SELECT id FROM service_categories WHERE name = 'Консультации специалистов'),
   (SELECT id FROM specializations WHERE name = 'Неврология')),

  ('Консультация терапевта',        'Первичный осмотр, расшифровка анализов, направления к специалистам',                        900,
   (SELECT id FROM service_categories WHERE name = 'Консультации специалистов'),
   (SELECT id FROM specializations WHERE name = 'Терапия')),

  ('Консультация гастроэнтеролога', 'Диагностика заболеваний ЖКТ, печени, поджелудочной железы, назначение лечения',            1600,
   (SELECT id FROM service_categories WHERE name = 'Консультации специалистов'),
   (SELECT id FROM specializations WHERE name = 'Гастроэнтерология')),

  ('Консультация эндокринолога',    'Диагностика нарушений эндокринной системы, коррекция гормонального фона',                  1500,
   (SELECT id FROM service_categories WHERE name = 'Консультации специалистов'),
   (SELECT id FROM specializations WHERE name = 'Эндокринология')),

  -- Диагностика
  ('ЭКГ',                           'Электрокардиограмма с расшифровкой врача',                                                   500,
   (SELECT id FROM service_categories WHERE name = 'Диагностика'),
   (SELECT id FROM specializations WHERE name = 'Кардиология')),

  ('УЗИ щитовидной железы',         'Ультразвуковое исследование щитовидной железы',                                             900,
   (SELECT id FROM service_categories WHERE name = 'Диагностика'),
   (SELECT id FROM specializations WHERE name = 'Эндокринология')),

  ('УЗИ органов брюшной полости',   'Ультразвуковое исследование печени, поджелудочной железы, желчного пузыря, селезёнки',     1200,
   (SELECT id FROM service_categories WHERE name = 'Диагностика'),
   (SELECT id FROM specializations WHERE name = 'Гастроэнтерология')),

  ('Рентген грудной клетки',        'Цифровой рентген лёгких и органов грудной клетки',                                          700,
   (SELECT id FROM service_categories WHERE name = 'Диагностика'),
   NULL),

  -- Маммология
  ('Маммография 3D',                'Цифровая маммография с технологией томосинтеза (GIOTTO 3D)',                                2200,
   (SELECT id FROM service_categories WHERE name = 'Маммология'),
   (SELECT id FROM specializations WHERE name = 'Онкология (маммология)')),

  ('Консультация онколога-маммолога','Осмотр, пальпация, разбор результатов маммографии и УЗИ молочных желёз',                  1800,
   (SELECT id FROM service_categories WHERE name = 'Маммология'),
   (SELECT id FROM specializations WHERE name = 'Онкология (маммология)')),

  ('УЗИ молочных желёз',            'Ультразвуковое исследование молочных желёз с регионарными лимфатическими узлами',          1000,
   (SELECT id FROM service_categories WHERE name = 'Маммология'),
   (SELECT id FROM specializations WHERE name = 'Онкология (маммология)')),

  -- Хирургия
  ('Консультация хирурга',          'Первичный осмотр, определение показаний к операции',                                        1200,
   (SELECT id FROM service_categories WHERE name = 'Хирургия'),
   (SELECT id FROM specializations WHERE name = 'Хирургия')),

  ('Удаление доброкачественных новообразований', 'Удаление папиллом, атером, липом под местной анестезией',                     3500,
   (SELECT id FROM service_categories WHERE name = 'Хирургия'),
   (SELECT id FROM specializations WHERE name = 'Хирургия'))

ON CONFLICT DO NOTHING;

-- -------------------------------------------------------------
-- Лицензии
-- -------------------------------------------------------------
INSERT INTO licenses (name, description) VALUES
  ('Лицензия №ЛО-42-01-006890',  'Лицензия на осуществление медицинской деятельности, выдана Департаментом здравоохранения Кемеровской области'),
  ('Свидетельство о регистрации', 'Свидетельство о государственной регистрации медицинского центра «МЕДЛАЙФ»'),
  ('Сертификат соответствия',     'Сертификат соответствия оборудования и услуг стандартам качества'),
  ('Лицензия на рентгенологию',   'Разрешение на использование источников ионизирующего излучения')
ON CONFLICT DO NOTHING;

-- -------------------------------------------------------------
-- Обратные звонки
-- -------------------------------------------------------------
INSERT INTO callback_requests (name, phone, message, status) VALUES
  ('Иван Петров',      '+7 (903) 123-45-67', 'Хочу записаться на консультацию к кардиологу',           'new'),
  ('Мария Сидорова',   '+7 (923) 987-65-43', 'Интересует маммография 3D, есть ли направление?',        'processed'),
  ('Алексей Козлов',   '+7 (3846) 111-22-33', NULL,                                                     'completed'),
  ('Елена Новикова',   '+7 (903) 444-55-66', 'Нужна информация о стоимости УЗИ',                       'new'),
  ('Татьяна Смирнова', '+7 (923) 321-00-11', 'Хочу записать маму на приём к терапевту, ей 72 года',    'new');

-- -------------------------------------------------------------
-- Записи на приём (тестовые)
-- Используем будущие даты относительно 2026-04-11
-- -------------------------------------------------------------
INSERT INTO appointments (patient_id, doctor_id, scheduled_at, status, notes) VALUES
  -- Иван Петров → Гасдинова (терапевт), понедельник 2026-04-13 09:00
  ((SELECT id FROM users WHERE email = 'patient1@test.ru'),
   (SELECT id FROM doctors WHERE fullname = 'Гасдинова Наталья Николаевна'),
   '2026-04-13 09:00:00', 'confirmed', 'Плановый осмотр, жалобы на давление'),

  -- Мария Сидорова → Андреев (маммолог), среда 2026-04-15 10:00
  ((SELECT id FROM users WHERE email = 'patient2@test.ru'),
   (SELECT id FROM doctors WHERE fullname = 'Андреев Виктор Иванович'),
   '2026-04-15 10:00:00', 'pending', 'Первичная консультация маммолога'),

  -- Алексей Козлов → Бурков (кардиолог), пятница 2026-04-17 11:00
  ((SELECT id FROM users WHERE email = 'patient3@test.ru'),
   (SELECT id FROM doctors WHERE fullname = 'Бурков Анатолий Николаевич'),
   '2026-04-18 11:00:00', 'confirmed', 'Боли в груди, одышка при нагрузке'),

  -- Иван Петров → Грищенко (невролог), вторник 2026-04-22 10:00
  ((SELECT id FROM users WHERE email = 'patient1@test.ru'),
   (SELECT id FROM doctors WHERE fullname = 'Грищенко Елена Юрьевна'),
   '2026-04-22 10:00:00', 'pending', 'Хронические головные боли'),

  -- Завершённая запись (прошедшая дата)
  ((SELECT id FROM users WHERE email = 'patient2@test.ru'),
   (SELECT id FROM doctors WHERE fullname = 'Гасдинова Наталья Николаевна'),
   '2026-04-01 09:00:00', 'completed', 'Профилактический осмотр');

-- Добавляем результат завершённой записи
UPDATE appointments
SET results = 'Состояние удовлетворительное. АД 120/80, ЧСС 72. Рекомендована диета, умеренные физические нагрузки. Повторный приём через 3 месяца.'
WHERE patient_id  = (SELECT id FROM users WHERE email = 'patient2@test.ru')
  AND doctor_id   = (SELECT id FROM doctors WHERE fullname = 'Гасдинова Наталья Николаевна')
  AND status      = 'completed';
