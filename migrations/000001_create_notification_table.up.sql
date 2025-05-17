CREATE TYPE notification_channel AS ENUM ('email', 'sms', 'push', 'webhook');

CREATE TABLE IF NOT EXISTS notifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  channel notification_channel NOT NULL,
  recipient TEXT NOT NULL,
  message TEXT NOT NULL,
  data TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_notifications_recipient ON notifications (recipient);
CREATE INDEX IF NOT EXISTS idx_notifications_channel ON notifications (channel);
