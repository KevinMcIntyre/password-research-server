--! This is the SQL script for initializing the TU Password Research database !--

-- CREATE DATABASE tupwresearch OWNER postgres;

--! TABLES !--

CREATE TABLE collections (
    id BIGSERIAL NOT NULL,
    label VARCHAR(60) NOT NULL,
    creation_date TIMESTAMP,
    CONSTRAINT pk_collections PRIMARY KEY (id)
);

CREATE TABLE random_stage_images (
    id BIGSERIAL NOT NULL,
    image BYTEA,
    image_type VARCHAR(20) NOT NULL,
    alias VARCHAR(40),
    test_config_id INTEGER NOT NULL,
    stage_number INTEGER NOT NULL,
    row_number INTEGER NOT NULL,
    column_number INTEGER NOT NULL,
    creation_date TIMESTAMP,
    replacement_alias VARCHAR(255),
    CONSTRAINT pk_random_stage_images PRIMARY KEY (id)
);

CREATE TABLE saved_images (
    id BIGSERIAL NOT NULL,
    image BYTEA NOT NULL,
    image_type VARCHAR(20) NOT NULL,
    subject_id INTEGER,
    collection_id INTEGER,
    alias VARCHAR(40),
    creation_date TIMESTAMP,
    CONSTRAINT pk_saved_images PRIMARY KEY (id)
);

CREATE TABLE subjects (
    id BIGSERIAL NOT NULL,
    username VARCHAR(30),
    email VARCHAR(30),
    password VARCHAR(60),
    password_entropy double precision,
    password_strength INTEGER,
    pin_number VARCHAR(10),
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    birth_date DATE,
    creation_date TIMESTAMP,
    notes TEXT,
    CONSTRAINT pk_subjects PRIMARY KEY (id)
);

CREATE TABLE test_config_stage_images (
    id BIGSERIAL NOT NULL,
    image BYTEA,
    image_type VARCHAR(20) NOT NULL,
    stage_id INTEGER NOT NULL,
    stage_number INTEGER NOT NULL,
    alias VARCHAR(40),
    row_number INTEGER NOT NULL,
    column_number INTEGER NOT NULL,
    creation_date TIMESTAMP,
    CONSTRAINT pk_test_config_stage_images PRIMARY KEY (id)
);

CREATE TABLE test_config_stages (
    id BIGSERIAL NOT NULL,
    test_config_id INTEGER NOT NULL,
    stage_number INTEGER NOT NULL,
    creation_date TIMESTAMP,
    CONSTRAINT pk_test_config_stages PRIMARY KEY (id)
);

CREATE TABLE test_configs (
    id BIGSERIAL NOT NULL,
    name VARCHAR(60),
    rows_in_matrix INTEGER,
    cols_in_matrix INTEGER,
    stage_count INTEGER,
    image_may_not_be_present boolean,
    creation_date TIMESTAMP,
    CONSTRAINT pk_test_configs PRIMARY KEY (id)
);

CREATE TABLE uploaded_images (
    id BIGSERIAL NOT NULL,
    image BYTEA NOT NULL,
    image_type VARCHAR(20) NOT NULL,
    subject_id INTEGER,
    collection_id INTEGER,
    alias VARCHAR(40),
    creation_date TIMESTAMP,
    CONSTRAINT pk_uploaded_images PRIMARY KEY (id)
);

CREATE TABLE image_trial_images (
    id BIGSERIAL NOT NULL,
    image BYTEA,
    image_type VARCHAR(20) NOT NULL,
    trial_id INTEGER NOT NULL,
    stage_number INTEGER NOT NULL,
    alias VARCHAR(40),
    row_number INTEGER NOT NULL,
    column_number INTEGER NOT NULL,
    is_user_image BOOLEAN,
    CONSTRAINT pk_image_trial_images PRIMARY KEY (id)
);

CREATE TABLE image_trial_stage_results (
    id BIGSERIAL NOT NULL,
    trial_id INTEGER NOT NULL,
    stage_number INTEGER NOT NULL,
    passed_auth BOOLEAN,
    selected_trial_image_id INTEGER,
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    CONSTRAINT pk_image_trial_stage_results PRIMARY KEY (id)
);

