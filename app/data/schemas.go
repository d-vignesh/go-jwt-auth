package data

// User table schema
var createUserTableSchema = `
create table if not exists users (
	id varchar(36) not null,
	email varchar(225) not null unique,
	username varchar(225),
	password varchar(225) not null,
	token varchar(15) not null,
	is_verified boolean default false,
	created_at timestamp not null,
	updated_at timestamp not null,
	primary key (id)
);
`

var createVerificationSchema = `
create table if not exists verifications (
	email 		Varchar(100) not null,
	code  		Varchar(10) not null,
	expiresat 	Timestamp not null,
	type        Varchar(10) not null,
	Primary Key (email),
	Constraint fk_user_email Foreign Key(email) References users(email)
		On Delete Cascade On Update Cascade
)
`
