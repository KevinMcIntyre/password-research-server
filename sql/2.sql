--! This is the SQL script for initializing the TU Password Research database !--

-- CREATE DATABASE tupwresearch OWNER postgres;

--! TABLES !--

CREATE TABLE subjects (
  id                BIGSERIAL NOT NULL,
  username          VARCHAR(30),
  email             VARCHAR(30) UNIQUE,
  password          VARCHAR(60),
  password_entropy  FLOAT,
  password_strength INTEGER,
  pin_number        VARCHAR(10),
  first_name        VARCHAR(50),
  last_name         VARCHAR(50),
  birth_date        DATE,
  creation_date     TIMESTAMP,
  notes             TEXT,
  CONSTRAINT pk_subjects PRIMARY KEY (id)
);

CREATE TABLE test_configs (
  id                       BIGSERIAL   NOT NULL,
  name                     VARCHAR(60) NOT NULL UNIQUE,
  rows_in_matrix           INTEGER,
  cols_in_matrix           INTEGER,
  stage_count              INTEGER,
  image_may_not_be_present BOOLEAN,
  creation_date            TIMESTAMP,
  CONSTRAINT pk_test_config PRIMARY KEY (id)
);

CREATE TABLE test_config_stages (
  id             BIGSERIAL NOT NULL,
  test_config_id INTEGER   NOT NULL,
  stage_number   INTEGER   NOT NULL,
  creation_date  TIMESTAMP,
  CONSTRAINT pk_custom_stages PRIMARY KEY (id)
);

CREATE TABLE test_config_stage_images (
  id            BIGSERIAL   NOT NULL,
  image         BYTEA       NOT NULL,
  image_type    VARCHAR(20) NOT NULL,
  stage_id      INTEGER     NOT NULL,
  alias         VARCHAR(40),
  row_number    INTEGER     NOT NULL,
  column_number INTEGER     NOT NULL,
  creation_date TIMESTAMP,
  CONSTRAINT pk_test_config_stage_images PRIMARY KEY (id)
);


CREATE TABLE random_stage_images (
  id             BIGSERIAL   NOT NULL,
  image          BYTEA       NOT NULL,
  image_type     VARCHAR(20) NOT NULL,
  alias          VARCHAR(40),
  test_config_id INTEGER     NOT NULL,
  stage_number   INTEGER     NOT NULL,
  row_number     INTEGER     NOT NULL,
  column_number  INTEGER     NOT NULL,
  creation_date  TIMESTAMP,
  CONSTRAINT pk_random_stage_images PRIMARY KEY (id)
);

CREATE TABLE image_trials (
  id             BIGSERIAL NOT NULL,
  subject_id     INTEGER   NOT NULL,
  test_config_id INTEGER   NOT NULL,
  passed_auth    BOOLEAN,
  start_time     TIMESTAMP,
  end_time       TIMESTAMP,
  notes          TEXT,
  CONSTRAINT pk_image_trials PRIMARY KEY (id)
);

CREATE TABLE image_trial_stage_results (
  id                 BIGSERIAL NOT NULL,
  trial_id           INTEGER   NOT NULL,
  stage_id           INTEGER   NOT NULL,
  selected_image_id  INTEGER   NOT NULL,
  correct_image_id   INTEGER,
  start_time         TIMESTAMP,
  end_time           TIMESTAMP,
  CONSTRAINT pk_image_trial_stage_results PRIMARY KEY (id)
);

CREATE TABLE collections (
  id            BIGSERIAL   NOT NULL,
  label         VARCHAR(60) NOT NULL UNIQUE,
  creation_date TIMESTAMP,
  CONSTRAINT pk_collections PRIMARY KEY (id)
);

CREATE TABLE saved_images (
  id            BIGSERIAL   NOT NULL,
  image         BYTEA       NOT NULL,
  image_type    VARCHAR(20) NOT NULL,
  subject_id    INTEGER,
  collection_id INTEGER,
  alias         VARCHAR(40),
  creation_date TIMESTAMP,
  CONSTRAINT pk_saved_images PRIMARY KEY (id)
);

CREATE TABLE uploaded_images (
  id            BIGSERIAL   NOT NULL,
  image         BYTEA       NOT NULL,
  image_type    VARCHAR(20) NOT NULL,
  subject_id    INTEGER,
  collection_id INTEGER,
  alias         VARCHAR(40),
  creation_date TIMESTAMP,
  CONSTRAINT pk_uploaded_images PRIMARY KEY (id)
);

