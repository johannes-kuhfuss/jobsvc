DROP TABLE joblist;

CREATE TABLE joblist (
	"id" varchar NOT NULL,
	"correlation_id" varchar NULL,
	"name" varchar NULL,
	"created_at" timestamptz NULL,
	"created_by" varchar NULL,
	"modified_at" timestamptz NULL,
	"modified_by" varchar NULL,
	"status" varchar NULL,
	"source" varchar NULL,
	"destination" varchar NULL,
	"type" varchar NULL,
	"sub_type" varchar NULL,
	"action" varchar NULL,
	"action_details" varchar NULL,
	"progress" int4 NULL,
	"history" jsonb NULL,
	"extra_data" varchar NULL,
	"priority" int4 NULL,
	"rank" int4 NULL,
	CONSTRAINT joblist_pk PRIMARY KEY (id)
);