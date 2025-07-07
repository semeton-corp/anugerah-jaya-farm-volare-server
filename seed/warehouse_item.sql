INSERT INTO warehouse_items (
    item_id,
    warehouse_id,
    quantity,
    estimation_run_out,
    created_at,
    created_by,
    updated_at,
    updated_by
) VALUES
-- Gudang Sidodadi A (id=10)
(10, 10, 500, '2025-08-01', NOW(), NULL, NOW(), NULL),
(11, 10, 150, '2025-07-25', NOW(), NULL, NOW(), NULL),
(12, 10,  80, '2025-07-20', NOW(), NULL, NOW(), NULL),
(13, 10, 600, '2025-09-05', NOW(), NULL, NOW(), NULL),  -- Jagung
(14, 10, 200, '2025-08-18', NOW(), NULL, NOW(), NULL),  -- Dedak

-- Gudang Sukamaju B (id=11)
(10, 11, 300, '2025-08-10', NOW(), NULL, NOW(), NULL),
(11, 11, 100, '2025-07-28', NOW(), NULL, NOW(), NULL),
(12, 11,  60, '2025-07-18', NOW(), NULL, NOW(), NULL),
(14, 11, 350, '2025-08-22', NOW(), NULL, NOW(), NULL),  -- Dedak
(15, 11,  25, '2025-11-15', NOW(), NULL, NOW(), NULL),  -- Karpet

-- Gudang Sukamaju C (id=12)
(10, 12, 400, '2025-08-15', NOW(), NULL, NOW(), NULL),
(11, 12, 130, '2025-07-30', NOW(), NULL, NOW(), NULL),
(12, 12,  90, '2025-07-22', NOW(), NULL, NOW(), NULL),
(13, 12, 450, '2025-09-12', NOW(), NULL, NOW(), NULL),  -- Jagung

-- Gudang Sidodadi D (id=13)
(10, 13, 450, '2025-08-12', NOW(), NULL, NOW(), NULL),
(11, 13, 120, '2025-07-27', NOW(), NULL, NOW(), NULL),
(12, 13,  70, '2025-07-19', NOW(), NULL, NOW(), NULL),
(14, 13, 280, '2025-08-20', NOW(), NULL, NOW(), NULL),  -- Dedak
(15, 13,  40, '2025-12-01', NOW(), NULL, NOW(), NULL),  -- Karpet

-- Gudang Sukamaju E (id=14)
(10, 14, 350, '2025-08-08', NOW(), NULL, NOW(), NULL),
(11, 14, 140, '2025-07-29', NOW(), NULL, NOW(), NULL),
(12, 14,  85, '2025-07-21', NOW(), NULL, NOW(), NULL),
(13, 14, 500, '2025-09-08', NOW(), NULL, NOW(), NULL),  -- Jagung
(15, 14,  30, '2025-11-20', NOW(), NULL, NOW(), NULL);  -- Karpet
