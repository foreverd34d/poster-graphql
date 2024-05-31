CREATE TABLE IF NOT EXISTS posts (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	title varchar(50) NOT NULL,
	author varchar(30) NOT NULL,
	content text NOT NULL,
	commentable boolean NOT NULL
);

CREATE TABLE IF NOT EXISTS comments (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	author varchar(30) NOT NULL,
	content varchar(100) NOT NULL,
	parent_comment_id uuid,
	post_id uuid REFERENCES posts
);
