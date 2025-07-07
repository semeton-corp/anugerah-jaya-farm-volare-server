INSERT INTO store_items (
    store_id,
    item_id,
    quantity,
    created_at,
    created_by,
    updated_at,
    updated_by
) VALUES
-- Toko A (store_id = 10)
(10, 10, 100, NOW(), NULL, NOW(), NULL), -- Telur OK
(10, 11,  50, NOW(), NULL, NOW(), NULL), -- Telur Retak
(10, 12,  30, NOW(), NULL, NOW(), NULL), -- Telur Bonyok

-- Toko B (store_id = 11)
(11, 10, 120, NOW(), NULL, NOW(), NULL),
(11, 11,  60, NOW(), NULL, NOW(), NULL),
(11, 12,  40, NOW(), NULL, NOW(), NULL),

-- Toko C (store_id = 12)
(12, 10,  90, NOW(), NULL, NOW(), NULL),
(12, 11,  45, NOW(), NULL, NOW(), NULL),
(12, 12,  25, NOW(), NULL, NOW(), NULL),

-- Toko D (store_id = 13)
(13, 10, 130, NOW(), NULL, NOW(), NULL),
(13, 11,  55, NOW(), NULL, NOW(), NULL),
(13, 12,  35, NOW(), NULL, NOW(), NULL),

-- Toko E (store_id = 14)
(14, 10, 110, NOW(), NULL, NOW(), NULL),
(14, 11,  65, NOW(), NULL, NOW(), NULL),
(14, 12,  28, NOW(), NULL, NOW(), NULL);
