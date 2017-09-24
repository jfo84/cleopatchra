require 'pg'

class Connection
  attr_reader :db_name, :template_connection, :connection

  def initialize(db_name)
    @db_name = db_name
    initialize_template
    initialize_db
  end

  def exec(*args)
    connection.exec(*args)
  end

  private

  def initialize_template
    # The template1 DB is always available
    @template_connection = PG::Connection.new(dbname: 'template1')
  end

  def initialize_db
    unless exists?
      template_connection.exec("CREATE DATABASE #{db_name}")
      connection = PG::Connection.new(dbname: db_name)
      initialize_tables(connection)
      index_tables(connection)
    end
    @connection = connection || PG::Connection.new(dbname: db_name)
  end

  def exists?
    result = template_connection.exec('SELECT * from pg_database where datname = $1', [ db_name ])
    result.ntuples == 1
  end

  # The tables are relatively static. I plan to keep these
  # set in stone and query into the JSON in C++ if I have to
  def initialize_tables(connection)
    connection.exec('CREATE TABLE repos (
      id integer PRIMARY KEY,
      data jsonb NOT NULL)')
    connection.exec('CREATE TABLE pulls (
      id integer PRIMARY KEY,
      data jsonb NOT NULL,
      repo_id text NOT NULL)')
    connection.exec("CREATE TABLE comments (
      id integer PRIMARY KEY,
      data jsonb NOT NULL,
      pull_id integer NOT NULL)")
    connection.exec("CREATE TABLE reviews (
      id integer PRIMARY KEY,
      data jsonb NOT NULL,
      pull_id integer NOT NULL)")
  end

  def index_tables(connection)
    ['repos', 'pulls', 'comments', 'reviews'].each do |table_name|
      connection.exec("CREATE UNIQUE INDEX index_#{table_name}_on_id ON #{table_name} (id)")
    end
  end
end