require 'virtus'

module DataMixin

  def self.included(model)
    model.attribute :data_hash, Hash
  end
  
  def id
    data_hash['id']
  end

  private
  
  def data
    data_hash.to_json
  end
end