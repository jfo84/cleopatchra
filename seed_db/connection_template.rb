require 'pg'

class ConnectionTemplate
  attr_reader :db_name, :template_connection, :connection

  def initialize(db_name)
    @db_name = db_name
    initialize_template
    initialize_db
  end

  private

  def initialize_template
    # The template1 DB is always available
    @template_connection = PG::Connection.new(dbname: 'template1')
  end

  def initialize_db
    unless exists?
      template_connection.exec("CREATE DATABASE #{db_name}")
    end
    @connection = PG::Connection.new(dbname: db_name)
  end

  def exists?
    result = template_connection.exec('SELECT * from pg_database where datname = $1', [ db_name ])
    result.ntuples == 1
  end
end