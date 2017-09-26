require 'virtus'
require_relative './connection'

module DataMixin

  def self.included(model)
    model.attribute :data_hash, Hash
  end
  
  def id
    data_hash['id']
  end

  private

  def connection
    @@connection ||= Connection.new('cleopatchra')
  end
  
  def data
    data_hash.to_json
  end

  def is_dup?
    result = connection.exec("SELECT id FROM #{table_name} WHERE id = $1", [id])
    result.ntuples == 1
  end

  def table_name
    "#{self.class.to_s.downcase}s"
  end
end