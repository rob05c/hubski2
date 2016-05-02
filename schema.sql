----------------------------------------
-- pub structure
----------------------------------------

create table if not exists "publication_types" (
  id               integer,
  publication_type text
);

insert into "publication_types" (id, publication_type) select 0, 'story' where not exists (select 1 from "publication_types" where id = 0);

insert into "publication_types" (id, publication_type) select 1, 'comment' where not exists (select 1 from "publication_types" where id = 1);

create table if not exists "publications" (
  id          integer,
  type_id     integer,
  username    text, -- username fk, 'by' in structure
  time        integer,
  date        text, -- \todo remove? calculate from time?
  url         text,
  title       text,
  mail        boolean,
  tag         text, -- \todo make fk?
  tag2        text, -- \todo make fk?
  text        text,
  md          text, -- \todo remove, should be calculated
  web_domain  text,
  score       integer,
  deleted     boolean,
  draft       boolean,
  parent_id   integer, -- fk into publication
  locked      boolean,
  no_kill     boolean
);

-- maps to cc field in pub structure
create table if not exists "publication_cc" (
  id       integer, -- fk into publication, NOT pk, one-to-many
  username text
);

-- maps to ctag field in pub structure
create table if not exists "publication_community_tags" (
  id   integer, -- fk into publication, NOT pk, one-to-many
  tag  text
);

-- \todo figure out why there are both 'ctag' and 'ctags'
-- maps to ctags field in pub structure
create table if not exists "publication_community_tagses" (
  id       integer, -- fk into publication, NOT pk, one-to-many
  username text,
  tag      text
);

-- \todo remove this, when a fulltext search solution exists
create table if not exists "publication_search_text" (
  id   integer, -- fk into publication, NOT a pk, one pub has many ctagses
  word text
);

-- \todo remove this, when a fulltext search solution exists
create table if not exists "publication_search_title" (
  id   integer, -- fk into publication, NOT a pk, one pub has many ctagses
  word text
);

-- \todo remove this, when a fulltext search solution exists
create table if not exists "publication_search_url" (
  id   integer, -- fk into publication, NOT a pk, one to many
  word text
);

