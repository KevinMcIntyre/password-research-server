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
    sex VARCHAR(1),
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
    image_type VARCHAR(20),
    trial_id INTEGER,
    stage_number INTEGER,
    alias VARCHAR(40),
    row_number INTEGER,
    column_number INTEGER,
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

CREATE TABLE password_trials (
  id                BIGSERIAL     NOT NULL,
  subject_id        INTEGER       NOT NULL,
  trial_type        VARCHAR(8)    NOT NULL,
  attempts_allowed  INTEGER       NOT NULL,
  passed_auth       BOOLEAN,
  start_time        TIMESTAMP,
  end_time          TIMESTAMP,
  creation_date     TIMESTAMP,
  notes             TEXT,
  CONSTRAINT pk_password_trials PRIMARY KEY (id)
);

CREATE TABLE passwords_submitted (
    id                  BIGSERIAL       NOT NULL,
    trial_id            INTEGER         NOT NULL,
    password_entered    TEXT            NOT NULL,
    attempt_number      INTEGER         NOT NULL,
    submission_time     TIMESTAMP       NOT NULL,
    CONSTRAINT pk_passwords_submitted PRIMARY KEY (id)
);

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
CREATE SEQUENCE password_trials_seq;
CREATE SEQUENCE passwords_submitted_seq;

--! FOREIGN KEY REFERENCES !--
 ALTER TABLE passwords_submitted ADD CONSTRAINT fk_passwords_submitted_trial_id_01 FOREIGN KEY (trial_id) REFERENCES password_trials(id);
 ALTER TABLE password_trials ADD CONSTRAINT fk_password_trials_subject_id_01 FOREIGN KEY (subject_id) REFERENCES subjects(id);
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

CREATE UNIQUE INDEX passwords_submitted_id ON passwords_submitted (id);
CREATE UNIQUE INDEX password_trial_id ON password_trials (id);
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

CREATE INDEX password_trial_type ON password_trials (trial_type);

--! SEED DATA !--

INSERT INTO subjects (id, username, email)
VALUES
  (0, 'NULL', 'NULL');

INSERT INTO collections (id, label)
VALUES
  (0, 'NULL');

INSERT INTO image_trial_images (id, image_type)
VALUES
    (0, 'IMAGE_NOT_PRESENT');

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
        ELSE regexp_replace(md5(random() :: TEXT || clock_timestamp() :: TEXT) :: VARCHAR(60), ''[^a-zA-Z0-9]+$|-'', '''', ''g'')
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

CREATE FUNCTION submit_image_selection(image_trial_id int, stage int, selected_alias VARCHAR(255), submit_time TIMESTAMP) RETURNS RECORD AS $$
DECLARE
    result RECORD;
BEGIN
    IF (selected_alias = 'no-pass-image')
        THEN
            UPDATE image_trial_stage_results
            SET selected_trial_image_id = subquery.image_id,
                passed_auth = subquery.passed_auth,
                end_time = $4
            FROM (
                SELECT 0 AS image_id,
                CASE WHEN (
                    SELECT id FROM image_trial_images WHERE trial_id = $1 AND stage_number = $2 AND is_user_image = TRUE
                ) IS NULL
                    THEN TRUE
                    ELSE FALSE
                END as passed_auth
            ) subquery
            WHERE trial_id = $1 AND stage_number = $2;
        ELSE
            UPDATE image_trial_stage_results
            SET selected_trial_image_id = subquery.image_id,
                passed_auth = subquery.passed_auth,
                end_time = $4
            FROM (
                SELECT
                id AS image_id,
                CASE WHEN is_user_image IS NULL
                    THEN FALSE
                    ELSE is_user_image
                END AS passed_auth
                FROM image_trial_images
                WHERE trial_id = $1
                AND stage_number = $2
                AND alias = $3
            ) subquery
            WHERE trial_id = $1 AND stage_number = $2;
    END IF;
    IF EXISTS (
        SELECT config.stage_count
        FROM image_trials trial
        JOIN test_configs config ON config.id = trial.test_config_id
        WHERE trial.id = $1 AND config.stage_count > $2
    ) THEN
        UPDATE image_trial_stage_results
        SET start_time = $4
        WHERE trial_id = $1 AND stage_number = ($2 + 1);
        SELECT FALSE AS trial_complete, FALSE AS succesful_auth INTO result;
    ELSE
        IF EXISTS (
            SELECT false
            FROM image_trial_stage_results
            WHERE trial_id = $1 AND (passed_auth = false OR passed_auth IS NULL)
            LIMIT 1
        ) THEN
            SELECT TRUE AS trial_complete, FALSE AS succesful_auth INTO result;
        ELSE
            SELECT TRUE AS trial_complete, TRUE AS succesful_auth INTO result;
        END IF;
    END IF;
    RETURN result;
END;
$$
LANGUAGE plpgsql;

CREATE FUNCTION submit_password_submission(trial_id int, password TEXT, submit_time TIMESTAMP) RETURNS RECORD AS $$
DECLARE
    result RECORD;
BEGIN
    --! Check if trial is already complete, if so, return the result !--
    IF (
	SELECT
	CASE WHEN passed_auth IS NULL
		THEN TRUE
		ELSE FALSE
	END AS incomplete
	FROM password_trials WHERE id = $1
    )
        THEN
		--! Insert the password submission !--
		INSERT INTO passwords_submitted(trial_id, password_entered, attempt_number, submission_time)
		SELECT
		$1,
		$2,
		attempt_count.attempts + 1,
		$3
		FROM (
		SELECT
		COUNT(id) AS attempts
		FROM passwords_submitted
		WHERE passwords_submitted.trial_id = $1
		) AS attempt_count;

		--! Check if authentication is successful and whether or not they have more attempts !--
		WITH submission_info AS (
		SELECT
		CASE WHEN attempts.count >= pt.attempts_allowed
		    THEN false
		    ELSE true
		END AS more_attempts,
		CASE WHEN pt.trial_type = 'password'
		    THEN
			CASE WHEN s.password IN (
			    SELECT 
			    password_entered
			    FROM passwords_submitted
			    WHERE passwords_submitted.trial_id = $1
			)
			    THEN true
			    ELSE false
			END
		    ELSE
			CASE WHEN s.pin_number IN (
			    SELECT 
			    password_entered
			    FROM passwords_submitted
			    WHERE passwords_submitted.trial_id = $1
			)
			    THEN true
			    ELSE false
			END
		END AS successful_auth
		FROM password_trials pt
		JOIN (
		    SELECT
		    $1 AS trial_id,
		    COUNT(id) AS count
		    FROM passwords_submitted
		    WHERE passwords_submitted.trial_id = $1
		) attempts ON attempts.trial_id = pt.id
		JOIN subjects s ON s.id = pt.subject_id
		WHERE pt.id = $1
		)
		SELECT
		CASE WHEN more_attempts = TRUE
		THEN
		    CASE WHEN successful_auth = TRUE
			THEN TRUE
			ELSE FALSE
		    END
		ELSE
		    TRUE
		END AS trial_complete,
		successful_auth 
		FROM submission_info INTO result;
		IF (result.trial_complete)
		THEN
			UPDATE password_trials
			SET end_time = now(),
			    passed_auth = result.successful_auth
			WHERE password_trials.id = $1;
		END IF;
	ELSE
		SELECT
		TRUE AS trial_complete,
		passed_auth AS successful_auth
		FROM password_trials
		WHERE id = $1 INTO result;
		
    END IF;
    RETURN result;
END;
$$
LANGUAGE plpgsql;