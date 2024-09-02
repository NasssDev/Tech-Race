

DROP TABLE IF EXISTS "collision";
DROP SEQUENCE IF EXISTS collision_id_seq;
CREATE SEQUENCE collision_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 785 CACHE 1;

CREATE TABLE "public"."collision" (
    "id" integer DEFAULT nextval('collision_id_seq') NOT NULL,
    "distance" numeric NOT NULL,
    "is_collision" boolean NOT NULL,
    "timestamp" timestamp NOT NULL,
    "id_session" integer NOT NULL,
    CONSTRAINT "collision_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

DROP TABLE IF EXISTS "linetracking";
DROP SEQUENCE IF EXISTS "lineTracking_id_seq";
CREATE SEQUENCE "lineTracking_id_seq" INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."linetracking" (
    "id" integer DEFAULT nextval('"lineTracking_id_seq"') NOT NULL,
    "line_tracking_value" integer NOT NULL,
    "id_session" integer NOT NULL,
    CONSTRAINT "lineTracking_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

DROP TABLE IF EXISTS "session";
DROP SEQUENCE IF EXISTS session_id_seq;
CREATE SEQUENCE session_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 25 CACHE 1;

CREATE TABLE "public"."session" (
    "id" integer DEFAULT nextval('session_id_seq') NOT NULL,
    "start_time" timestamp,
    "end_time" timestamp,
    "is_autopilot" boolean,
    CONSTRAINT "session_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

DROP TABLE IF EXISTS "speed";
DROP SEQUENCE IF EXISTS speed_id_seq;
CREATE SEQUENCE speed_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."speed" (
    "id" integer DEFAULT nextval('speed_id_seq') NOT NULL,
    "speed" integer NOT NULL,
    "timestamp" timestamp NOT NULL,
    "id_session" integer NOT NULL,
    CONSTRAINT "speed_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

DROP TABLE IF EXISTS "video";
DROP SEQUENCE IF EXISTS video_id_seq;
CREATE SEQUENCE video_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."video" (
    "id" integer DEFAULT nextval('video_id_seq') NOT NULL,
    "video_url" character varying NOT NULL,
    "id_session" integer NOT NULL,
    CONSTRAINT "video_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

ALTER TABLE ONLY "public"."collision" ADD CONSTRAINT "collision_id_session_fkey" FOREIGN KEY (id_session) REFERENCES session(id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."linetracking" ADD CONSTRAINT "lineTracking_id_session_fkey" FOREIGN KEY (id_session) REFERENCES session(id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."speed" ADD CONSTRAINT "speed_id_session_fkey" FOREIGN KEY (id_session) REFERENCES session(id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."video" ADD CONSTRAINT "video_id_session_fkey" FOREIGN KEY (id_session) REFERENCES session(id) NOT DEFERRABLE;

-- 2024-07-04 09:50:57.138037+00


