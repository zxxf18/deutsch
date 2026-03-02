-- Deutsch 初始化数据
-- 执行前需确保表已创建（通过应用 AutoMigrate 或手动建表）
-- 可重复执行：使用 INSERT IGNORE 或先检查再插入

-- ============================================================
-- 1. 联邦州（german_states）
-- ============================================================
INSERT INTO german_states (id, slug, name, name_cn, sort_order, created_at, updated_at) VALUES
(UUID(), 'baden-wuerttemberg', 'Baden-Württemberg', '巴登-符腾堡州', 1, NOW(), NOW()),
(UUID(), 'bayern', 'Bayern', '巴伐利亚州', 2, NOW(), NOW()),
(UUID(), 'berlin', 'Berlin', '柏林', 3, NOW(), NOW()),
(UUID(), 'brandenburg', 'Brandenburg', '勃兰登堡州', 4, NOW(), NOW()),
(UUID(), 'bremen', 'Bremen', '不来梅州', 5, NOW(), NOW()),
(UUID(), 'hamburg', 'Hamburg', '汉堡', 6, NOW(), NOW()),
(UUID(), 'hessen', 'Hessen', '黑森州', 7, NOW(), NOW()),
(UUID(), 'mecklenburg-vorpommern', 'Mecklenburg-Vorpommern', '梅克伦堡-前波莫瑞州', 8, NOW(), NOW()),
(UUID(), 'niedersachsen', 'Niedersachsen', '下萨克森州', 9, NOW(), NOW()),
(UUID(), 'nordrhein-westfalen', 'Nordrhein-Westfalen', '北莱茵-威斯特法伦州', 10, NOW(), NOW()),
(UUID(), 'rheinland-pfalz', 'Rheinland-Pfalz', '莱茵兰-普法尔茨州', 11, NOW(), NOW()),
(UUID(), 'saarland', 'Saarland', '萨尔州', 12, NOW(), NOW()),
(UUID(), 'sachsen', 'Sachsen', '萨克森州', 13, NOW(), NOW()),
(UUID(), 'sachsen-anhalt', 'Sachsen-Anhalt', '萨克森-安哈尔特州', 14, NOW(), NOW()),
(UUID(), 'schleswig-holstein', 'Schleswig-Holstein', '石勒苏益格-荷尔斯泰因州', 15, NOW(), NOW()),
(UUID(), 'thueringen', 'Thüringen', '图林根州', 16, NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at = NOW();

-- ============================================================
-- 2. 管理员用户（可选，首次部署时创建）
-- 密码为 123456 的 SHA256 哈希
-- ============================================================
-- INSERT INTO users (id, username, email, password_hash, role, is_enabled, created_at, updated_at)
-- VALUES (
--   UUID(),
--   'admin001',
--   'admin@example.com',
--   '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92',
--   'admin',
--   1,
--   NOW(),
--   NOW()
-- );

-- ============================================================
-- 3. 初始邀请码（可选，用于首次注册后改 role）
-- ============================================================
-- INSERT INTO invite_codes (id, code, expires_at, is_enabled, created_at)
-- VALUES (UUID(), 'ADMIN00001', DATE_ADD(NOW(), INTERVAL 30 DAY), 1, NOW());
