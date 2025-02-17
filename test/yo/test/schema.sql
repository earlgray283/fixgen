CREATE TABLE Users (
  id INT64 NOT NULL,
  name STRING(MAX) NOT NULL,  
  created_at TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
  updated_at TIMESTAMP,
) PRIMARY KEY(id);

CREATE TABLE Todos (
  id INT64 NOT NULL,
  title STRING(MAX) NOT NULL,  
  description STRING(MAX) NOT NULL,  
  tags ARRAY<STRING(MAX)>,
  created_at TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
  updated_at TIMESTAMP,
  done_at TIMESTAMP,
) PRIMARY KEY(id);

