INSERT INTO daily_works (id, name, description, role_id, start_time, end_time, created_at, updated_at) VALUES
(1, 'Ambil telur 1, cek kualitas, hitung jumlah', 'Dahulukan kandang produksi', 1, '11:00:00', '12:00:00', NOW(), NOW()),
(2, 'Ishoma', 'Di zona kuning', 1, '12:00:00', '13:30:00', NOW(), NOW()),
(3, 'Ambil telur 2, cek kualitas, hitung jumlah', '-', 1, '13:30:00', '14:30:00', NOW(), NOW()),
(4, 'Selesaikan pencatatan jumlah telur harian', '-', 1, '14:30:00', '14:45:00', NOW(), NOW()),
(5, 'Antar telur ke tempat penampungan, cek telur pecah', '-', 1, '14:45:00', '15:30:00', NOW(), NOW()),
(6, 'Pengepakan telur', '-', 1, '15:00:00', '17:00:00', NOW(), NOW()),

(7, 'Bersihkan tempat pakan & minum - Beri makan 1', '-', 2, '07:00:00', '10:00:00', NOW(), NOW()),
(8, 'Pengadukan pakan untuk konsumsi esok hari', '-', 2, '10:00:00', '11:30:00', NOW(), NOW()),
('Pengedropan pakan ke tiap kandang', '-', 2, '11:30:00', '12:00:00', NOW(), NOW()),
(9, 'Ishoma', 'Di zona kuning', 2, '12:00:00', '13:30:00', NOW(), NOW()),
('Bersihkan tempat pakan & minum - Beri makan 2', '-', 2, '13:30:00', '14:30:00', NOW(), NOW()),
(10, 'Selesaikan pencatatan performa ayam harian', '-', 2, '14:30:00', '15:00:00', NOW(), NOW()),

(11, 'Jaga malam', '-', 7, '21:00:00', '03:00:00', NOW(), NOW()),

(12, 'Supervisi pekerja kandang, pastikan pakan dan air cukup', '-', 3, '07:00:00', '10:00:00', NOW(), NOW()),
(13, 'Evaluasi performa ayam dan kesehatan', '-', 3, '10:30:00', '12:00:00', NOW(), NOW()),
(14, 'Cek pengobatan ayam sakit, pastikan SOP berjalan', '-', 3, '13:30:00', '15:00:00', NOW(), NOW()),
(15, 'Rekap laporan harian pekerja kandang dan telur', '-', 3, '15:00:00', '15:30:00', NOW(), NOW());
