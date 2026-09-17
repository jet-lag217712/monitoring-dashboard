ALTER TABLE interface_samples
  DROP COLUMN IF EXISTS out_packets,
  DROP COLUMN IF EXISTS in_packets;

ALTER TABLE interfaces
  DROP COLUMN IF EXISTS duplex;