CREATE TABLE image_trials (
    id BIGSERIAL NOT NULL,
    subject_id INTEGER NOT NULL,
    test_config_id INTEGER NOT NULL,
    notes TEXT,
    creation_date TIMESTAMP,
    CONSTRAINT pk_image_trials PRIMARY KEY (id)
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

--! SEQUENCES !--

CREATE SEQUENCE collections_seq;
CREATE SEQUENCE random_stage_images_seq;
CREATE SEQUENCE saved_images_seq;
CREATE SEQUENCE subjects_seq;
CREATE SEQUENCE test_config_stage_images_seq;
CREATE SEQUENCE test_config_stages_seq;
CREATE SEQUENCE test_configs_seq;
CREATE SEQUENCE uploaded_images_seq;
CREATE SEQUENCE image_trial_images_seq;
CREATE SEQUENCE image_trial_stage_results_seq;
CREATE SEQUENCE image_trials_seq;
-- CREATE SEQUENCE pin_trials_seq;
-- CREATE SEQUENCE string_trials_seq;

--! FOREIGN KEY REFERENCES !--
 ALTER TABLE image_trial_images ADD CONSTRAINT fk_image_trial_images_trial_id_01 FOREIGN KEY (trial_id) REFERENCES image_trials(id);
 ALTER TABLE image_trial_stage_results ADD CONSTRAINT fk_image_trial_stage_results_selected_image_id_01 FOREIGN KEY (selected_trial_image_id) REFERENCES image_trial_images(id);
 ALTER TABLE image_trial_stage_results ADD CONSTRAINT fk_image_trial_stage_results_trial_id_01 FOREIGN KEY (trial_id) REFERENCES image_trials(id);
 ALTER TABLE image_trials ADD CONSTRAINT fk_image_trials_config_id_01 FOREIGN KEY (test_config_id) REFERENCES test_configs(id);
 ALTER TABLE image_trials ADD CONSTRAINT fk_image_trials_subject_id_01 FOREIGN KEY (subject_id) REFERENCES subjects(id);
 ALTER TABLE random_stage_images ADD CONSTRAINT fk_random_stage_images_01 FOREIGN KEY (test_config_id) REFERENCES test_configs(id);
 ALTER TABLE saved_images ADD CONSTRAINT fk_saved_images_collection_id_01 FOREIGN KEY (collection_id) REFERENCES collections(id);
 ALTER TABLE saved_images ADD CONSTRAINT fk_saved_images_subject_id_01 FOREIGN KEY (subject_id) REFERENCES subjects(id);
 ALTER TABLE test_config_stage_images ADD CONSTRAINT fk_config_stage_images_stage_id_01 FOREIGN KEY (stage_id) REFERENCES test_config_stages(id);
 ALTER TABLE test_config_stages ADD CONSTRAINT fk_test_config_stages_01 FOREIGN KEY (test_config_id) REFERENCES test_configs(id);
 ALTER TABLE uploaded_images ADD CONSTRAINT fk_uploaded_images_collection_id_01 FOREIGN KEY (collection_id) REFERENCES collections(id);
 ALTER TABLE uploaded_images ADD CONSTRAINT fk_uploaded_images_subject_id_01 FOREIGN KEY (subject_id) REFERENCES subjects(id);


--! CREATE INDICES !--

CREATE UNIQUE INDEX image_trial_id ON image_trials (id);
CREATE UNIQUE INDEX image_trial_result_id ON image_trial_stage_results (id);
CREATE UNIQUE INDEX image_trial_image_id ON image_trial_images (id);
CREATE UNIQUE INDEX collection_id ON collections (id);
CREATE UNIQUE INDEX random_stage_image_id ON random_stage_images (id);
CREATE UNIQUE INDEX saved_image_id ON saved_images (id);
CREATE UNIQUE INDEX subject_id ON subjects (id);
CREATE UNIQUE INDEX test_config_id ON test_configs (id);
CREATE UNIQUE INDEX test_config_stage_id ON test_config_stages (id);
CREATE UNIQUE INDEX test_config_stage_image_id ON test_config_stage_images (id);
CREATE UNIQUE INDEX uploaded_image_id ON uploaded_images (id);

--! SEED DATA !--

INSERT INTO subjects (id, username, email)
VALUES
  (0, 'NULL', 'NULL');

INSERT INTO collections (id, label)
VALUES
  (0, 'NULL');

CREATE FUNCTION duplicate_config_images(trial_id int, config_id int) RETURNS VOID AS $$
BEGIN
    EXECUTE 'INSERT INTO image_trial_images (image, image_type, trial_id, stage_number, alias, row_number, column_number)
    SELECT
    image,
    image_type,
    $1,
    stage_number,
    CASE WHEN alias = ''user-img''
        THEN ''user-img''
        ELSE replace(md5(random() :: TEXT || clock_timestamp() :: TEXT), '-' :: TEXT, '' :: TEXT) :: VARCHAR(60)
    END AS alias,
    row_number,
    column_number
    FROM test_config_stage_images
    WHERE stage_id IN (
        SELECT id
        FROM test_config_stages
        WHERE test_config_id = $2
    );'
    USING trial_id, config_id;
END;
$$
LANGUAGE plpgsql;

CREATE FUNCTION create_image_trial(subject_id int, config_id int, number_of_stages int) RETURNS INTEGER AS $$
DECLARE
    trial_id int;
    current_stage_number int := 1;
BEGIN
	INSERT INTO image_trials (subject_id, test_config_id, creation_date)
	VALUES ($1, $2, now())
	RETURNING id INTO trial_id;

    LOOP
        EXIT WHEN current_stage_number > $3;
        INSERT INTO image_trial_stage_results (trial_id, stage_number)
        VALUES (trial_id, current_stage_number);
        current_stage_number := current_stage_number + 1;
    END LOOP;

    PERFORM duplicate_config_images(trial_id, $2);
    RETURN trial_id;
END;
$$
LANGUAGE plpgsql;

-- CREATE FUNCTION submit_image_selection(trial_id int, stage_number int, selected_alias VARCHAR(255), subject_id int) RETURNS void AS $$
-- DECLARE
--     test_config_stage_id int := (SELECT id FROM test_config_stages WHERE stage_number = $2);
-- BEGIN
--   IF EXISTS (SELECT id FROM saved_images WHERE alias = $3 AND subject_id = $4) THEN
--     UPDATE image_trial_stage_results
--     SET passed_auth = true,
--         correct_saved_image_id = (SELECT id FROM saved_images WHERE alias = $3 AND subject_id = $4),
--         end_time = now()
--     WHERE trial_id = $1 AND config_stage_id = test_config_stage_id;
--   ELSE
--     UPDATE image_trial_stage_results
--     SET passed_auth = false,
--         selected_trial_image_id = (SELECT id FROM test_config_stage_images WHERE alias = $3 AND stage_id = test_config_stage_id),
--         end_time = now()
--     WHERE trial_id = $1 AND config_stage_id = test_config_stage_id;
--   END IF;
-- END;
-- $$
-- LANGUAGE plpgsql;