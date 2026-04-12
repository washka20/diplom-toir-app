-- 000001_init.down.sql
-- Удаление всех таблиц в обратном порядке зависимостей

DROP TABLE IF EXISTS work_order_parts;
DROP TABLE IF EXISTS maintenance_logs;
DROP TABLE IF EXISTS work_orders;
DROP TABLE IF EXISTS maintenance_schedules;
DROP TABLE IF EXISTS repair_requests;
DROP TABLE IF EXISTS parts;
DROP TABLE IF EXISTS equipments;
DROP TABLE IF EXISTS users;
