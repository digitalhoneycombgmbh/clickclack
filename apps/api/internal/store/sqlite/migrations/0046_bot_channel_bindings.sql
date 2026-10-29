-- DHC: Bots fest an Channels binden. Ein Bot mit Bindung erreicht nur diese Channels (plus seine DMs).
CREATE TABLE IF NOT EXISTS bot_channel_bindings (
  bot_user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  channel_id TEXT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
  PRIMARY KEY (bot_user_id, channel_id)
);
