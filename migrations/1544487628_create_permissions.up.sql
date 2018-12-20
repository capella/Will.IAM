CREATE TABLE IF NOT EXISTS permissions (
	id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
  role_id UUID NOT NULL,
  ownership_level VARCHAR(10) NOT NULL,
  action VARCHAR(200) NOT NULL,
  service VARCHAR(200) NOT NULL,
  resource_hierarchy VARCHAR(1000),
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  FOREIGN KEY(role_id) REFERENCES roles (id)
);
