-- Write your migrate up statements here
CREATE TABLE project_techs
(
    id   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE project_category
(
    id   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE projects
(
    id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name              VARCHAR(255) NOT NULL UNIQUE,
    description       TEXT,
    preview_image_url TEXT,
    linked_link       VARCHAR(255),
    github_link       VARCHAR(255),
    date              TIMESTAMP,
    client            VARCHAR(255),
    images_links      TEXT[]
);

-- a project can have multiple techs, a tech can be used in multiple projects
CREATE TABLE project_techs_projects
(
    project_id uuid REFERENCES projects (id) ON DELETE CASCADE,
    tech_id    uuid REFERENCES project_techs (id),
    PRIMARY KEY (project_id, tech_id)
);

-- a project can have multiple categories, a category can be used in multiple projects
CREATE TABLE project_category_projects
(
    project_id  uuid REFERENCES projects (id) ON DELETE CASCADE,
    category_id uuid REFERENCES project_category (id),
    PRIMARY KEY (project_id, category_id)
);

---- create above / drop below ----
DROP TABLE IF EXISTS project_techs_projects;
DROP TABLE IF EXISTS project_category_projects;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS project_category;
DROP TABLE IF EXISTS project_techs;

