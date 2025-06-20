CREATE DATABASE IF NOT EXISTS "mawjood";

SET DATABASE  = "mawjood";

-- Create the contents table to store core content metadata
CREATE TABLE IF NOT EXISTS contents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    language VARCHAR(50),
    duration_seconds INT,
    published_at TIMESTAMPTZ,
    content_type VARCHAR(20) NOT NULL, -- 'podcast' or 'documentary'
    url VARCHAR(255),
    platform_name VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL -- For soft delete functionality
);

-- Create the tags table to store unique tags
CREATE TABLE IF NOT EXISTS tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE
);

-- Create the content_tags join table to associate content with tags
CREATE TABLE IF NOT EXISTS content_tags (
    content_id UUID NOT NULL REFERENCES contents(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (content_id, tag_id)
);

-- Create indexes for efficient querying
-- Index for searching content by title
CREATE INDEX IF NOT EXISTS idx_contents_title ON contents (title);

-- Full-text search indexes for title and description
CREATE INVERTED INDEX IF NOT EXISTS idx_contents_title_search ON contents (title gin_trgm_ops);
CREATE INVERTED INDEX IF NOT EXISTS idx_contents_description_search ON contents (description gin_trgm_ops);

-- Index for finding tags by name
CREATE INDEX IF NOT EXISTS idx_tags_name ON tags (name);

-- Indexes on the join table for efficient lookups in both directions
CREATE INDEX IF NOT EXISTS idx_content_tags_tag_id ON content_tags (tag_id);
CREATE INDEX IF NOT EXISTS idx_content_tags_content_id ON content_tags (content_id);

-- Index for content type filtering
CREATE INDEX IF NOT EXISTS idx_contents_content_type ON contents (content_type);

-- Index for language filtering
CREATE INDEX IF NOT EXISTS idx_contents_language ON contents (language);

-- Index for published date sorting
CREATE INDEX IF NOT EXISTS idx_contents_published_at ON contents (published_at DESC);

-- Index for created date sorting
CREATE INDEX IF NOT EXISTS idx_contents_created_at ON contents (created_at DESC);

-- Index for soft delete functionality
CREATE INDEX IF NOT EXISTS idx_contents_deleted_at ON contents (deleted_at);

-- Partial index for active (non-deleted) content for better performance
CREATE INDEX IF NOT EXISTS idx_contents_active ON contents (created_at DESC) WHERE deleted_at IS NULL;

-- Insert seed data for tags
INSERT INTO tags (name) VALUES 
    ('technology'),
    ('science'),
    ('history'),
    ('business'),
    ('health'),
    ('education'),
    ('entertainment'),
    ('news'),
    ('comedy'),
    ('true-crime'),
    ('nature'),
    ('space'),
    ('psychology'),
    ('philosophy'),
    ('art'),
    ('music'),
    ('sports'),
    ('politics'),
    ('environment'),
    ('culture')
ON CONFLICT (name) DO NOTHING;

-- Insert seed data for contents (podcasts and documentaries)
INSERT INTO contents (title, description, language, duration_seconds, published_at, content_type, url, platform_name) VALUES 
    -- Podcasts
    ('The Daily Tech Brief', 'Your daily dose of technology news and insights from around the world. Covering AI, startups, and the latest innovations.', 'English', 1800, '2024-01-15 08:00:00+00', 'podcast', 'https://www.youtube.com/watch?v=dQw4w9WgXcQ', 'YouTube'),
    ('Science Friday', 'Weekly discussions about science, technology, and other cool stuff. Hosted by Ira Flatow.', 'English', 3600, '2024-01-12 14:00:00+00', 'podcast', 'https://www.youtube.com/watch?v=dQw4w9WgXcQ', 'YouTube'),
    ('Hidden Brain', 'Shankar Vedantam uses science and storytelling to reveal the unconscious patterns that drive human behavior.', 'English', 2700, '2024-01-10 09:00:00+00', 'podcast', 'https://www.youtube.com/watch?v=dQw4w9WgXcQ', 'YouTube'),

    -- Documentaries
    ('The Blue Planet', 'A comprehensive exploration of the world''s oceans, revealing the extraordinary creatures that inhabit them.', 'English', 3600, '2024-01-20 20:00:00+00', 'documentary', 'https://www.youtube.com/watch?v=dQw4w9WgXcQ', 'YouTube'),
    ('Cosmos: A Spacetime Odyssey', 'Hosted by Neil deGrasse Tyson, this series explores how we discovered the laws of nature and found our coordinates in space and time.', 'English', 5400, '2024-01-18 19:00:00+00', 'documentary', 'https://www.youtube.com/watch?v=dQw4w9WgXcQ', 'YouTube'),
    ('The Civil War', 'Ken Burns'' epic documentary about the American Civil War, featuring archival photographs and first-person accounts.', 'English', 7200, '2024-01-16 21:00:00+00', 'documentary', 'https://www.youtube.com/watch?v=dQw4w9WgXcQ', 'YouTube'),
    ('Planet Earth II', 'Experience the world from the viewpoint of animals themselves, using cutting-edge technology.', 'English', 3600, '2024-01-14 18:00:00+00', 'documentary', 'https://www.youtube.com/watch?v=dQw4w9WgXcQ', 'YouTube'),
    ('The Last Dance', 'A 10-part documentary series about Michael Jordan and the Chicago Bulls dynasty of the 1990s.', 'English', 4800, '2024-01-12 22:00:00+00', 'documentary', 'https://www.youtube.com/watch?v=dQw4w9WgXcQ', 'YouTube'),
    
    -- More diverse content
    ('Arabic Tech Talk', 'تحدث عن التكنولوجيا باللغة العربية - مناقشات حول الذكاء الاصطناعي والابتكار', 'Arabic', 2400, '2024-01-22 10:00:00+00', 'podcast', 'https://www.youtube.com/watch?v=dQw4w9WgXcQ', 'YouTube'),
    ('Historia de México', 'Un podcast que explora la rica historia de México, desde los aztecas hasta la época moderna.', 'Spanish', 3000, '2024-01-20 16:00:00+00', 'podcast', 'https://www.youtube.com/watch?v=dQw4w9WgXcQ', 'YouTube'),
    ('Le Monde en Français', 'L''actualité internationale analysée et expliquée en français.', 'French', 2100, '2024-01-19 11:00:00+00', 'podcast', 'https://www.youtube.com/watch?v=dQw4w9WgXcQ', 'YouTube'),
    ('Die Deutsche Welle', 'Nachrichten und Hintergrundberichte aus Deutschland und der Welt.', 'German', 1800, '2024-01-17 13:00:00+00', 'podcast', 'https://www.youtube.com/watch?v=dQw4w9WgXcQ', 'YouTube')
ON CONFLICT DO NOTHING;

-- Associate content with tags
-- First, let's get the content and tag IDs to create the relationships
WITH content_ids AS (
    SELECT id, title FROM contents
),
tag_ids AS (
    SELECT id, name FROM tags
)
INSERT INTO content_tags (content_id, tag_id)
SELECT c.id, t.id
FROM content_ids c, tag_ids t
WHERE 
    -- The Daily Tech Brief
    (c.title = 'The Daily Tech Brief' AND t.name IN ('technology', 'news', 'business')) OR
    -- Science Friday
    (c.title = 'Science Friday' AND t.name IN ('science', 'education', 'technology')) OR
    -- Hidden Brain
    (c.title = 'Hidden Brain' AND t.name IN ('psychology', 'science', 'education')) OR
    -- Planet Money
    (c.title = 'Planet Money' AND t.name IN ('business', 'education', 'news')) OR
    -- Serial
    (c.title = 'Serial' AND t.name IN ('true-crime', 'news', 'entertainment')) OR
    -- The Blue Planet
    (c.title = 'The Blue Planet' AND t.name IN ('nature', 'science', 'environment')) OR
    -- Cosmos
    (c.title = 'Cosmos: A Spacetime Odyssey' AND t.name IN ('space', 'science', 'education')) OR
    -- The Civil War
    (c.title = 'The Civil War' AND t.name IN ('history', 'politics', 'culture')) OR
    -- Planet Earth II
    (c.title = 'Planet Earth II' AND t.name IN ('nature', 'environment', 'science')) OR
    -- The Last Dance
    (c.title = 'The Last Dance' AND t.name IN ('sports', 'entertainment', 'culture')) OR
    -- Arabic Tech Talk
    (c.title = 'Arabic Tech Talk' AND t.name IN ('technology', 'culture')) OR
    -- Historia de México
    (c.title = 'Historia de México' AND t.name IN ('history', 'culture')) OR
    -- Le Monde en Français
    (c.title = 'Le Monde en Français' AND t.name IN ('news', 'politics', 'culture')) OR
    -- Die Deutsche Welle
    (c.title = 'Die Deutsche Welle' AND t.name IN ('news', 'politics')) OR
    -- 日本の歴史
    (c.title = '日本の歴史' AND t.name IN ('history', 'culture'))
ON CONFLICT DO NOTHING;