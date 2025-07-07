INSERT INTO daily_works (id, description, role_id, start_time, end_time, created_by, created_at, updated_by, updated_at) VALUES
-- Role 1
(1, 'Ambil telur 1, cek kualitas, hitung jumlah. Dahulukan kandang produksi', 1, '11:00', '12:00', NULL, NOW(), NULL, NOW()),
(2, 'Ishoma di zona kuning', 1, '12:00', '13:30', NULL, NOW(), NULL, NOW()),
(3, 'Ambil telur 2, cek kualitas, hitung jumlah', 1, '13:30', '14:30', NULL, NOW(), NULL, NOW()),
(4, 'Selesaikan pencatatan jumlah telur harian', 1, '14:30', '14:45', NULL, NOW(), NULL, NOW()),
(5, 'Antar telur ke tempat penampungan, cek telur pecah', 1, '14:45', '15:30', NULL, NOW(), NULL, NOW()),
(6, 'Pengepakan telur', 1, '15:00', '17:00', NULL, NOW(), NULL, NOW()),

-- Role 2
(7, 'Bersihkan tempat pakan & minum - Beri makan 1', 2, '07:00', '10:00', NULL, NOW(), NULL, NOW()),
(8, 'Pengadukan pakan untuk konsumsi esok hari', 2, '10:00', '11:30', NULL, NOW(), NULL, NOW()),
(9, 'Pengedropan pakan ke tiap kandang', 2, '11:30', '12:00', NULL, NOW(), NULL, NOW()),
(10, 'Ishoma di zona kuning', 2, '12:00', '13:30', NULL, NOW(), NULL, NOW()),
(11, 'Bersihkan tempat pakan & minum - Beri makan 2', 2, '13:30', '14:30', NULL, NOW(), NULL, NOW()),
(12, 'Selesaikan pencatatan performa ayam harian', 2, '14:30', '15:00', NULL, NOW(), NULL, NOW()),

-- Role 7
(13, 'Jaga malam', 7, '21:00', '03:00', NULL, NOW(), NULL, NOW()),

-- Role 3
(14, 'Supervisi pekerja kandang, pastikan pakan dan air cukup', 3, '07:00', '10:00', NULL, NOW(), NULL, NOW()),
(15, 'Evaluasi performa ayam dan kesehatan', 3, '10:30', '12:00', NULL, NOW(), NULL, NOW()),
(16, 'Cek pengobatan ayam sakit, pastikan SOP berjalan', 3, '13:30', '15:00', NULL, NOW(), NULL, NOW()),
(17, 'Rekap laporan harian pekerja kandang dan telur', 3, '15:00', '15:30', NULL, NOW(), NULL, NOW());
