-- Hapus semua data dummy migration 000004 (urutan terbalik FK)
DELETE FROM denda        WHERE id LIKE 'f0000000-0000-0000-0000-%';
DELETE FROM pembayaran   WHERE id LIKE 'e0000000-0000-0000-0000-%';
DELETE FROM kewajiban_pajak WHERE id LIKE 'd0000000-0000-0000-0000-%';
DELETE FROM kendaraan    WHERE id LIKE 'c0000000-0000-0000-0000-%';
DELETE FROM users        WHERE id LIKE 'b0000000-0000-0000-0000-%';
DELETE FROM wajib_pajak  WHERE id LIKE 'a0000000-0000-0000-0000-%';