create table if not exists "publication_votes" (
  id       integer, -- fk into publication, NOT a pk, one to many
  vote_id  integer, -- fk int votes (which doesn't have a table yet)
  username text, -- fk into usernames. Probably duplicate data
  up       boolean, -- redundant? Do downvotes exist in hubski?
  num      integer -- \todo figure out what this is used for. Weight?
);

create table if not exists "publication_saved_by" (
  id       integer, -- fk into publication, NOT a pk, one to many
  username text 
);

create table if not exists "publication_shared_by" (
  id       integer, -- fk into publication, NOT a pk, one to many
  username text 
);

create table if not exists "publication_badged_by" (
  id       integer, -- fk into publication, NOT a pk, one to many
  username text 
);

-- \todo remove, replace with query
create table if not exists "publication_badged_kids" (
  id     integer, -- fk into publication, NOT a pk, one to many
  kid_id integer -- fk into publication
);

create table if not exists "publication_cubbed_by" (
  id       integer, -- fk into publication, NOT a pk, one to many
  username text
);

create table if not exists "publication_kids" (
  id     integer, -- fk into publication, NOT a pk, one to many
  kid_id integer  -- fk into publication
);

create table if not exists "donations" (
  username text, -- (will be) fk into users, NOT a pk, one to many
  donation_cents integer,
  donation_time timestamp
);

create table if not exists "passwords" (
  username text primary key, -- (will be) fk into users, NOT a pk, one to many
  encoding text,
  hash     text,
  salt     text
);

-- \todo make (token, username) composite key.
create table if not exists "tokens" (
  token text primary key,
  username text,-- (will be) fk into users, NOT a pk, one to many
	timestamp timestamp not null default current_timestamp
);


----------------------------------------
-- profile structure
----------------------------------------

-- \todo determine 'email' type
-- \todo put uncommon fields (like spammer) in their own table ?
create table if not exists "users" (
  id                      text, -- pk
  joined                  timestamp,
  inactive                boolean,
  clout                   double precision,
  word_count              boolean,
  average_com_numerator   integer,
  average_com_denominator integer,
  ignore_newbies          boolean,
  global_ignored          boolean,
  new_tabs                boolean,
  publication_tabs        boolean,
  reply_alerts            boolean,
  pm_alerts               boolean,
  follower_alerts         boolean,
  shout_outs              boolean,
  badge_alerts            boolean,
  saved_notifications     boolean,
  feed_times              boolean,
  share_counts            boolean,
  show_global_filtered    boolean,
  follows_badges          boolean,
  embed_videos            boolean,
  bio                     text,
  email                   boolean,
  hubski_style            text,
  homepage                text,
  follower_count          integer, -- \todo remove, duplicate of sum(users_followers)
  posts_count             integer, -- \todo remove, duplicate of sum(publications)
  shareds_count           integer,  -- \todo remove, duplicate of sum(publications_sharedby)
  unread_notifications    boolean,  -- \todo remove, duplicate of notifications
  last_com_time           timestamp,
  com_clout_date          timestamp,
  zen                     boolean,
  spammer                 boolean
);

create table if not exists "users_submitted" (
       id text, -- fk into users, NOT a pk, one to many
       publication_id integer -- fk into publications
);

create table if not exists "users_saved" (
       id text, -- fk into users, NOT a pk, one to many
       publication_id integer -- fk into publications
);

-- \todo determine if used, remove if not
create table if not exists "users_sticky" (
       id text, -- fk into users, NOT a pk, one to many
       publication_id integer -- fk into publications
);

-- \todo determine if used, remove if not
create table if not exists "users_hidden" (
       id text, -- fk into users, NOT a pk, one to many
       publication_id integer -- fk into publications
);

create table if not exists "users_mail" (
       id text, -- fk into users, NOT a pk, one to many
       publication_id integer -- fk into publications
);

create table if not exists "users_drafts" (
       id text, -- fk into users, NOT a pk, one to many
       publication_id integer -- fk into publications
);

create table if not exists "users_shareds" (
       id text, -- fk into users, NOT a pk, one to many
       publication_id integer -- fk into publications
);

create table if not exists "users_cubbeds" (
       id text, -- fk into users, NOT a pk, one to many
       publication_id integer -- fk into publications
);

-- \todo remove, duplicate of publication_votes
create table if not exists "users_votes" (
  id             text, -- fk into users, NOT a pk, one to many
  publication_id integer,
  vote_id        integer, -- fk int votes (which doesn't have a table yet)
  username       text, -- fk into usernames. Probably duplicate data
  web_domain     text,
  up             boolean -- redundant? Do downvotes exist in hubski?
);

-- \todo remove, duplicate of publication_community_tagses ?
create table if not exists "users_suggested_tags" (
       id text, -- fk into users, NOT a pk, one to many
       publication_id integer, -- fk into publications
       tag text
);

-- \todo remove, duplicate of publication_badged_by
create table if not exists "users_badged" (
       id text, -- fk into users, NOT a pk, one to many
       publication_id integer -- fk into publications
);

-- \todo remove, duplicate of publication_badged_by
create table if not exists "users_badges" (
       id text, -- fk into users, NOT a pk, one to many
       publication_id integer -- fk into publications
);

create table if not exists "users_ignoring" (
       id text, -- fk into users, NOT a pk, one to many
       ignoring_id text -- fk into users
);

create table if not exists "users_muting" (
       id text, -- fk into users, NOT a pk, one to many
       muting_id text -- fk into users
);

create table if not exists "users_hushing" (
       id text, -- fk into users, NOT a pk, one to many
       hushing_id text -- fk into users
);

create table if not exists "users_blocking" (
       id text, -- fk into users, NOT a pk, one to many
       blocking_id text -- fk into users
);

create table if not exists "users_ignoring_tag" (
       id text, -- fk into users, NOT a pk, one to many
       tag text
);

-- \todo determine if used
-- \todo determine what a 'dom' is, and its type
create table if not exists "users_ignoring_dom" (
       id text, -- fk into users, NOT a pk, one to many
       dom text
);

-- \todo remove, duplicate of users_ignoring
create table if not exists "users_ignored_by" (
       id text, -- fk into users, NOT a pk, one to many
       by_id text -- fk into users
);

-- \todo remove, duplicate of users_muting
create table if not exists "users_muted_by" (
       id text, -- fk into users, NOT a pk, one to many
       by_id text -- fk into users
);

-- \todo remove, duplicate of users_hushing
create table if not exists "users_hushed_by" (
       id text, -- fk into users, NOT a pk, one to many
       by_id text -- fk into users
);

-- \todo remove, duplicate of users_blocking
create table if not exists "users_blocked_by" (
       id text, -- fk into users, NOT a pk, one to many
       by_id text -- fk into users
);

create table if not exists "users_followed" (
       id text, -- fk into users, NOT a pk, one to many
       followed_id text -- fk into users
);

-- \todo remove, duplicate of users_followed
create table if not exists "users_follower" (
       id text, -- fk into users, NOT a pk, one to many
       follower_id text -- fk into users
);

-- \todo remove, duplicate of data in publications
create table if not exists "users_personal_tags" (
       id text, -- fk into users, NOT a pk, one to many
       tag text -- fk into users
);

create table if not exists "users_followed_tags" (
       id text, -- fk into users, NOT a pk, one to many
       tag text -- fk into users
);

-- \todo determine what 'dom' is, and its type
create table if not exists "users_followed_dom" (
       id text, -- fk into users, NOT a pk, one to many
       dom text -- fk into users
);

-- \todo determine if duplicate, if used
create table if not exists "users_notified" (
  id text, -- fk into users, NOT a pk, one to many
  notified_id integer -- fk into publications?
);

create table if not exists "users_password_hashes" (
  id text, -- pk, fk into users
  hash_type text,
  hash1 text,
  hash2 text
);

create table if not exists "users_cookies" (
  id text, -- pk, fk into users
  cookie text
);

create table if not exists "users_emails" (
  id text, -- pk, fk into users
  email text
);

CREATE VIEW "user_page_user_data" AS SELECT u.id, (SELECT count(1) AS followers FROM users_followed WHERE followed_id = u.id), (SELECT count(1) AS badges FROM users_badges where id = u.id), bio, (SELECT count(1) AS following FROM users_followed WHERE id = u.id), (SELECT count(1) AS followed_tags FROM users_followed_tags where id = u.id), (SELECT count(1) AS followed_domains FROM users_followed_dom where id = u.id), (SELECT count(1) AS badges_given FROM users_badged where id = u.id), CAST(clout AS INT) / 200 AS badges_total, joined, hubski_style AS style, (SELECT sum(donation_cents) AS donated_cents FROM donations where username = u.id), (SELECT COALESCE(sum(donation_cents), 0) AS donated_cents_this_year FROM donations where username = u.id AND donation_time > '2016-01-01') FROM users AS u;

CREATE VIEW "publication_posts" AS SELECT id,
(WITH RECURSIVE parents( id, parent_id )
AS (
  -- get leaf children
  SELECT id, parent_id
  FROM publications
  WHERE id = pp.id
  UNION ALL
  -- get all parents
  SELECT t.id, t.parent_id
  FROM parents p
  JOIN publications t
  ON p.parent_id = t.id
)
SELECT id from parents where parent_id is null) as post_id from publications as pp;

CREATE VIEW "publication_parents" AS SELECT id,
(WITH RECURSIVE parents( id, parent_id )
AS (
  -- get leaf children
  SELECT id, parent_id
  FROM publications
  WHERE id = pp.id
  UNION ALL
  -- get all parents
  SELECT t.id, t.parent_id
  FROM parents p
  JOIN publications t
  ON p.parent_id = t.id
)
SELECT id from parents) as parent_id from publications as pp;

CREATE VIEW "user_page_comments" AS SELECT
pubs.id,
pubs.username,
TIMESTAMP WITH TIME ZONE 'epoch' + pubs.time * INTERVAL '1 second' as time,
pubs.text,
pubs.score,
pubs.deleted,
pubs.draft,
pubs.parent_id,
posts.post_id as post_id,
pubs_post.title as post_title
FROM publications AS pubs
INNER JOIN publication_posts as posts ON posts.id = pubs.id
INNER JOIN publications AS pubs_post ON posts.post_id = pubs_post.id
WHERE pubs.type_id = (SELECT t.id FROM publication_types as t WHERE t.publication_type = 'comment') AND pubs.mail <> true;

create view "publication_community_tag" as select id, tag from (select counts.id, counts.tag, max(counts.count) as count from (select p.id, c.tag, count(c.tag) as count from publications as p, publication_community_tagses as c where p.id = c.id group by p.id, c.tag, p.username) AS counts group by counts.tag, counts.id) as maxes;

create index publication_badged_by_id_idx on publication_badged_by (id);
create index publication_badged_kids_id_idx on publication_badged_kids (id);
create index publication_cc_id_idx on publication_cc (id);
create index publication_community_tags_id_idx on publication_community_tags (id);
create index publication_community_tagses_id_idx on publication_community_tagses (id);
create index publication_cubbed_by_id_idx on publication_cubbed_by (id);
create index publication_kids_id_idx on publication_kids (id);
create index publication_saved_by_id_idx on publication_saved_by (id);
create index publication_search_text_id_idx on publication_search_text (id);
create index publication_search_title_id_idx on publication_search_title (id);
create index publication_search_url_id_idx on publication_search_url (id);
create index publication_shared_by_id_idx on publication_shared_by (id);
create index publication_votes_id_idx on publication_votes (id);
create index publications_id_idx on publications (id);
create index publications_parent_id_idx on publications (parent_id);
create index publications_mail_idx on publications (mail);
create index publications_deleted_idx on publications (deleted);
create index publications_draft_idx on publications (draft);
create index publication_community_tagses_username_idx on publication_community_tagses (username);
create index publication_community_tagses_tag_idx on publication_community_tagses (tag);
create index donations_username_idx on donations (username);
create index donations_donation_time_idx on donations (donation_time);

create index users_submitted_publication_id_idx on users_submitted (publication_id);
create index users_saved_publication_id_idx on users_saved (publication_id);
create index users_sticky_publication_id_idx on users_sticky (publication_id);
create index users_hidden_publication_id_idx on users_hidden (publication_id);
create index users_mail_publication_id_idx on users_mail (publication_id);
create index users_drafts_publication_id_idx on users_drafts (publication_id);
create index users_shareds_publication_id_idx on users_shareds (publication_id);
create index users_cubbeds_publication_id_idx on users_cubbeds (publication_id);
create index users_votes_id_idx on users_votes (id);
create index users_suggested_tags_publication_id_idx on users_suggested_tags (publication_id);
create index users_badged_publication_id_idx on users_badged (publication_id);
create index users_badges_publication_id_idx on users_badges (publication_id);
create index users_ignoring_ignoring_id_idx on users_ignoring (ignoring_id);
create index users_muting_muting_id_idx on users_muting (muting_id);
create index users_hushing_hushing_id_idx on users_hushing (hushing_id);
create index users_blocking_blocking_id_idx on users_blocking (blocking_id);
create index users_ignoring_tag_tag_idx on users_ignoring_tag (tag);
create index users_ignoring_dom_dom_idx on users_ignoring_dom (dom);
create index users_ignored_by_by_id_idx on users_ignored_by (by_id);
create index users_muted_by_by_id_idx on users_muted_by (by_id);
create index users_hushed_by_by_id_idx on users_hushed_by (by_id);
create index users_blocked_by_by_id_idx on users_blocked_by (by_id);
create index users_followed_followed_id_idx on users_followed (followed_id);
create index users_follower_follower_id_idx on users_follower (follower_id);
create index users_personal_tags_tag_idx on users_personal_tags (tag);
create index users_followed_tags_tag_idx on users_followed_tags (tag);
create index users_followed_dom_dom_idx on users_followed_dom (dom);
create index users_notified_notified_id_idx on users_notified (notified_id);

create index users_submitted_id_idx on users_submitted (id);
create index users_saved_id_idx on users_saved (id);
create index users_sticky_id_idx on users_sticky (id);
create index users_hidden_id_idx on users_hidden (id);
create index users_mail_id_idx on users_mail (id);
create index users_drafts_id_idx on users_drafts (id);
create index users_shareds_id_idx on users_shareds (id);
create index users_cubbeds_id_idx on users_cubbeds (id);
create index users_suggested_tags_id_idx on users_suggested_tags (id);
create index users_badged_id_idx on users_badged (id);
create index users_badges_id_idx on users_badges (id);
create index users_ignoring_id_idx on users_ignoring (id);
create index users_muting_id_idx on users_muting (id);
create index users_hushing_id_idx on users_hushing (id);
create index users_blocking_id_idx on users_blocking (id);
create index users_ignoring_tag_id_idx on users_ignoring_tag (id);
create index users_ignoring_dom_id_idx on users_ignoring_dom (id);
create index users_ignored_by_id_idx on users_ignored_by (id);
create index users_muted_by_id_idx on users_muted_by (id);
create index users_hushed_by_id_idx on users_hushed_by (id);
create index users_blocked_by_id_idx on users_blocked_by (id);
create index users_followed_id_idx on users_followed (id);
create index users_follower_id_idx on users_follower (id);
create index users_personal_tags_id_idx on users_personal_tags (id);
create index users_followed_tags_id_idx on users_followed_tags (id);
create index users_followed_dom_id_idx on users_followed_dom (id);
create index users_notified_id_idx on users_notified (id);

