require_relative './data_mixin'

class Comment
  include Virtus.model
  include DataMixin

  attribute :pull_id, String
  
  def record
    connection.exec("INSERT INTO comments (id, data, pull_id) VALUES ($1, $2, $3)", [id, data, pull_id])
  end
end