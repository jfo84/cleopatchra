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
    @connection ||= Connection.new('cleopatchra')
  end
  
  def data
    data_hash.to_s
  end
end