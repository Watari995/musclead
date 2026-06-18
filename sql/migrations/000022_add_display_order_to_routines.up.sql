ALTER TABLE routines
    ADD COLUMN display_order INT NOT NULL DEFAULT 0 AFTER name,
    ADD KEY idx_user_id_display_order (user_id, display_order);

-- 既存行は created_at 昇順で 0 始まりの連番を採番する(ユーザーごと)
UPDATE routines r
JOIN (
    SELECT id, ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY created_at ASC) - 1 AS rn
    FROM routines
) t ON r.id = t.id
SET r.display_order = t.rn;