-- CREATE TABLE pin_trials (
--   id                INTEGER       NOT NULL,
--   user_id           INTEGER       NOT NULL,
--   creation_date     TIMESTAMP,
--   notes             TEXT,
--   CONSTRAINT pk_pin_trials PRIMARY KEY (id)
-- );
--
-- CREATE TABLE string_trials (
--   id                INTEGER       NOT NULL,
--   user_id           INTEGER       NOT NULL,
--   creation_date     TIMESTAMP,
--   notes             TEXT,
--   CONSTRAINT pk_string_trials PRIMARY KEY (id)
-- );

--! JOIN TABLES !--


--! SEQUENCES !--

CREATE SEQUENCE subjects_seq;
CREATE SEQUENCE test_configs_seq;
CREATE SEQUENCE test_config_stage_images_seq;
CREATE SEQUENCE random_stage_images_seq;
CREATE SEQUENCE image_trials_seq;
CREATE SEQUENCE image_trial_stage_results_seq;
CREATE SEQUENCE collections_seq;
CREATE SEQUENCE saved_images_seq;
CREATE SEQUENCE uploaded_images_seq;
-- CREATE SEQUENCE pin_trials_seq;
-- CREATE SEQUENCE string_trials_seq;

--! FOREIGN KEY REFERENCES !--
ALTER TABLE test_config_stages ADD CONSTRAINT fk_test_config_stages_01 FOREIGN KEY (test_config_id) REFERENCES test_configs (id);
ALTER TABLE test_config_stage_images ADD CONSTRAINT fk_test_config_stage_images_01 FOREIGN KEY (stage_id) REFERENCES test_config_stages (id);
ALTER TABLE random_stage_images ADD CONSTRAINT fk_random_stage_images_01 FOREIGN KEY (test_config_id) REFERENCES test_configs (id);
ALTER TABLE image_trials ADD CONSTRAINT fk_image_trials_subject_id_01 FOREIGN KEY (subject_id) REFERENCES subjects (id);
ALTER TABLE image_trials ADD CONSTRAINT fk_image_trials_test_config_id_01 FOREIGN KEY (test_config_id) REFERENCES test_configs (id);
ALTER TABLE image_trial_stage_results ADD CONSTRAINT fk_image_trial_stage_results_trial_id_01 FOREIGN KEY (trial_id) REFERENCES image_trials (id);
ALTER TABLE image_trial_stage_results ADD CONSTRAINT fk_image_trial_stage_results_stage_id_01 FOREIGN KEY (stage_id) REFERENCES test_config_stages (id);
ALTER TABLE image_trial_stage_results ADD CONSTRAINT fk_image_trial_stage_results_selected_image_id_01 FOREIGN KEY (selected_image_id) REFERENCES test_config_stage_images (id);
ALTER TABLE image_trial_stage_results ADD CONSTRAINT fk_image_trial_stage_results_correct_image_id_01 FOREIGN KEY (correct_image_id) REFERENCES uploaded_images (id);
ALTER TABLE saved_images ADD CONSTRAINT fk_saved_images_subject_id_01 FOREIGN KEY (subject_id) REFERENCES subjects (id);
ALTER TABLE saved_images ADD CONSTRAINT fk_saved_images_collection_id_01 FOREIGN KEY (collection_id) REFERENCES collections (id);
ALTER TABLE uploaded_images ADD CONSTRAINT fk_uploaded_images_subject_id_01 FOREIGN KEY (subject_id) REFERENCES subjects (id);
ALTER TABLE uploaded_images ADD CONSTRAINT fk_uploaded_images_collection_id_01 FOREIGN KEY (collection_id) REFERENCES collections (id);


--! CREATE INDICES !--
CREATE UNIQUE INDEX subject_id ON subjects (id);
CREATE UNIQUE INDEX test_config_id ON test_configs (id);
CREATE UNIQUE INDEX test_config_stage_id ON test_config_stages (id);
CREATE UNIQUE INDEX test_config_stage_image_id ON test_config_stage_images (id);
CREATE UNIQUE INDEX random_stage_image_id ON random_stage_images (id);
CREATE UNIQUE INDEX image_trial_id ON image_trials (id);
CREATE UNIQUE INDEX image_trial_result_id ON image_trial_stage_results (id);
CREATE UNIQUE INDEX collection_id ON collections (id);
CREATE UNIQUE INDEX saved_image_id ON saved_images (id);
CREATE UNIQUE INDEX uploaded_image_id ON uploaded_images (id);


INSERT INTO subjects (id, username, email)
VALUES
  (0, 'NULL', 'NULL');

INSERT INTO collections (id, label)
VALUES
  (0, 'NULL');
