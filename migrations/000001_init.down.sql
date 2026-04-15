-- tables
DROP TABLE IF EXISTS message_reads;
DROP TABLE IF EXISTS message_files;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS pinned_rooms;
DROP TABLE IF EXISTS room_members;
DROP TABLE IF EXISTS rooms;
DROP TABLE IF EXISTS call_participants;
DROP TABLE IF EXISTS calls;
DROP TABLE IF EXISTS post_favorites;
DROP TABLE IF EXISTS post_media;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS user_blocks;
DROP TABLE IF EXISTS user_follows;
DROP TABLE IF EXISTS auth_providers;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS media;

-- types
DROP TYPE IF EXISTS file_type;
DROP TYPE IF EXISTS room_type;
DROP TYPE IF EXISTS call_participant_role;
DROP TYPE IF EXISTS call_joinable_by;
DROP TYPE IF EXISTS call_type;
DROP TYPE IF EXISTS auth_provider_type;
DROP TYPE IF EXISTS media_status;
DROP TYPE IF EXISTS media_type;
DROP TYPE IF EXISTS media_usage;

-- domains
DROP DOMAIN IF EXISTS ulid;
